package rest

import (
	"context"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/config"
	"github.com/aasumitro/tix/internal/domain"
	"github.com/aasumitro/tix/internal/domain/request"
	"github.com/aasumitro/tix/pkg/http/middleware"
	"github.com/aasumitro/tix/pkg/http/wrapper"
	"github.com/aasumitro/tix/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type AccountRESTHandler struct {
	Service domain.ITixService
}

// IN FE CHECK TOKEN/CALL VALIDATE EVERY 3 MINUTES
// EVERY CALL STORE LAST CHECK IN LOCAL STORAGE/SESSION STORAGE

// Validate godoc
// @Schemes
// @Summary Validate User Session
// @Description Get Magic Link Or Validate User Session
// @Tags 	Auth
// @Accept 	mpfd
// @Produce json
// @Param 	email formData string true "user email"
// @Router /api/v1/auth/validate [POST]
func (handler *AccountRESTHandler) Validate(ctx *gin.Context) {
	// VALIDATE AND GET ACCESS TOKEN (JWT) COOKIE
	jwt, err := ctx.Request.Cookie(common.AccessTokenCookieKey)
	if err != nil {
		handler.requestMagicLink(ctx)
		return
	}
	// EXTRACT AND VALIDATE IF ACCESS TOKEN (JWT) FROM COOKIE FOUND.
	if claim, err := token.ExtractAndValidateJWT(
		config.Instance.SupabaseJWTSecret, jwt.Value,
	); err != nil || claim == nil {
		handler.requestMagicLink(ctx)
		return
	}
	// IF THERE'S NO ERROR AND CLAIM NOT NULL THEN GIVE USER TO ACCESS THE FE
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, jwt.Value)
}

func (handler *AccountRESTHandler) requestMagicLink(ctx *gin.Context) {
	var body request.AuthRequestMakeMagicLink
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	ctxWT, cancel := context.WithTimeout(ctx.Request.Context(), common.ContextTimeout*time.Second)
	defer cancel()
	data := handler.Service.GenerateMagicLink(ctxWT, body.Email)
	wrapper.NewHTTPRespondWrapper(ctx, data.Code, data.Message)
}

// Verify godoc
// @Schemes
// @Summary 	Verify User Session
// @Description Verify and Set User Session
// @Tags 		Auth
// @Accept 		mpfd
// @Produce 	json
// @Param 		jwt formData string true "jwt from magic link"
// @Router /api/v1/auth/verify [POST]
func (handler *AccountRESTHandler) Verify(ctx *gin.Context) {
	var body request.AuthRequestMakeSession
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if body.Type == "invite" {
		handler.updateUserData(ctx, body.JWT)
	}
	ctx.SetCookie(common.AccessTokenCookieKey, body.JWT, 0, "/", "", false, true)
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusCreated, body.JWT)
}

func (handler *AccountRESTHandler) updateUserData(ctx *gin.Context, jwt string) {
	claim, err := token.ExtractAndValidateJWT(config.Instance.SupabaseJWTSecret, jwt)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	ctxWT, cancel := context.WithTimeout(ctx.Request.Context(), common.ContextTimeout*time.Second)
	defer cancel()
	if err := handler.Service.SetUserAsVerified(ctxWT, claim.Email); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err.Error())
		return
	}
}

// Profile godoc
// @Schemes
// @Summary 	User Profile
// @Description Get user profile from session
// @Tags 		Auth
// @Accept 		mpfd
// @Produce 	json
// @Router /api/v1/auth/profile [GET]
func (handler *AccountRESTHandler) Profile(ctx *gin.Context) {
	var email, username, uuid string
	if userEmail, ok := ctx.MustGet("user_email").(string); ok {
		email = userEmail
		split := strings.Split(userEmail, "@")
		username = split[0]
	}
	if userUUID, ok := ctx.MustGet("user_uuid").(string); ok {
		uuid = userUUID
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, map[string]any{
		"uuid":     uuid,
		"email":    email,
		"username": username,
	})
}

// SignOut godoc
// @Schemes
// @Summary 	Remove User Session
// @Description Logged User Out and Clear Session
// @Tags 		Auth
// @Accept 		mpfd
// @Produce 	json
// @Router /api/v1/auth/logout [POST]
func (handler *AccountRESTHandler) SignOut(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    common.AccessTokenCookieKey,
		Value:   "",
		MaxAge:  -1,
		Domain:  "http://localhost:8000",
		Path:    "/",
		Expires: time.Now().Add(-time.Hour),
	})
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusUnauthorized, nil)
}

func NewAccountRESTHandler(
	router *gin.RouterGroup,
	service domain.ITixService,
) {
	handler := &AccountRESTHandler{service}
	router = router.Group("/auth")
	router.POST("/validate", handler.Validate)
	router.POST("/verify", handler.Verify)
	router.Use(middleware.Auth(config.Instance.SupabaseJWTSecret)).
		GET("/profile", handler.Profile)
	router.Use(middleware.Auth(config.Instance.SupabaseJWTSecret)).
		POST("/logout", handler.SignOut)
}
