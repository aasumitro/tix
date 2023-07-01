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
	"testing"
)

type userHandlerTestSuite struct {
	suite.Suite
}

func (s *userHandlerTestSuite) SetupSuite() {
	viper.Reset()
	viper.SetConfigFile("../../../.example.env")
	viper.SetConfigType("dotenv")
	config.LoadEnv()

	svcMock := new(mocks.ITixService)
	eg := gin.Default().Group("test/user")
	rest.NewUserRESTHandler(eg, svcMock)
}

func (s *userHandlerTestSuite) Test_Fetch_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("FetchUsers", mock.Anything, mock.Anything).
		Return([]*response.UserResponse{{
			ID:         1,
			UUID:       "123",
			Username:   "lorem",
			Email:      "lorem@ipsum.id",
			IsVerified: true,
		}}, nil).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/users", http.NoBody)
	ctx.Request = req
	handler := rest.UserRESTHandler{Service: svcMock}
	handler.Fetch(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusOK, writer.Code)
	s.Equal(http.StatusOK, got.Code)
	s.Equal(http.StatusText(http.StatusOK), got.Status)
}
func (s *userHandlerTestSuite) Test_Fetch_ShouldError() {
	svcMock := new(mocks.ITixService)
	svcMock.On("FetchUsers", mock.Anything, mock.Anything).
		Return(nil, errors.New("lorem")).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/users", http.NoBody)
	ctx.Request = req
	handler := rest.UserRESTHandler{Service: svcMock}
	handler.Fetch(ctx)
	var got wrapper.ErrorRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusBadRequest, writer.Code)
	s.Equal(http.StatusBadRequest, got.Code)
	s.Equal(http.StatusText(http.StatusBadRequest), got.Status)
}

func (s *userHandlerTestSuite) Test_Invite_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("InviteUserByEmail", mock.Anything, mock.Anything).
		Return(&response.ServiceSingleRespond{
			Code:    http.StatusCreated,
			Message: "User invited",
		}).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
		"email": "hello@tix.id",
	})
	handler := rest.UserRESTHandler{Service: svcMock}
	handler.Invite(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusCreated, writer.Code)
	s.Equal(http.StatusCreated, got.Code)
	s.Equal(http.StatusText(http.StatusCreated), got.Status)
}
func (s *userHandlerTestSuite) Test_Invite_ShouldError() {
	s.T().Run("ERROR ENTITY", func(t *testing.T) {
		svcMock := new(mocks.ITixService)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{})
		handler := rest.UserRESTHandler{Service: svcMock}
		handler.Invite(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusUnprocessableEntity, writer.Code)
		s.Equal(http.StatusUnprocessableEntity, got.Code)
		s.Equal(http.StatusText(http.StatusUnprocessableEntity), got.Status)
	})
	s.T().Run("ERROR SERVICE", func(t *testing.T) {
		svcMock := new(mocks.ITixService)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		svcMock.On("InviteUserByEmail", mock.Anything, mock.Anything).
			Return(&response.ServiceSingleRespond{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}).Once()
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
			"email": "hello@tix.id",
		})
		handler := rest.UserRESTHandler{Service: svcMock}
		handler.Invite(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusInternalServerError, writer.Code)
		s.Equal(http.StatusInternalServerError, got.Code)
		s.Equal(http.StatusText(http.StatusInternalServerError), got.Status)
	})
}

func (s *userHandlerTestSuite) Test_Remove_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	svcMock.On("DeleteUser", mock.Anything, mock.Anything).Return(nil).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	ctx.Params = []gin.Param{{Key: "uuid", Value: "123"}}
	ctx.Set("user_uuid", "456")
	tests.MockJSONRequest(ctx, http.MethodDelete, "application/json", nil)
	handler := rest.UserRESTHandler{Service: svcMock}
	handler.Remove(ctx)
	s.Equal(http.StatusNoContent, writer.Code)
}
func (s *userHandlerTestSuite) Test_Remove_ShouldError() {
	s.T().Run("CANNOT DELETE OWN ACCOUNT", func(t *testing.T) {
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		ctx.Params = []gin.Param{{Key: "uuid", Value: "123"}}
		ctx.Set("user_uuid", "123")
		tests.MockJSONRequest(ctx, http.MethodDelete, "application/json", nil)
		handler := rest.UserRESTHandler{Service: new(mocks.ITixService)}
		handler.Remove(ctx)
		s.Equal(http.StatusNotAcceptable, writer.Code)
	})
	s.T().Run("ERROR SERVICE", func(t *testing.T) {
		svcMock := new(mocks.ITixService)
		svcMock.On("DeleteUser", mock.Anything, mock.Anything).Return(errors.New("lorem")).Once()
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		ctx.Params = []gin.Param{{Key: "uuid", Value: "123"}}
		ctx.Set("user_uuid", "456")
		tests.MockJSONRequest(ctx, http.MethodDelete, "application/json", nil)
		handler := rest.UserRESTHandler{Service: svcMock}
		handler.Remove(ctx)
		s.Equal(http.StatusBadRequest, writer.Code)
	})
}

func TestUserHandlerService(t *testing.T) {
	suite.Run(t, new(userHandlerTestSuite))
}
