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
	"github.com/aasumitro/tix/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type accountHandlerTestSuite struct {
	suite.Suite
}

func (s *accountHandlerTestSuite) SetupSuite() {
	viper.Reset()
	viper.SetConfigFile("../../../.example.env")
	viper.SetConfigType("dotenv")
	config.LoadEnv()

	svcMock := new(mocks.ITixService)
	eg := gin.Default().Group("test/auth")
	rest.NewAccountRESTHandler(eg, svcMock)
}

func (s *accountHandlerTestSuite) Test_Validate_ShouldSuccess() {
	s.T().Run("NO COOKIE", func(t *testing.T) {
		svcMock := new(mocks.ITixService)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		svcMock.On("GenerateMagicLink", mock.Anything, mock.Anything).
			Return(&response.ServiceSingleRespond{
				Code:    201,
				Message: "Created",
			}).Once()
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
			"email": "hello@tix.id",
		})
		handler := rest.AccountRESTHandler{Service: svcMock}
		handler.Validate(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusCreated, writer.Code)
		s.Equal(http.StatusCreated, got.Code)
		s.Equal(http.StatusText(http.StatusCreated), got.Status)
	})
	s.T().Run("COOKIE EXPIRED", func(t *testing.T) {
		svcMock := new(mocks.ITixService)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		ctx.Request.AddCookie(&http.Cookie{Name: "access_token", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhdXRoZW50aWNhdGVkIiwiZXhwIjoxNjg0MTQyODU2LCJzdWIiOiI3ZDc3MjA0MS1mNTI0LTQ2NzAtOTZlOS1hYjY4ODFkMmJiOTkiLCJlbWFpbCI6ImFhc3VtaXRyb0BnbWFpbC5jb20iLCJwaG9uZSI6IiIsImFwcF9tZXRhZGF0YSI6eyJwcm92aWRlciI6ImVtYWlsIiwicHJvdmlkZXJzIjpbImVtYWlsIl19LCJ1c2VyX21ldGFkYXRhIjp7fSwicm9sZSI6ImF1dGhlbnRpY2F0ZWQiLCJhYWwiOiJhYWwxIiwiYW1yIjpbeyJtZXRob2QiOiJvdHAiLCJ0aW1lc3RhbXAiOjE2ODQwNTY0NTZ9XSwic2Vzc2lvbl9pZCI6ImQ4ZmI0ZmVlLWY0NzAtNGQyMS04ZGI2LWNmNWQxYWY3MGI1ZSJ9.uo8I5zZhfIzjAJaV2SDkkdspxu7yaVlZromzbsyRzw0"})
		svcMock.On("GenerateMagicLink", mock.Anything, mock.Anything).
			Return(&response.ServiceSingleRespond{
				Code:    201,
				Message: "Created",
			}).Once()
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
			"email": "hello@tix.id",
		})
		handler := rest.AccountRESTHandler{Service: svcMock}
		handler.Validate(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusCreated, writer.Code)
		s.Equal(http.StatusCreated, got.Code)
		s.Equal(http.StatusText(http.StatusCreated), got.Status)
	})
	s.T().Run("FRESH TOKEN", func(t *testing.T) {
		jwt := token.JSONWebToken{
			Issuer:    "MIDDLEWARE_TEST",
			SecretKey: []byte(config.Instance.SupabaseJWTSecret),
			IssuedAt:  time.Now(),
			ExpiredAt: time.Now().Add(1 * time.Minute),
		}
		accessToken, err := jwt.Claim(nil)
		assert.Nil(s.T(), err)
		svcMock := new(mocks.ITixService)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		ctx.Request.AddCookie(&http.Cookie{Name: "access_token", Value: accessToken})
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
			"email": "hello@tix.id",
		})
		handler := rest.AccountRESTHandler{Service: svcMock}
		handler.Validate(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusOK, writer.Code)
		s.Equal(http.StatusOK, got.Code)
		s.Equal(http.StatusText(http.StatusOK), got.Status)
	})
}
func (s *accountHandlerTestSuite) Test_Validate_ShouldError() {
	s.T().Run("ERROR BODY", func(t *testing.T) {
		svcMock := new(mocks.ITixService)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		svcMock.On("GenerateMagicLink", mock.Anything, mock.Anything).
			Return(&response.ServiceSingleRespond{
				Code:    201,
				Message: "Created",
			}).Once()
		tests.MockJSONRequest(ctx, "POST", "application/json", nil)
		handler := rest.AccountRESTHandler{Service: svcMock}
		handler.Validate(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusUnprocessableEntity, writer.Code)
		s.Equal(http.StatusUnprocessableEntity, got.Code)
		s.Equal(http.StatusText(http.StatusUnprocessableEntity), got.Status)
	})
}

func (s *accountHandlerTestSuite) Test_Verify_ShouldSuccess() {
	svcMock := new(mocks.ITixService)
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	svcMock.On("SetUserAsVerified", mock.Anything, mock.Anything).
		Return(nil).Once()
	jwt := token.JSONWebToken{
		Issuer:    "MIDDLEWARE_TEST",
		SecretKey: []byte(config.Instance.SupabaseJWTSecret),
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(1 * time.Minute),
	}
	accessToken, err := jwt.Claim(nil)
	assert.Nil(s.T(), err)
	tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
		"jwt":  accessToken,
		"type": "invite",
	})
	handler := rest.AccountRESTHandler{Service: svcMock}
	handler.Verify(ctx)
	var got wrapper.CommonRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	s.Equal(http.StatusCreated, writer.Code)
	s.Equal(http.StatusCreated, got.Code)
	s.Equal(http.StatusText(http.StatusCreated), got.Status)
}
func (s *accountHandlerTestSuite) Test_Verify_ShouldError() {
	s.T().Run("ERROR NO BODY", func(t *testing.T) {
		svcMock := new(mocks.ITixService)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		tests.MockJSONRequest(ctx, "POST", "application/json", nil)
		handler := rest.AccountRESTHandler{Service: svcMock}
		handler.Verify(ctx)
		var got wrapper.CommonRespond
		_ = json.Unmarshal(writer.Body.Bytes(), &got)
		s.Equal(http.StatusUnprocessableEntity, writer.Code)
		s.Equal(http.StatusUnprocessableEntity, got.Code)
		s.Equal(http.StatusText(http.StatusUnprocessableEntity), got.Status)
	})
	s.T().Run("ERROR EXTRACT JWT", func(t *testing.T) {
		svcMock := new(mocks.ITixService)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
			"jwt":  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhdXRoZW50aWNhdGVkIiwiZXhwIjoxNjg0MTQyODU2LCJzdWIiOiI3ZDc3MjA0MS1mNTI0LTQ2NzAtOTZlOS1hYjY4ODFkMmJiOTkiLCJlbWFpbCI6ImFhc3VtaXRyb0BnbWFpbC5jb20iLCJwaG9uZSI6IiIsImFwcF9tZXRhZGF0YSI6eyJwcm92aWRlciI6ImVtYWlsIiwicHJvdmlkZXJzIjpbImVtYWlsIl19LCJ1c2VyX21ldGFkYXRhIjp7fSwicm9sZSI6ImF1dGhlbnRpY2F0ZWQiLCJhYWwiOiJhYWwxIiwiYW1yIjpbeyJtZXRob2QiOiJvdHAiLCJ0aW1lc3RhbXAiOjE2ODQwNTY0NTZ9XSwic2Vzc2lvbl9pZCI6ImQ4ZmI0ZmVlLWY0NzAtNGQyMS04ZGI2LWNmNWQxYWY3MGI1ZSJ9.uo8I5zZhfIzjAJaV2SDkkdspxu7yaVlZromzbsyRzw0",
			"type": "invite",
		})
		handler := rest.AccountRESTHandler{Service: svcMock}
		handler.Verify(ctx)
	})
	s.T().Run("ERROR SERVICE", func(t *testing.T) {
		svcMock := new(mocks.ITixService)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{Header: make(http.Header)}
		svcMock.On("SetUserAsVerified", mock.Anything, mock.Anything).
			Return(errors.New("lorem")).Once()
		jwt := token.JSONWebToken{
			Issuer:    "MIDDLEWARE_TEST",
			SecretKey: []byte(config.Instance.SupabaseJWTSecret),
			IssuedAt:  time.Now(),
			ExpiredAt: time.Now().Add(1 * time.Minute),
		}
		accessToken, err := jwt.Claim(nil)
		assert.Nil(s.T(), err)
		tests.MockJSONRequest(ctx, "POST", "application/json", map[string]interface{}{
			"jwt":  accessToken,
			"type": "invite",
		})
		handler := rest.AccountRESTHandler{Service: svcMock}
		handler.Verify(ctx)
	})
}

func (s *accountHandlerTestSuite) Test_Profile() {
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	ctx.Set("user_email", "hello@tix.id")
	ctx.Set("user_uuid", "123")
	req, _ := http.NewRequest("GET", "/api/v1/profile", http.NoBody)
	ctx.Request = req
	handler := rest.AccountRESTHandler{Service: new(mocks.ITixService)}
	handler.Profile(ctx)
	var got wrapper.CommonRespond
	s.NotNil(got)
}

func (s *accountHandlerTestSuite) Test_SignOut() {
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	req, _ := http.NewRequest("POST", "/api/v1/logout", http.NoBody)
	ctx.Request = req
	handler := rest.AccountRESTHandler{Service: new(mocks.ITixService)}
	handler.SignOut(ctx)
	var got wrapper.ErrorRespond
	s.NotNil(got)
}

func TestAccountHandlerService(t *testing.T) {
	suite.Run(t, new(accountHandlerTestSuite))
}
