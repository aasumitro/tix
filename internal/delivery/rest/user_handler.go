package rest

import (
	"context"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/config"
	"github.com/aasumitro/tix/internal/domain"
	"github.com/aasumitro/tix/internal/domain/request"
	"github.com/aasumitro/tix/pkg/http/middleware"
	"github.com/aasumitro/tix/pkg/http/wrapper"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserRESTHandler struct {
	Service domain.ITixService
}

func (handler *UserRESTHandler) Fetch(ctx *gin.Context) {
	email := ctx.Query("email")
	ctxWT, cancel := context.WithTimeout(ctx.Request.Context(), common.ContextTimeout*time.Second)
	defer cancel()
	data, err := handler.Service.FetchUsers(ctxWT, email)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, data)
}

func (handler *UserRESTHandler) Invite(ctx *gin.Context) {
	var body request.AuthRequestInvite
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	ctxWT, cancel := context.WithTimeout(ctx.Request.Context(), common.ContextTimeout*time.Second)
	defer cancel()
	data := handler.Service.InviteUserByEmail(ctxWT, body.Email)
	wrapper.NewHTTPRespondWrapper(ctx, data.Code, data.Message)
}

func (handler *UserRESTHandler) Remove(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == ctx.MustGet("user_uuid").(string) {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusNotAcceptable,
			"You cannot delete your own account!")
		return
	}
	ctxWT, cancel := context.WithTimeout(ctx.Request.Context(), common.ContextTimeout*time.Second)
	defer cancel()
	err := handler.Service.DeleteUser(ctxWT, uuid)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusNoContent, nil)
}

func NewUserRESTHandler(
	router *gin.RouterGroup,
	service domain.ITixService,
) {
	handler := &UserRESTHandler{service}
	router = router.Group("/users")
	router.Use(middleware.Auth(config.Instance.SupabaseJWTSecret))
	router.GET(common.EmptyPath, handler.Fetch)
	router.POST("/invite", handler.Invite)
	router.DELETE("/remove/:uuid", handler.Remove)
}
