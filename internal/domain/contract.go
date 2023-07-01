package domain

import (
	"context"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain/entity"
	"github.com/aasumitro/tix/internal/domain/request"
	"github.com/aasumitro/tix/internal/domain/response"
	"google.golang.org/api/forms/v1"
)

type (
	IAuthRESTRepository interface {
		SendMagicLink(
			ctx context.Context,
			email string,
		) (resp *response.SupabaseRespond, err error)
		InviteUserByEmail(
			ctx context.Context,
			email string,
		) (resp *response.SupabaseRespond, err error)
		DeleteUser(
			ctx context.Context,
			uuid string,
		) (resp *response.SupabaseRespond, err error)
	}

	IGoogleServiceRepository interface {
		GetEvent(
			ctx context.Context,
			formID string,
		) (*forms.Form, error)
		GetResponses(
			ctx context.Context,
			formID string,
		) (*forms.ListFormResponsesResponse, error)
	}

	IPostgreSQLRepository interface {
		CountUsers(ctx context.Context) int
		GetAllUsers(ctx context.Context, email string) (users []*entity.User, err error)
		GetUserByEmail(ctx context.Context, email string) (user *entity.User, err error)
		UpdateUserVerifiedTime(ctx context.Context, email string) error
		DeleteUser(ctx context.Context, email string) error

		GetAllEvents(ctx context.Context) (events []*entity.Event, err error)
		GetEventByGoogleFormID(ctx context.Context, googleFormID string) (event *entity.Event, err error)
		InsertNewEvent(ctx context.Context, param *request.EventRequestMakeNew) (event *entity.Event, err error)

		CountParticipants(
			ctx context.Context,
			eventID int32,
			participantStatus common.EventParticipantStatus,
			startBetween, endBetween int64,
		) int
		GetAllParticipants(
			ctx context.Context,
			eventID int32, filter string,
			startBetween, endBetween int64, limit int32,
			sortKey, sortDir string,
		) (
			participants []*entity.Participant,
			err error,
		)
		GetParticipantByEmailAndEventID(
			ctx context.Context,
			email string, eventID int32,
		) (
			participant *entity.Participant,
			err error,
		)
		GetParticipantByIDAndEventID(
			ctx context.Context,
			participantID, eventID int32,
		) (
			participant *entity.Participant,
			err error,
		)
		InsertManyParticipants(
			ctx context.Context,
			participants []*entity.Participant,
			createdAt int64,
		) error
		UpdateParticipants(
			ctx context.Context,
			approvedAt, declinedAt *int64,
			declinedReason *string,
			id int32,
		) error
	}

	ITixService interface {
		FetchForms(
			ctx context.Context,
			formID string,
		) (items []*response.GoogleFormQuestion, err error)
		FetchResponds(
			ctx context.Context,
			formID string,
		) (items []*response.GoogleFormRespond, err error)
		SyncRespondData(
			ctx context.Context,
			formID string,
		) error

		FetchEvents(ctx context.Context) (
			items []*response.EventResponse,
			err error,
		)
		StoreEvent(
			ctx context.Context,
			form *request.EventRequestMakeNew,
		) (
			item *response.EventResponse,
			err error,
		)
		FetchOverview(
			ctx context.Context,
			googleFormID string,
		) (item *response.EventOverviewResponse, err error)
		FetchParticipants(
			ctx context.Context,
			googleFormID string,
		) (
			items []*response.ParticipantResponse,
			err error,
		)

		GenerateMagicLink(
			ctx context.Context,
			email string,
		) *response.ServiceSingleRespond
		FetchUsers(
			ctx context.Context,
			email string,
		) (
			items []*response.UserResponse,
			err error,
		)
		SetUserAsVerified(
			ctx context.Context,
			email string,
		) error
		InviteUserByEmail(
			ctx context.Context,
			email string,
		) *response.ServiceSingleRespond
		DeleteUser(
			ctx context.Context,
			uuid string,
		) error
		UpdateParticipantStatus(
			ctx context.Context,
			googleFormID string,
			participantID int32,
			form *request.EventRequestUpdateParticipant,
		) error
		PublishSyncEventDataQueue(
			ctx context.Context,
			googleFormID string,
		) error
		PublishExportEventDataQueue(
			ctx context.Context,
			googleFormID, exportType, email string,
		) error
		PublishGenerateEventTicketQueue(
			ctx context.Context,
			googleFormID string,
			participantID int32,
		) error

		ExportEvent(
			ctx context.Context,
			googleFormID, exportFileType, targetEmail string,
		) error
		GenerateTicket(
			ctx context.Context,
			googleFormID string,
			participantID int32,
		) error
	}
)
