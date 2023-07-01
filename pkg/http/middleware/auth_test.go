package middleware_test

import (
	"github.com/aasumitro/tix/config"
	"github.com/aasumitro/tix/pkg/http/middleware"
	"github.com/aasumitro/tix/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddleware(t *testing.T) {
	viper.Reset()
	viper.SetConfigFile("../../../.example.env")
	viper.SetConfigType("dotenv")
	config.LoadEnv()
	router := gin.Default()
	router.Use(middleware.Auth(config.Instance.SupabaseJWTSecret))
	t.Run("ERROR COOKIE", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("ERROR PARSE TOKEN WITH CLAIM", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJQT1NCRV9URVNUIiwiZXhwIjotNjIxMzU1OTY4MDAsImlhdCI6LTYyMTM1NTk2ODAwLCJwYXlsb2FkIjp7ImRhdGEiOiJoZWxsbyB3b3JsZCJ9fQ.-_tfeKKhqSRP2H_pVg4f_spkX_Z1Lo1nuiu09OFFvO0"
		req.AddCookie(&http.Cookie{Name: "access_token", Value: accessToken})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("SUCCESS", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		jwt := token.JSONWebToken{
			Issuer:    "MIDDLEWARE_TEST",
			SecretKey: []byte(config.Instance.SupabaseJWTSecret),
			IssuedAt:  time.Now(),
			ExpiredAt: time.Now().Add(1 * time.Minute),
		}
		accessToken, err := jwt.Claim(nil)
		assert.Nil(t, err)
		req.AddCookie(&http.Cookie{Name: "access_token", Value: accessToken})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
