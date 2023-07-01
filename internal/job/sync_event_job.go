package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain"
	"github.com/aasumitro/tix/internal/domain/response"
	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
	"time"
)

type syncEventJob struct {
	service     domain.ITixService
	redisClient *redis.Client
}

func NewEventJob(
	service domain.ITixService,
	redisClient *redis.Client,
) {
	event := &syncEventJob{service, redisClient}
	event.regisCronJob()
	event.regisQueueSubscriber()
}

func (e *syncEventJob) regisCronJob() {
	scheduler := gocron.NewScheduler(time.UTC)
	_, _ = scheduler.Every(common.EventRemovalScheduleTime).Minute().Do(func() {
		// Retrieve data from Redis cache
		cacheData, err := e.redisClient.Get(context.Background(), common.AutoSyncEventKey).Result()
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				return
			}
			cacheData = ""
		}
		if cacheData == "" {
			return
		}
		// Unmarshal cache data into a slice of AutoSyncRespond
		var cacheItems []*response.AutoSyncRespond
		if err = json.Unmarshal([]byte(cacheData), &cacheItems); err != nil {
			return
		}
		// remove expired event
		for i, event := range cacheItems {
			if int64(event.EventDate) < time.Now().Unix() {
				cacheItems = append(cacheItems[:i], cacheItems[i+1:]...)
			}
		}
		// convert to string
		jsonData, err := json.Marshal(cacheItems)
		if err != nil {
			return
		}
		// Set the updated cache data in Redis
		if err = e.redisClient.Set(
			context.Background(),
			common.AutoSyncEventKey,
			jsonData, -1,
		).Err(); err != nil {
			return
		}
	})
	_, _ = scheduler.Every(common.EventSyncScheduleTime).Minute().Do(func() {
		// Retrieve data from Redis cache
		cacheData, err := e.redisClient.Get(context.Background(), common.AutoSyncEventKey).Result()
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				return
			}
			cacheData = ""
		}
		if cacheData == "" {
			return
		}
		// Unmarshal cache data into a slice of AutoSyncRespond
		var cacheItems []*response.AutoSyncRespond
		if err = json.Unmarshal([]byte(cacheData), &cacheItems); err != nil {
			return
		}
		// remove expired event
		go func() {
			for _, event := range cacheItems {
				_ = e.service.SyncRespondData(context.Background(), event.FormID)
			}
		}()
	})
	scheduler.StartAsync()
}

func (e *syncEventJob) regisQueueSubscriber() {
	e.subscribeSyncEvent()
	e.subscribeGenerateTicket()
	e.subscribeExportData()
}

func (e *syncEventJob) subscribeSyncEvent() {
	subscriber := e.redisClient.Subscribe(context.Background(), common.ReqSyncEventQueueKey)
	go func() {
		for {
			msg, err := subscriber.ReceiveMessage(context.Background())
			if err != nil {
				ptn := "[%d] - SYNC_EVENT_ERR (QUEUE): %s"
				msg := fmt.Sprintf(ptn, time.Now().Unix(), err.Error())
				sentry.CaptureMessage(msg)
				continue
			}

			var eventData struct {
				GoogleFormID string `json:"google_form_id"`
			}

			if err := json.Unmarshal([]byte(msg.Payload), &eventData); err != nil {
				ptn := "[%d] - SYNC_EVENT_ERR (DECODE): %s"
				msg := fmt.Sprintf(ptn, time.Now().Unix(), err.Error())
				sentry.CaptureMessage(msg)
				continue
			}

			if err := e.service.SyncRespondData(context.Background(), eventData.GoogleFormID); err != nil {
				ptn := "[%d] - SYNC_EVENT_ERR (ACTION): %s"
				msg := fmt.Sprintf(ptn, time.Now().Unix(), err.Error())
				sentry.CaptureMessage(msg)
				continue
			}
		}
	}()
}

func (e *syncEventJob) subscribeGenerateTicket() {
	subscriber := e.redisClient.Subscribe(context.Background(), common.ReqGenEventTixQueueKey)
	go func() {
		for {
			msg, err := subscriber.ReceiveMessage(context.Background())
			if err != nil {
				ptn := "[%d] - GEN_TIX_ERR (QUEUE): %s"
				msg := fmt.Sprintf(ptn, time.Now().Unix(), err.Error())
				sentry.CaptureMessage(msg)
				continue
			}

			var eventData struct {
				GoogleFormID  string `json:"google_form_id"`
				ParticipantID int32  `json:"participant_id"`
			}

			if err := json.Unmarshal([]byte(msg.Payload), &eventData); err != nil {
				ptn := "[%d] - GEN_TIX_ERR (DECODE): %s"
				msg := fmt.Sprintf(ptn, time.Now().Unix(), err.Error())
				sentry.CaptureMessage(msg)
				continue
			}

			if err := e.service.GenerateTicket(
				context.Background(), eventData.GoogleFormID,
				eventData.ParticipantID,
			); err != nil {
				ptn := "[%d] - EXPORT_DATA_ERR (ACTION): %s"
				msg := fmt.Sprintf(ptn, time.Now().Unix(), err.Error())
				sentry.CaptureMessage(msg)
				continue
			}
		}
	}()
}

func (e *syncEventJob) subscribeExportData() {
	subscriber := e.redisClient.Subscribe(context.Background(), common.ReqExpEventDataQueueKey)
	go func() {
		for {
			msg, err := subscriber.ReceiveMessage(context.Background())
			if err != nil {
				ptn := "[%d] - EXPORT_DATA_ERR (QUEUE): %s"
				msg := fmt.Sprintf(ptn, time.Now().Unix(), err.Error())
				sentry.CaptureMessage(msg)
				continue
			}

			var eventData struct {
				GoogleFormID string `json:"google_form_id"`
				ExportType   string `json:"export_type"`
				Email        string `json:"email"`
			}

			if err := json.Unmarshal([]byte(msg.Payload), &eventData); err != nil {
				ptn := "[%d] - EXPORT_DATA_ERR (DECODE): %s"
				msg := fmt.Sprintf(ptn, time.Now().Unix(), err.Error())
				sentry.CaptureMessage(msg)
				continue
			}

			if err := e.service.ExportEvent(
				context.Background(), eventData.GoogleFormID,
				eventData.ExportType, eventData.Email,
			); err != nil {
				ptn := "[%d] - EXPORT_DATA_ERR (ACTION): %s"
				msg := fmt.Sprintf(ptn, time.Now().Unix(), err.Error())
				sentry.CaptureMessage(msg)
				continue
			}
		}
	}()
}
