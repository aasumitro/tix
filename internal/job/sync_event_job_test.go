package job_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain/response"
	"github.com/aasumitro/tix/internal/job"
	"github.com/aasumitro/tix/mocks"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type tixJobTestSuite struct {
	suite.Suite
}

func (s *tixJobTestSuite) TestEventCronJob_success() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	jsonData, err := json.Marshal([]*response.AutoSyncRespond{
		{
			FormID:    "qwe",
			EventDate: int32(time.Now().Add(-1 * time.Hour).Unix()),
		},
		{
			FormID:    "qwe",
			EventDate: int32(time.Now().Add(1 * time.Hour).Unix()),
		},
	})
	if err != nil {
		return
	}
	redisClient.Set(context.TODO(), common.AutoSyncEventKey, jsonData, 1)
	miniRedis.CheckGet(s.T(), common.AutoSyncEventKey, string(jsonData))
	if !miniRedis.Exists(common.AutoSyncEventKey) {
		s.Error(errors.New("key not exists"))
	}
	tixService.On("SyncRespondData", mock.Anything, mock.Anything).Return(nil).Once()
	tixService.On("SyncRespondData", mock.Anything, mock.Anything).Return(nil).Once()
	job.NewEventJob(tixService, redisClient)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func (s *tixJobTestSuite) TestEventCronJob_ErrorDecodeData() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	jsonData, err := json.Marshal(map[string]interface{}{
		"form_id":    123,
		"event_date": int32(time.Now().Add(-1 * time.Hour).Unix()),
	})
	if err != nil {
		return
	}
	redisClient.Set(context.TODO(), common.AutoSyncEventKey, jsonData, 1)
	miniRedis.CheckGet(s.T(), common.AutoSyncEventKey, string(jsonData))
	if !miniRedis.Exists(common.AutoSyncEventKey) {
		s.Error(errors.New("key not exists"))
	}
	job.NewEventJob(tixService, redisClient)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func (s *tixJobTestSuite) TestEventCronJob_ErrorNoData() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	redisClient.Del(context.Background(), common.AutoSyncEventKey)
	job.NewEventJob(tixService, redisClient)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func (s *tixJobTestSuite) TestEventCronJob_ErrorNilData() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
	redisClient.Set(context.Background(), common.AutoSyncEventKey, nil, 1)
	job.NewEventJob(tixService, redisClient)
	miniRedis.Close()
}

func (s *tixJobTestSuite) TestEventStreamer_SYNC_Success() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	job.NewEventJob(tixService, redisClient)
	jsonData, err := json.Marshal(map[string]string{
		"google_form_id": "asd",
	})
	if err != nil {
		s.Error(err)
	}
	tixService.On("SyncRespondData", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
	if err := redisClient.Publish(context.TODO(), common.ReqSyncEventQueueKey, jsonData).Err(); err != nil {
		s.Error(err)
	}
	time.Sleep(100 * time.Millisecond)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func (s *tixJobTestSuite) TestEventStreamer_SYNC_ErrorService() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	job.NewEventJob(tixService, redisClient)
	jsonData, err := json.Marshal(map[string]string{
		"google_form_id": "asd",
	})
	if err != nil {
		s.Error(err)
	}
	tixService.On("SyncRespondData", mock.Anything, mock.Anything).Return(errors.New("lorem")).Once()
	if err := redisClient.Publish(context.TODO(), common.ReqSyncEventQueueKey, jsonData); err == nil {
		s.Error(err.Err())
	}
	time.Sleep(100 * time.Millisecond)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func (s *tixJobTestSuite) TestEventStreamer_SYNC_ErrorMarshal() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	job.NewEventJob(tixService, redisClient)
	jsonData, err := json.Marshal(1)
	if err != nil {
		return
	}
	redisClient.Publish(context.TODO(), common.ReqSyncEventQueueKey, jsonData)
	time.Sleep(100 * time.Millisecond)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func (s *tixJobTestSuite) TestEventStreamer_GEN_Success() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	job.NewEventJob(tixService, redisClient)
	jsonData, err := json.Marshal(map[string]any{
		"google_form_id": "asd",
		"participant_id": 1,
	})
	if err != nil {
		s.Error(err)
	}
	tixService.On("GenerateTicket", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	if err := redisClient.Publish(context.TODO(), common.ReqGenEventTixQueueKey, jsonData).Err(); err != nil {
		s.Error(err)
	}

	time.Sleep(100 * time.Millisecond)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func (s *tixJobTestSuite) TestEventStreamer_GEN_ErrorService() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	job.NewEventJob(tixService, redisClient)
	jsonData, err := json.Marshal(map[string]any{
		"google_form_id": "asd",
		"participant_id": 1,
	})
	if err != nil {
		s.Error(err)
	}
	tixService.On("GenerateTicket", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("lorem")).Once()
	if err := redisClient.Publish(context.TODO(), common.ReqGenEventTixQueueKey, jsonData); err == nil {
		s.Error(err.Err())
	}
	time.Sleep(100 * time.Millisecond)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func (s *tixJobTestSuite) TestEventStreamer_GEN_ErrorMarshal() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	job.NewEventJob(tixService, redisClient)
	jsonData, err := json.Marshal(1)
	if err != nil {
		return
	}
	redisClient.Publish(context.TODO(), common.ReqGenEventTixQueueKey, jsonData)
	time.Sleep(100 * time.Millisecond)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func (s *tixJobTestSuite) TestEventStreamer_EXP_Success() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	job.NewEventJob(tixService, redisClient)
	jsonData, err := json.Marshal(map[string]string{
		"google_form_id": "asd",
		"export_type":    "pdf",
		"email":          "hello@tix.id",
	})
	if err != nil {
		s.Error(err)
	}
	tixService.On("ExportEvent", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
	if err := redisClient.Publish(context.TODO(), common.ReqExpEventDataQueueKey, jsonData).Err(); err != nil {
		s.Error(err)
	}
	time.Sleep(100 * time.Millisecond)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func (s *tixJobTestSuite) TestEventStreamer_EXP_ErrorService() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	job.NewEventJob(tixService, redisClient)
	jsonData, err := json.Marshal(map[string]string{
		"google_form_id": "asd",
		"export_type":    "pdf",
		"email":          "hello@tix.id",
	})
	if err != nil {
		s.Error(err)
	}
	tixService.On("ExportEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("lorem")).Once()
	if err := redisClient.Publish(context.TODO(), common.ReqExpEventDataQueueKey, jsonData).Err(); err != nil {
		s.Error(err)
	}
	time.Sleep(100 * time.Millisecond)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func (s *tixJobTestSuite) TestEventStreamer_EXP_ErrorMarshal() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	tixService := new(mocks.ITixService)
	job.NewEventJob(tixService, redisClient)
	jsonData, err := json.Marshal(1)
	if err != nil {
		return
	}
	if err := redisClient.Publish(context.TODO(), common.ReqExpEventDataQueueKey, jsonData).Err(); err != nil {
		s.Error(err)
	}
	time.Sleep(100 * time.Millisecond)
	miniRedis.Close()
	if err := redisClient.Close(); err != nil {
		s.Error(err)
	}
}

func TestTixJob(t *testing.T) {
	suite.Run(t, new(tixJobTestSuite))
}
