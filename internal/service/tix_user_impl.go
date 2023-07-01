package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain/response"
	"net/http"
)

func (service *tixService) GenerateMagicLink(
	ctx context.Context,
	email string,
) *response.ServiceSingleRespond {
	user, err := service.postgreSQLRepository.GetUserByEmail(ctx, email)
	if err != nil && errors.Is(err, sql.ErrNoRows) && user == nil {
		return &response.ServiceSingleRespond{
			Code:    http.StatusInternalServerError,
			Message: common.ErrUserNotFound.Error(),
		}
	}

	data, err := service.authRESTRepository.SendMagicLink(ctx, email)
	if err != nil {
		return &response.ServiceSingleRespond{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return &response.ServiceSingleRespond{
		Code:    data.Code,
		Message: data.Message,
	}
}

func (service *tixService) FetchUsers(
	ctx context.Context,
	email string,
) (
	items []*response.UserResponse,
	err error,
) {
	data, err := service.postgreSQLRepository.GetAllUsers(ctx, email)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		for _, user := range data {
			items = append(items, &response.UserResponse{
				ID:         user.ID,
				UUID:       user.UUID,
				Username:   user.Username,
				Email:      user.Email,
				IsVerified: user.EmailVerifiedAt.Valid,
			})
		}
	}

	return items, nil
}

func (service *tixService) SetUserAsVerified(
	ctx context.Context,
	email string,
) error {
	return service.postgreSQLRepository.UpdateUserVerifiedTime(ctx, email)
}

func (service *tixService) InviteUserByEmail(
	ctx context.Context,
	email string,
) *response.ServiceSingleRespond {
	user, err := service.postgreSQLRepository.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &response.ServiceSingleRespond{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if user != nil {
		return &response.ServiceSingleRespond{
			Code:    http.StatusInternalServerError,
			Message: common.ErrUserAlreadyInvited.Error(),
		}
	}

	data, err := service.authRESTRepository.InviteUserByEmail(ctx, email)
	if err != nil {
		return &response.ServiceSingleRespond{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return &response.ServiceSingleRespond{
		Code:    data.Code,
		Message: data.Message,
	}
}

func (service *tixService) DeleteUser(
	ctx context.Context,
	uuid string,
) error {
	if _, err := service.authRESTRepository.DeleteUser(ctx, uuid); err != nil {
		return err
	}

	return service.postgreSQLRepository.DeleteUser(ctx, uuid)
}
