package rest_test

import (
	"encoding/json"
	"errors"
	"github.com/aasumitro/tix/config"
	"github.com/aasumitro/tix/internal/delivery/rest"
	"github.com/aasumitro/tix/internal/domain/response"
	"github.com/aasumitro/tix/mocks"
	"github.com/aasumitro/tix/pkg/http/tests"
	"github.com/aasumitro/tix/pkg/http/wrapper"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type eventHandlerTestSuite struct {
	suite.Suite
}

func (s *eventHandlerTestSuite) SetupSuite() {
	viper.Reset()
	viper.SetConfigFile("../../../.example.env")
	viper.SetConfigType("dotenv")
	config.LoadEnv()

	svcMock := new(mocks.ITixService)
	eg := gin.Default().Group("test/events")
	rest.NewEventRESTHandler(eg, svcMock)
}

func (s *eventHandlerTestSuite) Test_Fetch_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("FetchEvents", mock.Anything, mock.Anything).
		Return([]*response.EventResponse{{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "tix",
			Location:          "jln tix",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
			IsActive:          true,
		}}, nil).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/events", http.NoBody)
	ctx.Request = req
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Fetch(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusOK, writer.Code)
	s.Equal(http.StatusOK, got.Code)
	s.Equal(http.StatusText(http.StatusOK), got.Status)
}
func (s *eventHandlerTestSuite) Test_Fetch_ShouldError() {
	svcMock := new(mocks.ITixService)
	svcMock.On("FetchEvents", mock.Anything, mock.Anything).
		Return(nil, errors.New("lorem")).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/events", http.NoBody)
	ctx.Request = req
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Fetch(ctx)
	var got wrapper.ErrorRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusBadRequest, writer.Code)
	s.Equal(http.StatusBadRequest, got.Code)
	s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
}

func (s *eventHandlerTestSuite) Test_Store_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("StoreEvent", mock.Anything, mock.Anything).
		Return(&response.EventResponse{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "tix",
			Location:          "jln tix",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
			IsActive:          true,
		}, nil).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
		"google_form_id":   "asd",
		"name":             "tix",
		"preregister_date": strconv.FormatInt(time.Now().Unix(), 16),
		"event_date":       strconv.FormatInt(time.Now().Unix(), 16),
		"location":         "jln tix",
	})
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Store(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusCreated, writer.Code)
	s.Equal(http.StatusCreated, got.Code)
	s.Equal(http.StatusText(http.StatusCreated), got.Status)
}
func (s *eventHandlerTestSuite) Test_Store_ShouldError() {
	svcMock := new(mocks.ITixService)
	s.T().Run("ERROR ENTITY", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
		handler := rest.EventRESTHandler{Service: svcMock}
		handler.Store(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusUnprocessableEntity, writer.Code)
		s.Equal(http.StatusUnprocessableEntity, got.Code)
		s.Equal(http.StatusText(http.StatusUnprocessableEntity), got.Status)
	})
	s.T().Run("ERROR SERVICE", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		svcMock.On("StoreEvent", mock.Anything, mock.Anything).
			Return(nil, errors.New("lorem")).Once()
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
			"google_form_id":   "asd",
			"name":             "tix",
			"preregister_date": strconv.FormatInt(time.Now().Unix(), 16),
			"event_date":       strconv.FormatInt(time.Now().Unix(), 16),
			"location":         "jln tix",
		})
		handler := rest.EventRESTHandler{Service: svcMock}
		handler.Store(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusBadRequest, writer.Code)
		s.Equal(http.StatusBadRequest, got.Code)
		s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
	})
}

func (s *eventHandlerTestSuite) Test_Validate_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("FetchForms", mock.Anything, mock.Anything).
		Return([]*response.GoogleFormQuestion{{
			ID:    "1",
			Title: "name",
		}}, nil).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
		"google_form_id": "asd",
	})
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Validate(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusOK, writer.Code)
	s.Equal(http.StatusOK, got.Code)
	s.Equal(http.StatusText(http.StatusOK), got.Status)
}
func (s *eventHandlerTestSuite) Test_Validate_ShouldError() {
	svcMock := new(mocks.ITixService)
	s.T().Run("ERROR ENTITY", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
		handler := rest.EventRESTHandler{Service: svcMock}
		handler.Validate(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusUnprocessableEntity, writer.Code)
		s.Equal(http.StatusUnprocessableEntity, got.Code)
		s.Equal(http.StatusText(http.StatusUnprocessableEntity), got.Status)
	})
	s.T().Run("ERROR SERVICE", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		svcMock.On("FetchForms", mock.Anything, mock.Anything).
			Return(nil, errors.New("lorem")).Once()
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
			"google_form_id": "asd",
		})
		handler := rest.EventRESTHandler{Service: svcMock}
		handler.Validate(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusBadRequest, writer.Code)
		s.Equal(http.StatusBadRequest, got.Code)
		s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
	})
}

func (s *eventHandlerTestSuite) Test_Overview_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("FetchOverview", mock.Anything, mock.Anything).
		Return(&response.EventOverviewResponse{}, nil).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/events/asd/overview", http.NoBody)
	ctx.Request = req
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Overview(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusOK, writer.Code)
	s.Equal(http.StatusOK, got.Code)
	s.Equal(http.StatusText(http.StatusOK), got.Status)
}
func (s *eventHandlerTestSuite) Test_Overview_ShouldError() {
	svcMock := new(mocks.ITixService)
	svcMock.On("FetchOverview", mock.Anything, mock.Anything).
		Return(nil, errors.New("lorem")).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/events/asd/overview", http.NoBody)
	ctx.Request = req
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Overview(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusBadRequest, writer.Code)
	s.Equal(http.StatusBadRequest, got.Code)
	s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
}

func (s *eventHandlerTestSuite) Test_Participant_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("FetchParticipants", mock.Anything, mock.Anything).
		Return([]*response.ParticipantResponse{{}}, nil).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/events/asd/participants", http.NoBody)
	ctx.Request = req
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Participants(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusOK, writer.Code)
	s.Equal(http.StatusOK, got.Code)
	s.Equal(http.StatusText(http.StatusOK), got.Status)
}
func (s *eventHandlerTestSuite) Test_Participant_ShouldError() {
	svcMock := new(mocks.ITixService)
	svcMock.On("FetchParticipants", mock.Anything, mock.Anything).
		Return(nil, errors.New("lorem")).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/events/asd/participants", http.NoBody)
	ctx.Request = req
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Participants(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusBadRequest, writer.Code)
	s.Equal(http.StatusBadRequest, got.Code)
	s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
}

func (s *eventHandlerTestSuite) Test_Sync_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("PublishSyncEventDataQueue", mock.Anything, mock.Anything).
		Return(nil).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Sync(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusOK, writer.Code)
	s.Equal(http.StatusOK, got.Code)
	s.Equal(http.StatusText(http.StatusOK), got.Status)
}
func (s *eventHandlerTestSuite) Test_Sync_ShouldError() {
	svcMock := new(mocks.ITixService)
	svcMock.On("PublishSyncEventDataQueue", mock.Anything, mock.Anything).
		Return(errors.New("lorem")).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Sync(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusBadRequest, writer.Code)
	s.Equal(http.StatusBadRequest, got.Code)
	s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
}

func (s *eventHandlerTestSuite) Test_Status_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("UpdateParticipantStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	ctx.AddParam("participant_id", "1")
	tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
		"status": "approved",
	})
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Status(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusOK, writer.Code)
	s.Equal(http.StatusOK, got.Code)
	s.Equal(http.StatusText(http.StatusOK), got.Status)
}
func (s *eventHandlerTestSuite) Test_Status_ShouldError() {
	svcMock := new(mocks.ITixService)
	s.T().Run("error parse", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		ctx.AddParam("participant_id", "asd")
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
		handler := rest.EventRESTHandler{Service: svcMock}
		handler.Status(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusBadRequest, writer.Code)
		s.Equal(http.StatusBadRequest, got.Code)
		s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
	})
	s.T().Run("error bind", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		ctx.AddParam("participant_id", "1")
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
		handler := rest.EventRESTHandler{Service: svcMock}
		handler.Status(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusUnprocessableEntity, writer.Code)
		s.Equal(http.StatusUnprocessableEntity, got.Code)
		s.Equal(http.StatusText(http.StatusUnprocessableEntity), got.Status)
	})
	s.T().Run("error no decline reason", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		ctx.AddParam("participant_id", "1")
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
			"status": "declined",
		})
		handler := rest.EventRESTHandler{Service: svcMock}
		handler.Status(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusUnprocessableEntity, writer.Code)
		s.Equal(http.StatusUnprocessableEntity, got.Code)
		s.Equal(http.StatusText(http.StatusUnprocessableEntity), got.Status)
	})
	s.T().Run("error service", func(t *testing.T) {
		svcMock.On("UpdateParticipantStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("lorem")).Once()
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		ctx.AddParam("participant_id", "1")
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
			"status": "approved",
		})
		handler := rest.EventRESTHandler{Service: svcMock}
		handler.Status(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusBadRequest, writer.Code)
		s.Equal(http.StatusBadRequest, got.Code)
		s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
	})
}

func (s *eventHandlerTestSuite) Test_Generate_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("PublishGenerateEventTicketQueue", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	ctx.AddParam("participant_id", "1")
	tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Generate(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusOK, writer.Code)
	s.Equal(http.StatusOK, got.Code)
	s.Equal(http.StatusText(http.StatusOK), got.Status)
}
func (s *eventHandlerTestSuite) Test_Generate_ShouldError() {
	svcMock := new(mocks.ITixService)
	s.T().Run("error parse", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		ctx.AddParam("participant_id", "asd")
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
		handler := rest.EventRESTHandler{Service: svcMock}
		handler.Generate(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusBadRequest, writer.Code)
		s.Equal(http.StatusBadRequest, got.Code)
		s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
	})
	s.T().Run("error service", func(t *testing.T) {
		svcMock.On("PublishGenerateEventTicketQueue", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("lorem")).Once()
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		ctx.AddParam("participant_id", "1")
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
		handler := rest.EventRESTHandler{Service: svcMock}
		handler.Generate(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusBadRequest, writer.Code)
		s.Equal(http.StatusBadRequest, got.Code)
		s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
	})
}

func (s *eventHandlerTestSuite) Test_Export_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("PublishExportEventDataQueue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	ctx.Set("user_email", "hello@tix.id")
	tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Export(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusOK, writer.Code)
	s.Equal(http.StatusOK, got.Code)
	s.Equal(http.StatusText(http.StatusOK), got.Status)
}
func (s *eventHandlerTestSuite) Test_Export_ShouldError() {
	svcMock := new(mocks.ITixService)
	svcMock.On("PublishExportEventDataQueue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("lorem")).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	ctx.Set("user_email", "hello@tix.id")
	tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
	handler := rest.EventRESTHandler{Service: svcMock}
	handler.Export(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusBadRequest, writer.Code)
	s.Equal(http.StatusBadRequest, got.Code)
	s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
}

func TestEventHandlerService(t *testing.T) {
	suite.Run(t, new(eventHandlerTestSuite))
}
