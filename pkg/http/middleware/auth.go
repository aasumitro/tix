package middleware

import (
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func Auth(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string

		if cookie, err := ctx.Request.Cookie(common.AccessTokenCookieKey); err == nil {
			accessToken = cookie.Value
		}

		if authHeader := ctx.Request.Header.Get("Authorization"); authHeader != "" {
			header := strings.Split(authHeader, " ")
			accessToken = header[1]
		}

		if accessToken == "" {
			http.SetCookie(ctx.Writer, &http.Cookie{
				Name:    common.AccessTokenCookieKey,
				Value:   "",
				MaxAge:  -1,
				Domain:  "http://localhost:8000",
				Path:    "/",
				Expires: time.Now().Add(-time.Hour),
			})
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "ACCESS_TOKEN_NOT_PROVIDE")
			return
		}

		claim, err := token.ExtractAndValidateJWT(secret, accessToken)
		if err != nil {
			http.SetCookie(ctx.Writer, &http.Cookie{
				Name:    common.AccessTokenCookieKey,
				Value:   "",
				MaxAge:  -1,
				Domain:  "http://localhost:8000",
				Path:    "/",
				Expires: time.Now().Add(-time.Hour),
			})
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		ctx.Set("user_uuid", claim.Subject)
		ctx.Set("user_email", claim.Email)
		ctx.Set("user_session_id", claim.SessionID)
		ctx.Next()
	}
}
