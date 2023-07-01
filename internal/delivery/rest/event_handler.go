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
	"strconv"
	"strings"
	"time"
)

type EventRESTHandler struct {
	Service domain.ITixService
}

func (handler *EventRESTHandler) Fetch(ctx *gin.Context) {
	ctxWT, cancel := context.WithTimeout(ctx.Request.Context(), common.ContextTimeout*time.Second)
	defer cancel()
	data, err := handler.Service.FetchEvents(ctxWT)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, data)
}

func (handler *EventRESTHandler) Store(ctx *gin.Context) {
	var body request.EventRequestMakeNew
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	ctxWT, cancel := context.WithTimeout(ctx.Request.Context(), common.ContextTimeout*time.Second)
	defer cancel()
	data, err := handler.Service.StoreEvent(ctxWT, &body)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusCreated, data)
}

func (handler *EventRESTHandler) Validate(ctx *gin.Context) {
	var body request.EventValidationRequest
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	ctxWT, cancel := context.WithTimeout(ctx.Request.Context(), common.ContextTimeout*time.Second)
	defer cancel()
	data, err := handler.Service.FetchForms(ctxWT, body.GoogleFormID)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, data)
}

func (handler *EventRESTHandler) Overview(ctx *gin.Context) {
	id := ctx.Param("google_form_id")
	ctxWT, cancel := context.WithTimeout(
		ctx.Request.Context(),
		common.ContextTimeout*time.Second)
	defer cancel()
	data, err := handler.Service.FetchOverview(ctxWT, id)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, data)
}

func (handler *EventRESTHandler) Participants(ctx *gin.Context) {
	id := ctx.Param("google_form_id")
	ctxWT, cancel := context.WithTimeout(ctx.Request.Context(), common.ContextTimeout*time.Second)
	defer cancel()
	data, err := handler.Service.FetchParticipants(ctxWT, id)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, data)
}

func (handler *EventRESTHandler) Sync(ctx *gin.Context) {
	googleFormID := ctx.Param("google_form_id")
	ctxWT, cancel := context.WithTimeout(
		ctx.Request.Context(),
		common.ContextTimeout*time.Second)
	defer cancel()
	if err := handler.Service.PublishSyncEventDataQueue(
		ctxWT, googleFormID,
	); err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, common.MsgWaitSync)
}

func (handler *EventRESTHandler) Status(ctx *gin.Context) {
	googleFormID := ctx.Param("google_form_id")
	participantID := ctx.Param("participant_id")
	pid, err := strconv.ParseInt(participantID, 10, 32)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var body request.EventRequestUpdateParticipant
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if strings.EqualFold(strings.ToLower(body.Status), string(common.ParticipantRequestDeclined)) && body.DeclinedReason == "" {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity,
			common.ErrDeclineReasonNotProvide.Error())
		return
	}
	ctxWT, cancel := context.WithTimeout(
		ctx.Request.Context(),
		common.ContextTimeout*time.Second)
	defer cancel()
	if err := handler.Service.UpdateParticipantStatus(
		ctxWT, googleFormID, int32(pid), &body,
	); err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, common.MsgWaitGenTix)
}

func (handler *EventRESTHandler) Generate(ctx *gin.Context) {
	googleFormID := ctx.Param("google_form_id")
	participantID := ctx.Param("participant_id")
	pid, err := strconv.ParseInt(participantID, 10, 32)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctxWT, cancel := context.WithTimeout(
		ctx.Request.Context(),
		common.ContextTimeout*time.Second)
	defer cancel()
	if err := handler.Service.PublishGenerateEventTicketQueue(
		ctxWT, googleFormID, int32(pid),
	); err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, common.MsgWaitGenTix)
}

func (handler *EventRESTHandler) Export(ctx *gin.Context) {
	googleFormID := ctx.Param("google_form_id")
	exportType := ctx.Param("export_type")
	email := ctx.MustGet("user_email").(string)
	ctxWT, cancel := context.WithTimeout(
		ctx.Request.Context(),
		common.ContextTimeout*time.Second)
	defer cancel()
	if err := handler.Service.PublishExportEventDataQueue(
		ctxWT, googleFormID, exportType, email,
	); err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, common.MsgWaitExport)
}

func NewEventRESTHandler(
	router *gin.RouterGroup,
	service domain.ITixService,
) {
	handler := &EventRESTHandler{service}
	router = router.Group("/events")
	router.Use(middleware.Auth(config.Instance.SupabaseJWTSecret))
	router.GET(common.EmptyPath, handler.Fetch)
	router.POST(common.EmptyPath, handler.Store)
	router.POST("/validate", handler.Validate)
	router.GET("/:google_form_id/overview", handler.Overview)
	router.GET("/:google_form_id/participants", handler.Participants)
	router.POST("/:google_form_id/sync", handler.Sync)
	router.PATCH("/:google_form_id/participants/:participant_id/status", handler.Status)
	router.POST("/:google_form_id/participants/:participant_id/ticket", handler.Generate)
	router.POST("/:google_form_id/export/:export_type", handler.Export)
}
