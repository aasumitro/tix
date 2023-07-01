package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain/entity"
	"github.com/aasumitro/tix/internal/domain/request"
	"github.com/aasumitro/tix/internal/domain/response"
	"github.com/aasumitro/tix/pkg/dt"
	"github.com/redis/go-redis/v9"
	"strings"
	"sync"
	"time"
)

func (service *tixService) FetchEvents(
	ctx context.Context,
) (
	items []*response.EventResponse,
	err error,
) {
	data, err := service.postgreSQLRepository.GetAllEvents(ctx)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, event := range data {
		wg.Add(1)
		go func(event *entity.Event) {
			defer wg.Done()
			items = append(items, &response.EventResponse{
				ID:                event.ID,
				GoogleFormID:      event.GoogleFormID,
				Name:              event.Name,
				Location:          event.Location,
				PreregisterDate:   event.PreregisterDate,
				EventDate:         event.EventDate,
				TotalParticipants: event.TotalParticipants,
				IsActive:          false,
			})
		}(event)
	}
	wg.Wait()

	cacheData, err := service.redisCache.Get(ctx, common.AutoSyncEventKey).Result()
	if err == nil {
		var cacheItems []*response.AutoSyncRespond
		if cacheData != "" {
			_ = json.Unmarshal([]byte(cacheData), &cacheItems)
		}
		if len(cacheItems) > 0 {
			for _, eventCache := range cacheItems {
				for _, event := range items {
					if eventCache.FormID == event.GoogleFormID {
						event.IsActive = true
						break
					}
				}
			}
		}
	}

	return items, nil
}

func (service *tixService) StoreEvent(
	ctx context.Context,
	form *request.EventRequestMakeNew,
) (
	item *response.EventResponse,
	err error,
) {
	data, err := service.postgreSQLRepository.InsertNewEvent(ctx, form)
	if err != nil {
		return nil, err
	}

	if time.Now().Unix() < int64(data.EventDate) {
		service.updateAutoSyncEvent(ctx, data.GoogleFormID, data.EventDate)
	}

	return &response.EventResponse{
		ID:                data.ID,
		GoogleFormID:      data.GoogleFormID,
		Name:              data.Name,
		Location:          data.Location,
		PreregisterDate:   data.PreregisterDate,
		EventDate:         data.EventDate,
		TotalParticipants: 0,
	}, nil
}

func (service *tixService) FetchOverview(
	ctx context.Context,
	googleFormID string,
) (item *response.EventOverviewResponse, err error) {
	cacheKey := fmt.Sprintf("overview-%s", googleFormID)

	data, err := func() (any, error) {
		if valueCache, errCache := service.redisCache.Get(
			ctx, cacheKey,
		).Result(); errCache == nil {
			return valueCache, nil
		}

		var wg sync.WaitGroup
		errCh := make(chan error)
		now := time.Now()
		limit := 5

		event, err := service.postgreSQLRepository.GetEventByGoogleFormID(ctx, googleFormID)
		if err != nil {
			return nil, err
		}

		data := &response.EventOverviewResponse{
			EventResponse: &response.EventResponse{
				ID:                event.ID,
				GoogleFormID:      event.GoogleFormID,
				Name:              event.Name,
				Location:          event.Location,
				PreregisterDate:   event.PreregisterDate,
				EventDate:         event.EventDate,
				TotalParticipants: event.TotalParticipants,
			},
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			startOfDay, endOfDay := dt.CurrentDayStartToEnd(now)
			participants, err := service.postgreSQLRepository.GetAllParticipants(
				ctx, event.ID, "", startOfDay.Unix(), endOfDay.Unix(), int32(limit), "", "")
			if err != nil {
				errCh <- err
				return
			}
			var respondentsToday []*response.ParticipantResponse
			for _, participant := range participants {
				wg.Add(1)
				go func(participant *entity.Participant) {
					defer wg.Done()
					respondentsToday = append(respondentsToday, &response.ParticipantResponse{
						ID:      participant.ID,
						EventID: participant.EventID,
						Name:    participant.Name,
						Email:   participant.Email,
						Phone:   participant.Phone,
						Job:     participant.Job,
						PoP:     participant.PoP,
						DoB:     participant.DoB,
						ApprovedAt: func() *int32 {
							if participant.ApprovedAt.Valid {
								return &participant.ApprovedAt.Int32
							}
							return nil
						}(),
						DeclinedAt: func() *int32 {
							if participant.DeclinedAt.Valid {
								return &participant.DeclinedAt.Int32
							}
							return nil
						}(),
						DeclinedReason: func() string {
							if participant.DeclinedReason.Valid {
								return participant.DeclinedReason.String
							}
							return ""
						}(),
					})
				}(participant)
			}
			data.LatestRespondents = respondentsToday
		}()

		wg.Add(3)
		go func() {
			defer wg.Done()
			data.TotalApprovedParticipant = service.postgreSQLRepository.CountParticipants(
				ctx, event.ID, common.ParticipantRequestApproved, 0, 0)
		}()
		go func() {
			defer wg.Done()
			data.TotalWaitingApprovalParticipant = service.postgreSQLRepository.CountParticipants(
				ctx, event.ID, common.ParticipantRequestWaiting, 0, 0)
		}()
		go func() {
			defer wg.Done()
			data.TotalDeclinedParticipant = service.postgreSQLRepository.CountParticipants(
				ctx, event.ID, common.ParticipantRequestDeclined, 0, 0)
		}()

		var weeklyOverview []*response.WeeklyOverviewResponse
		for _, week := range dt.WeekDayStartToEnd(now) {
			wg.Add(1)
			go func(week *dt.Weekly) {
				defer wg.Done()
				weeklyOverview = append(weeklyOverview, &response.WeeklyOverviewResponse{
					Name: fmt.Sprintf("%d %s", week.Start.Day(), week.Start.Month().String()),
					Total: service.postgreSQLRepository.CountParticipants(
						ctx, event.ID, common.ParticipantStatusNone, week.Start.Unix(), week.End.Unix()),
				})
			}(week)
		}
		data.WeeklyOverview = weeklyOverview

		go func() {
			wg.Wait()
			close(errCh)
		}()

		select {
		case <-ctx.Done():
			wg.Wait()
			return nil, ctx.Err()
		case err := <-errCh:
			wg.Wait()
			return nil, err
		default:
		}

		if jsonData, err := json.Marshal(data); err == nil {
			service.redisCache.Set(ctx, cacheKey, jsonData, common.EventDataCacheTimeDuration)
		}

		return data, nil
	}()

	if err != nil {
		return nil, err
	}

	if data, ok := data.(*response.EventOverviewResponse); ok {
		item = data
	}

	if data, ok := data.(string); ok {
		var overview *response.EventOverviewResponse
		if err := json.Unmarshal([]byte(data), &overview); err != nil {
			return nil, err
		}
		item = overview
	}

	return item, nil
}

func (service *tixService) updateAutoSyncEvent(
	ctx context.Context,
	googleFormID string,
	eventDate int32,
) {
	// Retrieve data from Redis cache
	cacheData, err := service.redisCache.Get(ctx, common.AutoSyncEventKey).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return
		}
		cacheData = ""
	}

	// Unmarshal cache data into a slice of AutoSyncRespond
	var cacheItems []*response.AutoSyncRespond
	if cacheData != "" {
		err = json.Unmarshal([]byte(cacheData), &cacheItems)
		if err != nil {
			return
		}
	}

	var found bool
	for _, event := range cacheItems {
		if event.FormID == googleFormID {
			found = true
			break
		}
	}

	// Append the new event to the cache items
	if !found {
		cacheItems = append(cacheItems, &response.AutoSyncRespond{
			FormID:    googleFormID,
			EventDate: eventDate,
		})
	}

	// Marshal the updated cache items back into JSON
	jsonData, err := json.Marshal(cacheItems)
	if err != nil {
		return
	}

	// Set the updated cache data in Redis
	if err = service.redisCache.Set(ctx, common.AutoSyncEventKey, jsonData, -1).Err(); err != nil {
		return
	}
}

func (service *tixService) FetchParticipants(
	ctx context.Context,
	googleFormID string,
) (
	items []*response.ParticipantResponse,
	err error,
) {
	cacheKey := fmt.Sprintf("participants-%s", googleFormID)

	data, err := func() (any, error) {
		if valueCache, errCache := service.redisCache.Get(
			ctx, cacheKey,
		).Result(); errCache == nil {
			return valueCache, nil
		}

		event, err := service.postgreSQLRepository.GetEventByGoogleFormID(ctx, googleFormID)
		if err != nil {
			return nil, err
		}

		participants, err := service.postgreSQLRepository.GetAllParticipants(
			ctx, event.ID, "", 0, 0, 0, "", "")
		if err != nil {
			return nil, err
		}

		var allParticipants []*response.ParticipantResponse
		for _, participant := range participants {
			allParticipants = append(allParticipants, &response.ParticipantResponse{
				ID:      participant.ID,
				EventID: participant.EventID,
				Name:    participant.Name,
				Email:   participant.Email,
				Phone:   participant.Phone,
				Job:     participant.Job,
				PoP:     participant.PoP,
				DoB:     participant.DoB,
				ApprovedAt: func() *int32 {
					if participant.ApprovedAt.Valid {
						return &participant.ApprovedAt.Int32
					}
					return nil
				}(),
				DeclinedAt: func() *int32 {
					if participant.DeclinedAt.Valid {
						return &participant.DeclinedAt.Int32
					}
					return nil
				}(),
				DeclinedReason: func() string {
					if participant.DeclinedReason.Valid {
						return participant.DeclinedReason.String
					}
					return ""
				}(),
				Status: func() string {
					if participant.ApprovedAt.Valid {
						return "approved"
					}
					if participant.DeclinedAt.Valid {
						return "declined"
					}
					return "waiting approval"
				}(),
			})
		}

		if jsonData, err := json.Marshal(allParticipants); err == nil {
			service.redisCache.Set(ctx, cacheKey, jsonData, common.EventParticipantCacheTimeDuration)
		}

		return allParticipants, nil
	}()

	if err != nil {
		return nil, err
	}

	if data, ok := data.([]*response.ParticipantResponse); ok {
		items = data
	}

	if data, ok := data.(string); ok {
		var participants []*response.ParticipantResponse
		if err := json.Unmarshal([]byte(data), &participants); err != nil {
			return nil, err
		}
		items = participants
	}

	return items, nil
}

func (service *tixService) PublishSyncEventDataQueue(
	ctx context.Context,
	googleFormID string,
) error {
	cacheKey := fmt.Sprintf("%s-%s",
		common.ReqSyncEventQueueKey, googleFormID)

	cache, err := service.redisCache.Get(ctx, cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if cache != "" {
		return common.ErrRateLimitingPushQueue
	}

	payload, err := json.Marshal(map[string]string{
		"google_form_id": googleFormID,
	})
	if err != nil {
		return err
	}

	if err := service.redisCache.Publish(
		ctx, common.ReqSyncEventQueueKey, payload,
	).Err(); err != nil {
		return err
	}

	return service.redisCache.Set(
		ctx, cacheKey, payload, time.Minute*1,
	).Err()
}

func (service *tixService) UpdateParticipantStatus(
	ctx context.Context,
	googleFormID string,
	participantID int32,
	form *request.EventRequestUpdateParticipant,
) error {
	now := time.Now().Unix()
	if err := service.postgreSQLRepository.UpdateParticipants(
		ctx,
		func() *int64 {
			if strings.EqualFold(strings.ToLower(form.Status), string(common.ParticipantRequestApproved)) {
				return &now
			}
			return nil
		}(),
		func() *int64 {
			if strings.EqualFold(strings.ToLower(form.Status), string(common.ParticipantRequestDeclined)) {
				return &now
			}
			return nil
		}(),
		func() *string {
			if strings.EqualFold(strings.ToLower(form.Status), string(common.ParticipantRequestDeclined)) {
				return &form.DeclinedReason
			}
			return nil
		}(),
		participantID,
	); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("participants-%s", googleFormID)
	service.redisCache.Del(ctx, cacheKey)

	if strings.EqualFold(strings.ToLower(form.Status), string(common.ParticipantRequestDeclined)) {
		return nil
	}

	return service.PublishGenerateEventTicketQueue(
		ctx, googleFormID, participantID,
	)
}

func (service *tixService) PublishExportEventDataQueue(
	ctx context.Context,
	googleFormID, exportType, email string,
) error {
	cacheKey := fmt.Sprintf("%s-%s-%s",
		common.ReqExpEventDataQueueKey, googleFormID, email)

	cache, err := service.redisCache.Get(ctx, cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if cache != "" {
		return common.ErrRateLimitingPushQueue
	}

	payload, err := json.Marshal(map[string]string{
		"google_form_id": googleFormID,
		"export_type":    exportType,
		"email":          email,
	})
	if err != nil {
		return err
	}

	if err := service.redisCache.Publish(
		ctx, common.ReqExpEventDataQueueKey, payload,
	).Err(); err != nil {
		return err
	}

	return service.redisCache.Set(
		ctx, cacheKey, payload, time.Minute*1,
	).Err()
}

func (service *tixService) PublishGenerateEventTicketQueue(
	ctx context.Context,
	googleFormID string,
	participantID int32,
) error {
	cacheKey := fmt.Sprintf("%s-%s-%d",
		common.ReqGenEventTixQueueKey, googleFormID, participantID)
	cache, err := service.redisCache.Get(ctx, cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if cache != "" {
		return common.ErrRateLimitingPushQueue
	}

	payload, err := json.Marshal(map[string]any{
		"google_form_id": googleFormID,
		"participant_id": participantID,
	})
	if err != nil {
		return err
	}

	if err := service.redisCache.Publish(
		ctx, common.ReqGenEventTixQueueKey, payload,
	).Err(); err != nil {
		return err
	}

	return service.redisCache.Set(
		ctx, cacheKey, payload, time.Minute*1,
	).Err()
}
