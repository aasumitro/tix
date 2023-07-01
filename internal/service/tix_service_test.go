package service_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain/entity"
	"github.com/aasumitro/tix/internal/domain/request"
	"github.com/aasumitro/tix/internal/domain/response"
	"github.com/aasumitro/tix/internal/service"
	"github.com/aasumitro/tix/mocks"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gomail.v2"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

type tixServiceTestSuite struct {
	suite.Suite
}

// TIX USER IMPL
func (s *tixServiceTestSuite) Test_GenerateMagicLink_ShouldSuccess() {
	restRepo := new(mocks.IAuthRESTRepository)
	sqlRepo := new(mocks.IPostgreSQLRepository)
	sqlRepo.On("GetUserByEmail", mock.Anything, mock.Anything).
		Return(&entity.User{}, nil).Once()
	restRepo.On("SendMagicLink", mock.Anything, mock.Anything).
		Return(&response.SupabaseRespond{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		}, nil).Once()
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(sqlRepo),
		service.WithAuthRESTRepository(restRepo))
	resp := svc.GenerateMagicLink(context.TODO(), "hello@tix.id")
	s.NotNil(resp)
	s.Equal(resp.Code, http.StatusOK)
	s.Equal(resp.Message, http.StatusText(http.StatusOK))
	restRepo.AssertExpectations(s.T())
	sqlRepo.AssertExpectations(s.T())
}
func (s *tixServiceTestSuite) Test_GenerateMagicLink_ShouldError() {
	restRepo := new(mocks.IAuthRESTRepository)
	sqlRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(sqlRepo),
		service.WithAuthRESTRepository(restRepo))
	s.T().Run("error user not found", func(t *testing.T) {
		sqlRepo.On("GetUserByEmail", mock.Anything, mock.Anything).
			Return(nil, sql.ErrNoRows).Once()
		resp := svc.GenerateMagicLink(context.TODO(), "hello@tix.id")
		s.NotNil(resp)
		s.Equal(resp.Code, http.StatusInternalServerError)
		s.Equal(resp.Message, common.ErrUserNotFound.Error())
		restRepo.AssertExpectations(s.T())
		sqlRepo.AssertExpectations(s.T())
	})
	s.T().Run("error send magic link", func(t *testing.T) {
		sqlRepo.On("GetUserByEmail", mock.Anything, mock.Anything).
			Return(&entity.User{}, nil).Once()
		restRepo.On("SendMagicLink", mock.Anything, mock.Anything).
			Return(nil, errors.New("lorem")).Once()
		resp := svc.GenerateMagicLink(context.TODO(), "hello@tix.id")
		s.NotNil(resp)
		s.Equal(resp.Code, http.StatusInternalServerError)
		s.Equal(resp.Message, "lorem")
		restRepo.AssertExpectations(s.T())
		sqlRepo.AssertExpectations(s.T())
	})
}

func (s *tixServiceTestSuite) Test_FetchUsers_ShouldSuccess() {
	repo := new(mocks.IPostgreSQLRepository)
	repo.On("GetAllUsers", mock.Anything, mock.Anything).
		Return([]*entity.User{{
			ID:       1,
			UUID:     "123",
			Username: "hello",
			Email:    "hello@tix.id",
			EmailVerifiedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}}, nil).Once()
	svc := service.NewTixService(service.WithPostgreSQLRepository(repo))
	resp, err := svc.FetchUsers(context.TODO(), "hello@tix.id")
	s.NotNil(resp)
	s.Nil(err)
	s.Equal(len(resp), 1)
	repo.AssertExpectations(s.T())
}
func (s *tixServiceTestSuite) Test_FetchUsers_ShouldError() {
	repo := new(mocks.IPostgreSQLRepository)
	repo.On("GetAllUsers", mock.Anything, mock.Anything).
		Return(nil, errors.New("lorem")).Once()
	svc := service.NewTixService(service.WithPostgreSQLRepository(repo))
	resp, err := svc.FetchUsers(context.TODO(), "hello@tix.id")
	s.Nil(resp)
	s.NotNil(err)
	s.Equal(err.Error(), "lorem")
	repo.AssertExpectations(s.T())
}

func (s *tixServiceTestSuite) Test_SetUserAsVerified_ShouldSuccess() {
	repo := new(mocks.IPostgreSQLRepository)
	repo.On("UpdateUserVerifiedTime", mock.Anything, mock.Anything).
		Return(nil).Once()
	svc := service.NewTixService(service.WithPostgreSQLRepository(repo))
	err := svc.SetUserAsVerified(context.TODO(), "hello@tix.id")
	s.Nil(err)
	repo.AssertExpectations(s.T())
}
func (s *tixServiceTestSuite) Test_SetUserAsVerified_ShouldError() {
	repo := new(mocks.IPostgreSQLRepository)
	repo.On("UpdateUserVerifiedTime", mock.Anything, mock.Anything).
		Return(errors.New("lorem")).Once()
	svc := service.NewTixService(service.WithPostgreSQLRepository(repo))
	err := svc.SetUserAsVerified(context.TODO(), "hello@tix.id")
	s.NotNil(err)
	repo.AssertExpectations(s.T())
}

func (s *tixServiceTestSuite) Test_InviteUserByEmail_ShouldSuccess() {
	restRepo := new(mocks.IAuthRESTRepository)
	sqlRepo := new(mocks.IPostgreSQLRepository)
	sqlRepo.On("GetUserByEmail", mock.Anything, mock.Anything).
		Return(nil, sql.ErrNoRows).Once()
	restRepo.On("InviteUserByEmail", mock.Anything, mock.Anything).
		Return(&response.SupabaseRespond{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		}, nil).Once()
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(sqlRepo),
		service.WithAuthRESTRepository(restRepo))
	resp := svc.InviteUserByEmail(context.TODO(), "hello@tix.id")
	s.NotNil(resp)
	s.Equal(resp.Code, http.StatusOK)
	s.Equal(resp.Message, http.StatusText(http.StatusOK))
	restRepo.AssertExpectations(s.T())
	sqlRepo.AssertExpectations(s.T())
}
func (s *tixServiceTestSuite) Test_InviteUserByEmail_ShouldError() {
	restRepo := new(mocks.IAuthRESTRepository)
	sqlRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(sqlRepo),
		service.WithAuthRESTRepository(restRepo))
	s.T().Run("error get user data", func(t *testing.T) {
		sqlRepo.On("GetUserByEmail", mock.Anything, mock.Anything).
			Return(nil, errors.New("lorem")).Once()
		resp := svc.InviteUserByEmail(context.TODO(), "hello@tix.id")
		s.NotNil(resp)
		s.Equal(resp.Code, http.StatusInternalServerError)
		s.Equal(resp.Message, "lorem")
		restRepo.AssertExpectations(s.T())
		sqlRepo.AssertExpectations(s.T())
	})
	s.T().Run("error user already invited", func(t *testing.T) {
		sqlRepo.On("GetUserByEmail", mock.Anything, mock.Anything).
			Return(&entity.User{}, nil).Once()
		resp := svc.InviteUserByEmail(context.TODO(), "hello@tix.id")
		s.NotNil(resp)
		s.Equal(resp.Code, http.StatusInternalServerError)
		s.Equal(resp.Message, common.ErrUserAlreadyInvited.Error())
		restRepo.AssertExpectations(s.T())
		sqlRepo.AssertExpectations(s.T())
	})
	s.T().Run("error invite user", func(t *testing.T) {
		sqlRepo.On("GetUserByEmail", mock.Anything, mock.Anything).
			Return(nil, nil).Once()
		restRepo.On("InviteUserByEmail", mock.Anything, mock.Anything).
			Return(nil, errors.New("lorem")).Once()
		resp := svc.InviteUserByEmail(context.TODO(), "hello@tix.id")
		s.NotNil(resp)
		s.Equal(resp.Code, http.StatusInternalServerError)
		s.Equal(resp.Message, "lorem")
		restRepo.AssertExpectations(s.T())
		sqlRepo.AssertExpectations(s.T())
	})
}

func (s *tixServiceTestSuite) Test_DeleteUser_ShouldSuccess() {
	restRepo := new(mocks.IAuthRESTRepository)
	pqRepo := new(mocks.IPostgreSQLRepository)
	restRepo.On("DeleteUser", mock.Anything, mock.Anything).
		Return(&response.SupabaseRespond{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		}, nil).Once()
	pqRepo.On("DeleteUser", mock.Anything, mock.Anything).
		Return(nil).Once()
	svc := service.NewTixService(
		service.WithAuthRESTRepository(restRepo),
		service.WithPostgreSQLRepository(pqRepo))
	err := svc.DeleteUser(context.TODO(), "12345")
	s.Nil(err)
	restRepo.AssertExpectations(s.T())
	pqRepo.AssertExpectations(s.T())
}
func (s *tixServiceTestSuite) Test_DeleteUser_ShouldError() {
	restRepo := new(mocks.IAuthRESTRepository)
	pqRepo := new(mocks.IPostgreSQLRepository)
	s.T().Run("error from rest", func(t *testing.T) {
		restRepo.On("DeleteUser", mock.Anything, mock.Anything).
			Return(nil, errors.New("lorem")).Once()
		svc := service.NewTixService(
			service.WithAuthRESTRepository(restRepo),
			service.WithPostgreSQLRepository(pqRepo))
		err := svc.DeleteUser(context.TODO(), "12345")
		s.NotNil(err)
		restRepo.AssertExpectations(s.T())
		pqRepo.AssertExpectations(s.T())
	})
	s.T().Run("error from postgre", func(t *testing.T) {
		restRepo.On("DeleteUser", mock.Anything, mock.Anything).
			Return(&response.SupabaseRespond{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
			}, nil).Once()
		pqRepo.On("DeleteUser", mock.Anything, mock.Anything).
			Return(errors.New("lorem")).Once()
		svc := service.NewTixService(
			service.WithAuthRESTRepository(restRepo),
			service.WithPostgreSQLRepository(pqRepo))
		err := svc.DeleteUser(context.TODO(), "12345")
		s.NotNil(err)
		restRepo.AssertExpectations(s.T())
		pqRepo.AssertExpectations(s.T())
	})
}

// TIX GOOGLE FORM IMPL
func (s *tixServiceTestSuite) Test_FetchForms_ShouldSuccess() {}
func (s *tixServiceTestSuite) Test_FetchForms_ShouldError()   {}

func (s *tixServiceTestSuite) Test_FetchResponds_ShouldSuccess() {}
func (s *tixServiceTestSuite) Test_FetchResponds_ShouldError()   {}

func (s *tixServiceTestSuite) Test_SyncRespondData_ShouldSuccess() {}
func (s *tixServiceTestSuite) Test_SyncRespondData_ShouldError()   {}

// TIX EXPORT IMPL
func (s *tixServiceTestSuite) Test_ExportEvent_ShouldSuccess() {
	pqRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(pqRepo),
		service.WithMailer(&gomail.Dialer{}))

	s.T().Run("success excel", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(&entity.Event{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "asd",
			Location:          "asd",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}, nil).Once()
		pqRepo.On("GetAllParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*entity.Participant{{
			ID:      1,
			EventID: 1,
			Name:    "asd",
			Email:   "asd@asd.id",
			Phone:   "123",
			Job:     "asd",
			PoP:     "asd",
			DoB:     "asd",
			ApprovedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			DeclinedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			DeclinedReason: sql.NullString{
				String: "asd",
				Valid:  true,
			},
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}}, nil).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()

		err := svc.ExportEvent(context.TODO(), "asd", string(common.ExportTypeXLS), "asd")
		s.Nil(err)
		pqRepo.AssertExpectations(t)
	})

	s.T().Run("success pdf", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(&entity.Event{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "asd",
			Location:          "asd",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}, nil).Once()
		pqRepo.On("GetAllParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*entity.Participant{{
			ID:      1,
			EventID: 1,
			Name:    "asd",
			Email:   "asd@asd.id",
			Phone:   "123",
			Job:     "asd",
			PoP:     "asd",
			DoB:     "asd",
			ApprovedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			DeclinedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			DeclinedReason: sql.NullString{
				String: "asd",
				Valid:  true,
			},
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}}, nil).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		dir := "./temps/exports/"
		filename := "asd.pdf"
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			s.T().Fatalf("Failed to create directory: %s", err)
		}
		filePath := filepath.Join(dir, filename)
		file, err := os.Create(filePath)
		if err != nil {
			s.T().Fatalf("Failed to create file: %s", err)
		}
		defer func() { _ = file.Close() }()
		errSvc := svc.ExportEvent(context.TODO(), "asd", string(common.ExportTypePDF), "asd")
		s.Nil(errSvc)
		if err = os.RemoveAll("./temps"); err != nil {
			s.T().Fatalf("Failed to remove directory: %s", err)
		}
		pqRepo.AssertExpectations(t)
	})
}
func (s *tixServiceTestSuite) Test_ExportEvent_ShouldError() {
	pqRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(pqRepo))
	s.T().Run("error get event", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(nil, errors.New("lorem")).Once()
		err := svc.ExportEvent(context.TODO(), "asd", string(common.ExportTypePDF), "asd")
		s.NotNil(err)
		pqRepo.AssertExpectations(t)
	})
	s.T().Run("error get participant", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(&entity.Event{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "asd",
			Location:          "asd",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}, nil).Once()
		pqRepo.On("GetAllParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("lorem")).Once()
		err := svc.ExportEvent(context.TODO(), "asd", string(common.ExportTypePDF), "asd")
		s.NotNil(err)
		pqRepo.AssertExpectations(t)
	})
}

// TIX GENERATE IMPL
func (s *tixServiceTestSuite) Test_GenerateTicket_ShouldSuccess() {
	pqRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(pqRepo),
		service.WithMailer(&gomail.Dialer{}))
	pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(&entity.Event{
		ID:                1,
		GoogleFormID:      "asd",
		Name:              "asd",
		Location:          "asd",
		PreregisterDate:   int32(time.Now().Unix()),
		EventDate:         int32(time.Now().Unix()),
		TotalParticipants: 1,
		CreatedAt: sql.NullInt32{
			Int32: int32(time.Now().Unix()),
			Valid: true,
		},
		UpdatedAt: sql.NullInt32{
			Int32: int32(time.Now().Unix()),
			Valid: true,
		},
	}, nil).Once()
	pqRepo.On("GetParticipantByIDAndEventID", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Participant{
		ID:    1,
		Name:  "lorem",
		Email: "lorem@lorem.id",
	}, nil).Once()
	dir := "./temps/exports/"
	filename := "gen11tix.pdf"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		s.T().Fatalf("Failed to create directory: %s", err)
	}
	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		s.T().Fatalf("Failed to create file: %s", err)
	}
	defer func() { _ = file.Close() }()
	errSvc := svc.GenerateTicket(context.TODO(), "asd", 1)
	s.Nil(errSvc)
	if err = os.RemoveAll("./temps"); err != nil {
		s.T().Fatalf("Failed to remove directory: %s", err)
	}
	pqRepo.AssertExpectations(s.T())
}
func (s *tixServiceTestSuite) Test_GenerateTicket_ShouldError() {
	pqRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(pqRepo),
		service.WithMailer(&gomail.Dialer{}))
	s.T().Run("error get event", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(nil, errors.New("lorem")).Once()
		errSvc := svc.GenerateTicket(context.TODO(), "asd", 1)
		s.NotNil(errSvc)
		pqRepo.AssertExpectations(s.T())
	})
	s.T().Run("error get participant", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(&entity.Event{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "asd",
			Location:          "asd",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}, nil).Once()
		pqRepo.On("GetParticipantByIDAndEventID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("lorem")).Once()
		errSvc := svc.GenerateTicket(context.TODO(), "asd", 1)
		s.NotNil(errSvc)
		pqRepo.AssertExpectations(s.T())
	})
	s.T().Run("error generate attachment", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(&entity.Event{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "asd",
			Location:          "asd",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}, nil).Once()
		pqRepo.On("GetParticipantByIDAndEventID", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Participant{
			ID:    1,
			Name:  "lorem",
			Email: "lorem@lorem.id",
		}, nil).Once()
		errSvc := svc.GenerateTicket(context.TODO(), "asd", 1)
		s.NotNil(errSvc)
		pqRepo.AssertExpectations(s.T())
	})
}

// TIX EVENT IMPL
func (s *tixServiceTestSuite) Test_FetchEvents_ShouldSuccess() {
	rc := redis.NewClient(&redis.Options{
		Addr: miniredis.RunT(s.T()).Addr(),
	})
	repo := new(mocks.IPostgreSQLRepository)
	repo.On("GetAllEvents", mock.Anything).
		Return([]*entity.Event{{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "tix",
			Location:          "jln tix",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
		}}, nil).Once()
	rc.Set(context.TODO(), common.AutoSyncEventKey, `[{"form_id":"asd","event_date":1690819200}]`, 1)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(repo),
		service.WithRedisCache(rc))
	data, err := svc.FetchEvents(context.TODO())
	s.NotNil(data)
	s.Nil(err)
	repo.AssertExpectations(s.T())
}
func (s *tixServiceTestSuite) Test_FetchEvents_ShouldError() {
	repo := new(mocks.IPostgreSQLRepository)
	repo.On("GetAllEvents", mock.Anything).
		Return(nil, errors.New("lorem")).Once()
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(repo))
	data, err := svc.FetchEvents(context.TODO())
	s.Nil(data)
	s.NotNil(err)
	repo.AssertExpectations(s.T())
}

func (s *tixServiceTestSuite) Test_StoreEvent_ShouldSuccess() {
	repo := new(mocks.IPostgreSQLRepository)
	rc := redis.NewClient(&redis.Options{
		Addr: miniredis.RunT(s.T()).Addr(),
	})
	repo.On("InsertNewEvent", mock.Anything, mock.Anything).
		Return(&entity.Event{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "tix",
			Location:          "jln tix",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Add(1 * time.Hour).Unix()),
			TotalParticipants: 1,
		}, nil).Once()
	rc.Set(context.TODO(), common.AutoSyncEventKey, `[{"form_id":"qwe","event_date":1690819200}]`, 1)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(repo),
		service.WithRedisCache(rc))
	data, err := svc.StoreEvent(context.TODO(), &request.EventRequestMakeNew{
		GoogleFormID:    "asd",
		Name:            "tix",
		PreregisterDate: strconv.FormatInt(time.Now().Unix(), 16),
		EventDate:       strconv.FormatInt(time.Now().Unix(), 16),
		Location:        "jln tix",
	})
	s.NotNil(data)
	s.Nil(err)
	repo.AssertExpectations(s.T())
}
func (s *tixServiceTestSuite) Test_StoreEvent_ShouldError() {
	repo := new(mocks.IPostgreSQLRepository)
	repo.On("InsertNewEvent", mock.Anything, mock.Anything).
		Return(nil, errors.New("lorem")).Once()
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(repo))
	data, err := svc.StoreEvent(context.TODO(), &request.EventRequestMakeNew{})
	s.NotNil(err)
	s.Nil(data)
	repo.AssertExpectations(s.T())
}

func (s *tixServiceTestSuite) Test_FetchOverview_ShouldSuccess() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	pqRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(pqRepo),
		service.WithRedisCache(redisClient))
	s.T().Run("from db", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(&entity.Event{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "asd",
			Location:          "asd",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}, nil).Once()
		pqRepo.On("GetAllParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*entity.Participant{{
			ID:      1,
			EventID: 1,
			Name:    "asd",
			Email:   "asd@asd.id",
			Phone:   "123",
			Job:     "asd",
			PoP:     "asd",
			DoB:     "asd",
			ApprovedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			DeclinedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			DeclinedReason: sql.NullString{
				String: "asd",
				Valid:  true,
			},
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}}, nil).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		data, err := svc.FetchOverview(context.TODO(), "asd")
		s.NotNil(data)
		s.Nil(err)
		time.Sleep(500 * time.Millisecond)
		pqRepo.AssertExpectations(t)
	})
	s.T().Run("from mem", func(t *testing.T) {
		participants := &response.EventOverviewResponse{
			EventResponse: &response.EventResponse{
				ID:                1,
				GoogleFormID:      "asd",
				Name:              "asd",
				Location:          "asd",
				PreregisterDate:   int32(time.Now().Unix()),
				EventDate:         int32(time.Now().Unix()),
				TotalParticipants: 1,
			},
		}
		cacheKey := fmt.Sprintf("overview-%s", "asd")
		if jsonData, err := json.Marshal(participants); err == nil {
			redisClient.Set(context.TODO(), cacheKey, jsonData, 1)
		}
		data, err := svc.FetchOverview(context.TODO(), "asd")
		s.NotNil(data)
		s.Nil(err)
		pqRepo.AssertExpectations(t)
	})
}
func (s *tixServiceTestSuite) Test_FetchOverview_ShouldError() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	pqRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(pqRepo),
		service.WithRedisCache(redisClient))
	s.T().Run("error get event from db", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(nil, errors.New("lorem")).Once()
		data, err := svc.FetchOverview(context.TODO(), "asd")
		s.Nil(data)
		s.NotNil(err)
		pqRepo.AssertExpectations(t)
	})
	s.T().Run("error get participant from db", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(&entity.Event{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "asd",
			Location:          "asd",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}, nil).Once()
		pqRepo.On("GetAllParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("lorem")).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		pqRepo.On("CountParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1).Once()
		data, err := svc.FetchOverview(context.TODO(), "asd")
		s.NotNil(data)
		s.Nil(err)
		time.Sleep(500 * time.Millisecond)
		pqRepo.AssertExpectations(t)
	})
}

func (s *tixServiceTestSuite) Test_FetchParticipants_ShouldSuccess() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	pqRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(pqRepo),
		service.WithRedisCache(redisClient))
	s.T().Run("from db", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(&entity.Event{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "asd",
			Location:          "asd",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}, nil).Once()
		pqRepo.On("GetAllParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*entity.Participant{{
			ID:      1,
			EventID: 1,
			Name:    "asd",
			Email:   "asd@asd.id",
			Phone:   "123",
			Job:     "asd",
			PoP:     "asd",
			DoB:     "asd",
			ApprovedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			DeclinedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			DeclinedReason: sql.NullString{
				String: "asd",
				Valid:  true,
			},
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}}, nil).Once()
		data, err := svc.FetchParticipants(context.TODO(), "asd")
		s.Nil(err)
		s.NotNil(data)
		pqRepo.AssertExpectations(s.T())
	})
	s.T().Run("from redis", func(t *testing.T) {
		participants := []*entity.Participant{{
			ID:      1,
			EventID: 1,
			Name:    "asd",
			Email:   "asd@asd.id",
			Phone:   "123",
			Job:     "asd",
			PoP:     "asd",
			DoB:     "asd",
			ApprovedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			DeclinedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			DeclinedReason: sql.NullString{
				String: "asd",
				Valid:  true,
			},
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}}
		cacheKey := fmt.Sprintf("participants-%s", "asd")
		if jsonData, err := json.Marshal(participants); err == nil {
			redisClient.Set(context.TODO(), cacheKey, jsonData, common.EventParticipantCacheTimeDuration)
		}
		data, err := svc.FetchParticipants(context.TODO(), "asd")
		s.Nil(err)
		s.NotNil(data)
		pqRepo.AssertExpectations(s.T())
	})
}
func (s *tixServiceTestSuite) Test_FetchParticipants_ShouldError() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	pqRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(pqRepo),
		service.WithRedisCache(redisClient))
	s.T().Run("error get event", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(nil, errors.New("lorem")).Once()
		data, err := svc.FetchParticipants(context.TODO(), "asd")
		s.NotNil(err)
		s.Nil(data)
		pqRepo.AssertExpectations(s.T())
	})
	s.T().Run("error get participant", func(t *testing.T) {
		pqRepo.On("GetEventByGoogleFormID", mock.Anything, mock.Anything).Return(&entity.Event{
			ID:                1,
			GoogleFormID:      "asd",
			Name:              "asd",
			Location:          "asd",
			PreregisterDate:   int32(time.Now().Unix()),
			EventDate:         int32(time.Now().Unix()),
			TotalParticipants: 1,
			CreatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
			UpdatedAt: sql.NullInt32{
				Int32: int32(time.Now().Unix()),
				Valid: true,
			},
		}, nil).Once()
		pqRepo.On("GetAllParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("lorem")).Once()
		data, err := svc.FetchParticipants(context.TODO(), "asd")
		s.NotNil(err)
		s.Nil(data)
		pqRepo.AssertExpectations(s.T())
	})
	s.T().Run("error decode", func(t *testing.T) {
		cacheKey := fmt.Sprintf("participants-%s", "asd")
		if jsonData, err := json.Marshal(map[string]string{
			"asd": "asd",
		}); err == nil {
			redisClient.Set(context.TODO(), cacheKey, jsonData, common.EventParticipantCacheTimeDuration)
		}
		data, err := svc.FetchParticipants(context.TODO(), "asd")
		s.NotNil(err)
		s.Nil(data)
		pqRepo.AssertExpectations(s.T())
	})
}

func (s *tixServiceTestSuite) Test_PublishSyncEventDataQueue_ShouldSuccess() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	svc := service.NewTixService(
		service.WithRedisCache(redisClient))
	err := svc.PublishSyncEventDataQueue(context.TODO(), "asd")
	s.Nil(err)
}
func (s *tixServiceTestSuite) Test_PublishSyncEventDataQueue_ShouldError() {
	s.T().Run("error get", func(t *testing.T) {
		miniRedis := miniredis.RunT(s.T())
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		miniRedis.Close()
		_ = redisClient.Close()
		svc := service.NewTixService(service.WithRedisCache(redisClient))
		err := svc.PublishSyncEventDataQueue(context.TODO(), "asd")
		s.NotNil(err)
	})
	s.T().Run("error rate limiter", func(t *testing.T) {
		miniRedis := miniredis.RunT(s.T())
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		cacheKey := fmt.Sprintf("%s-%s",
			common.ReqSyncEventQueueKey, "asd")
		redisClient.Set(context.TODO(), cacheKey, "asd", 1)
		svc := service.NewTixService(service.WithRedisCache(redisClient))
		err := svc.PublishSyncEventDataQueue(context.TODO(), "asd")
		s.NotNil(err)
		s.Equal(err, common.ErrRateLimitingPushQueue)
		redisClient.Del(context.TODO(), cacheKey)
		miniRedis.Close()
		_ = redisClient.Close()
	})
}

func (s *tixServiceTestSuite) Test_UpdateParticipantStatus_ShouldSuccess() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	pqRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(pqRepo),
		service.WithRedisCache(redisClient))
	s.T().Run("approved", func(t *testing.T) {
		pqRepo.On("UpdateParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		err := svc.UpdateParticipantStatus(context.TODO(), "asd", 1, &request.EventRequestUpdateParticipant{
			Status: string(common.ParticipantRequestApproved),
		})
		s.Nil(err)
	})
	s.T().Run("declined", func(t *testing.T) {
		pqRepo.On("UpdateParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		err := svc.UpdateParticipantStatus(context.TODO(), "asd", 1, &request.EventRequestUpdateParticipant{
			Status:         string(common.ParticipantRequestDeclined),
			DeclinedReason: "lorem",
		})
		s.Nil(err)
	})
}
func (s *tixServiceTestSuite) Test_UpdateParticipantStatus_ShouldError() {
	pqRepo := new(mocks.IPostgreSQLRepository)
	svc := service.NewTixService(
		service.WithPostgreSQLRepository(pqRepo))
	pqRepo.On("UpdateParticipants", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("lorem")).Once()
	err := svc.UpdateParticipantStatus(context.TODO(), "asd", 1, &request.EventRequestUpdateParticipant{
		Status: string(common.ParticipantRequestApproved),
	})
	s.NotNil(err)
}

func (s *tixServiceTestSuite) Test_PublishExportEventDataQueue_ShouldSuccess() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	svc := service.NewTixService(service.WithRedisCache(redisClient))
	err := svc.PublishExportEventDataQueue(context.TODO(), "asd", "pdf", "asd@hello.id")
	s.Nil(err)
}
func (s *tixServiceTestSuite) Test_PublishExportEventDataQueue_ShouldError() {
	s.T().Run("error get", func(t *testing.T) {
		miniRedis := miniredis.RunT(s.T())
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		miniRedis.Close()
		_ = redisClient.Close()
		svc := service.NewTixService(service.WithRedisCache(redisClient))
		err := svc.PublishExportEventDataQueue(context.TODO(), "asd", "pdf", "asd@hello.id")
		s.NotNil(err)
	})
	s.T().Run("error rate limiter", func(t *testing.T) {
		miniRedis := miniredis.RunT(s.T())
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		cacheKey := fmt.Sprintf("%s-%s-%s",
			common.ReqExpEventDataQueueKey, "asd", "asd@hello.id")

		redisClient.Set(context.TODO(), cacheKey, "asd", 1)
		svc := service.NewTixService(service.WithRedisCache(redisClient))
		err := svc.PublishExportEventDataQueue(context.TODO(), "asd", "pdf", "asd@hello.id")
		s.NotNil(err)
		s.Equal(err, common.ErrRateLimitingPushQueue)
		redisClient.Del(context.TODO(), cacheKey)
		miniRedis.Close()
		_ = redisClient.Close()
	})
}

func (s *tixServiceTestSuite) Test_PublishGenerateEventTicketQueue_ShouldSuccess() {
	miniRedis := miniredis.RunT(s.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	svc := service.NewTixService(service.WithRedisCache(redisClient))
	err := svc.PublishGenerateEventTicketQueue(context.TODO(), "asd", 1)
	s.Nil(err)
}
func (s *tixServiceTestSuite) Test_PublishGenerateEventTicketQueue_ShouldError() {
	s.T().Run("error get", func(t *testing.T) {
		miniRedis := miniredis.RunT(s.T())
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		miniRedis.Close()
		_ = redisClient.Close()
		svc := service.NewTixService(service.WithRedisCache(redisClient))
		err := svc.PublishGenerateEventTicketQueue(context.TODO(), "asd", 1)
		s.NotNil(err)
	})
	s.T().Run("error rate limiter", func(t *testing.T) {
		miniRedis := miniredis.RunT(s.T())
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		cacheKey := fmt.Sprintf("%s-%s-%d",
			common.ReqGenEventTixQueueKey, "asd", 1)
		redisClient.Set(context.TODO(), cacheKey, "asd", 1)
		svc := service.NewTixService(service.WithRedisCache(redisClient))
		err := svc.PublishGenerateEventTicketQueue(context.TODO(), "asd", 1)
		s.NotNil(err)
		s.Equal(err, common.ErrRateLimitingPushQueue)
		redisClient.Del(context.TODO(), cacheKey)
		miniRedis.Close()
		_ = redisClient.Close()
	})
}

func TestTixService(t *testing.T) {
	suite.Run(t, new(tixServiceTestSuite))
}
