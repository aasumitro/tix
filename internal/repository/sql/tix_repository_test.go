package sql_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain"
	"github.com/aasumitro/tix/internal/domain/entity"
	"github.com/aasumitro/tix/internal/domain/request"
	sqlRepository "github.com/aasumitro/tix/internal/repository/sql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"strconv"
	"testing"
	"time"
)

type tixSQLRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo domain.IPostgreSQLRepository
}

func (s *tixSQLRepositoryTestSuite) SetupSuite() {
	var err error
	s.db, s.mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(
			sqlmock.QueryMatcherRegexp))
	require.NoError(s.T(), err)
	s.repo = sqlRepository.NewTixPostgreSQLRepository(s.db)
}

func (s *tixSQLRepositoryTestSuite) AfterTest(_, _ string) {
	s.NoError(s.mock.ExpectationsWereMet())
}

// ===============================================================
// PART OF USER TEST CASE
// ===============================================================
func (s *tixSQLRepositoryTestSuite) Test_CountUsers_ShouldSuccess() {
	count := s.mock.
		NewRows([]string{"total"}).
		AddRow(1)
	query := "SELECT COUNT(*) AS total FROM users"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnRows(count)
	res := s.repo.CountUsers(context.TODO())
	s.NotZero(res)
	s.Equal(1, res)
}
func (s *tixSQLRepositoryTestSuite) Test_CountUsers_ShouldError() {
	query := "SELECT COUNT(*) AS total FROM users"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New("lorem"))
	res := s.repo.CountUsers(context.TODO())
	s.Zero(res)
	s.Equal(0, res)
}

func (s *tixSQLRepositoryTestSuite) Test_GetAllUser_ShouldSuccess() {
	dataMock := s.mock.
		NewRows([]string{"id", "uuid", "username", "email", "email_verified_at"}).
		AddRow(1, "123", "tix", "hello@tix.id", nil)
	query := "SELECT id, uuid, username, email, email_verified_at FROM users WHERE email LIKE %$1%"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
	data, err := s.repo.GetAllUsers(context.TODO(), "hello@tix.id")
	s.NotNil(data)
	s.NoError(err)
	s.Equal(len(data), 1)
}
func (s *tixSQLRepositoryTestSuite) Test_GetAllUser_ShouldError() {
	query := "SELECT id, uuid, username, email, email_verified_at FROM users WHERE email LIKE %$1%"
	expectedQuery := regexp.QuoteMeta(query)
	s.T().Run("ERROR FROM QUERY", func(t *testing.T) {
		s.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New("lorem"))
		data, err := s.repo.GetAllUsers(context.TODO(), "hello@tix.id")
		s.Nil(data)
		s.Error(err)
	})
	s.T().Run("ERROR FROM SCAN", func(t *testing.T) {
		dataMock := s.mock.
			NewRows([]string{"id", "uuid", "username", "email", "email_verified_at"}).
			AddRow(1, "123", "tix", "hello@tix.id", nil).
			AddRow(2, nil, nil, nil, nil)
		s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
		data, err := s.repo.GetAllUsers(context.TODO(), "hello@tix.id")
		s.Nil(data)
		s.Error(err)
	})
}

func (s *tixSQLRepositoryTestSuite) Test_GetUserByEmail_ShouldSuccess() {
	dataMock := s.mock.
		NewRows([]string{"id", "uuid", "username", "email", "email_verified_at"}).
		AddRow(1, "123", "tix", "hello@tix.id", nil)
	query := "SELECT id, uuid, username, email, email_verified_at FROM users WHERE email = $1 LIMIT 1"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
	data, err := s.repo.GetUserByEmail(context.TODO(), "hello@tix.id")
	s.NotNil(data)
	s.NoError(err)
	s.Equal(data.Email, "hello@tix.id")
}
func (s *tixSQLRepositoryTestSuite) Test_GetUserByEmail_ShouldError() {
	dataMock := s.mock.
		NewRows([]string{"id", "uuid", "username", "email", "email_verified_at"}).
		AddRow(1, nil, nil, "hello@tix.id", nil)
	query := "SELECT id, uuid, username, email, email_verified_at FROM users WHERE email = $1 LIMIT 1"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
	data, err := s.repo.GetUserByEmail(context.TODO(), "hello@tix.id")
	s.Nil(data)
	s.Error(err)
}

func (s *tixSQLRepositoryTestSuite) Test_UpdateUserVerifiedTime_ShouldSuccess() {
	dataMock := s.mock.NewRows([]string{"id"}).AddRow(1)
	query := "UPDATE users SET email_verified_at = $1, updated_at = $2 WHERE email = $3 RETURNING id"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).
		WithArgs(time.Now().Unix(), time.Now().Unix(), "hello@tix.id").
		WillReturnRows(dataMock)
	err := s.repo.UpdateUserVerifiedTime(context.TODO(), "hello@tix.id")
	s.Nil(err)
	s.NoError(err)
}
func (s *tixSQLRepositoryTestSuite) Test_UpdateUserVerifiedTime_ShouldError() {
	query := "UPDATE users SET email_verified_at = $1, updated_at = $2 WHERE email = $3 RETURNING id"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).
		WithArgs(time.Now().Unix(), time.Now().Unix(), "hello@tix.id").
		WillReturnError(errors.New("lorem"))
	err := s.repo.UpdateUserVerifiedTime(context.TODO(), "hello@tix.id")
	s.NotNil(err)
	s.Error(err)
}

func (s *tixSQLRepositoryTestSuite) Test_DeleteUser_ShouldSuccess() {
	query := "DELETE FROM users WHERE uuid = $1"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectExec(expectedQuery).
		WithArgs("123-asd-456").
		WillReturnResult(sqlmock.NewResult(0, 1))
	err := s.repo.DeleteUser(context.TODO(), "123-asd-456")
	s.Nil(err)
	s.NoError(err)
}
func (s *tixSQLRepositoryTestSuite) Test_DeleteUser_ShouldError() {
	query := "DELETE FROM users WHERE uuid = $1"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectExec(expectedQuery).
		WithArgs("123-asd-456").
		WillReturnError(errors.New("lorem"))
	err := s.repo.DeleteUser(context.TODO(), "123-asd-456")
	s.NotNil(err)
	s.Error(err)
}

// ===============================================================
// PART OF EVENT TEST CASE
// ===============================================================
func (s *tixSQLRepositoryTestSuite) Test_GetAllEvent_ShouldSuccess() {
	dataMock := s.mock.
		NewRows([]string{"id", "google_form_id", "name", "location", "preregister_date", "event_date", "total_participants"}).
		AddRow(1, "123", "tix", "jalan tix", time.Now().Unix(), time.Now().Unix(), 10)
	query := `
		SELECT 
		    events.id, 
		    events.google_form_id, 
		    events.name, 
		    events.location, 
		    events.preregister_date, 
		    events.event_date,
			COUNT(participants.id) AS total_participants
		FROM events 
		LEFT JOIN participants on events.id = participants.event_id
		GROUP BY events.id ORDER BY events.id DESC;`
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
	data, err := s.repo.GetAllEvents(context.TODO())
	s.NotNil(data)
	s.NoError(err)
	s.Equal(len(data), 1)
}
func (s *tixSQLRepositoryTestSuite) Test_GetAllEvent_ShouldError() {
	query := `
		SELECT 
		    events.id, 
		    events.google_form_id, 
		    events.name, 
		    events.location, 
		    events.preregister_date, 
		    events.event_date,
			COUNT(participants.id) AS total_participants
		FROM events 
		LEFT JOIN participants on events.id = participants.event_id
		GROUP BY events.id ORDER BY events.id DESC;`
	expectedQuery := regexp.QuoteMeta(query)
	s.T().Run("ERROR FROM QUERY", func(t *testing.T) {
		s.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New("lorem"))
		data, err := s.repo.GetAllEvents(context.TODO())
		s.Nil(data)
		s.Error(err)
	})
	s.T().Run("ERROR FROM SCAN", func(t *testing.T) {
		dataMock := s.mock.
			NewRows([]string{"id", "google_form_id", "name", "location", "preregister_date", "event_date", "total_participants"}).
			AddRow(2, nil, nil, nil, nil, nil, nil)
		s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
		data, err := s.repo.GetAllEvents(context.TODO())
		s.Nil(data)
		s.Error(err)
	})
}

func (s *tixSQLRepositoryTestSuite) Test_GetEventByGoogleFormID_ShouldSuccess() {
	dataMock := s.mock.
		NewRows([]string{"id", "google_form_id", "name", "location", "preregister_date", "event_date", "total_participants"}).
		AddRow(1, "123", "tix", "jalan tix", time.Now().Unix(), time.Now().Unix(), 10)
	query := `
		SELECT 
		    events.id, 
		    events.google_form_id, 
		    events.name, 
		    events.location, 
		    events.preregister_date, 
		    events.event_date,
		    COUNT(participants.id) AS total_participants
		FROM events
		LEFT JOIN participants on events.id = participants.event_id
		WHERE google_form_id = $1
		GROUP BY events.id
		LIMIT 1;`
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
	data, err := s.repo.GetEventByGoogleFormID(context.TODO(), "123")
	s.NotNil(data)
	s.NoError(err)
	s.Equal(data.GoogleFormID, "123")
}
func (s *tixSQLRepositoryTestSuite) Test_GetEventByGoogleFormID_ShouldError() {
	query := `
		SELECT 
		    events.id, 
		    events.google_form_id, 
		    events.name, 
		    events.location, 
		    events.preregister_date, 
		    events.event_date,
		    COUNT(participants.id) AS total_participants
		FROM events
		LEFT JOIN participants on events.id = participants.event_id
		WHERE google_form_id = $1
		GROUP BY events.id
		LIMIT 1;`
	expectedQuery := regexp.QuoteMeta(query)
	s.T().Run("ERROR FROM QUERY", func(t *testing.T) {
		s.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New("lorem"))
		data, err := s.repo.GetEventByGoogleFormID(context.TODO(), "123")
		s.Nil(data)
		s.Error(err)
	})
	s.T().Run("ERROR FROM SCAN", func(t *testing.T) {
		dataMock := s.mock.
			NewRows([]string{"id", "google_form_id", "name", "location", "preregister_date", "event_date", "total_participants"}).
			AddRow(2, nil, nil, nil, nil, nil, nil)
		s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
		data, err := s.repo.GetEventByGoogleFormID(context.TODO(), "123")
		s.Nil(data)
		s.Error(err)
	})
}

func (s *tixSQLRepositoryTestSuite) Test_InsertNewEvent_ShouldSuccess() {
	now := time.Now().Unix()
	ns := strconv.FormatInt(now, 10)
	rows := s.mock.NewRows([]string{"id", "google_form_id", "name", "location", "preregister_date", "event_date"}).
		AddRow(1, "123", "tix", "jalan tix", ns, ns)
	query := `
		INSERT INTO events (google_form_id, name, location, preregister_date, event_date) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id, google_form_id, name, location, preregister_date, event_date`
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).
		WithArgs("123", "tix", "jalan tix", ns, ns).
		WillReturnRows(rows)
	res, err := s.repo.InsertNewEvent(context.TODO(), &request.EventRequestMakeNew{
		GoogleFormID:    "123",
		Name:            "tix",
		PreregisterDate: ns,
		EventDate:       ns,
		Location:        "jalan tix",
	})
	s.Nil(err)
	s.NoError(err)
	s.NotNil(res)
}
func (s *tixSQLRepositoryTestSuite) Test_InsertNewEvent_ShouldError() {
	now := time.Now().Unix()
	ns := strconv.FormatInt(now, 10)
	query := `
		INSERT INTO events (google_form_id, name, location, preregister_date, event_date) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id, google_form_id, name, location, preregister_date, event_date`
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).
		WithArgs("123", "tix", "jalan tix", ns, ns).
		WillReturnError(errors.New("lorem"))
	res, err := s.repo.InsertNewEvent(context.TODO(), &request.EventRequestMakeNew{
		GoogleFormID:    "123",
		Name:            "tix",
		PreregisterDate: ns,
		EventDate:       ns,
		Location:        "jalan tix",
	})
	s.Nil(res)
	s.Error(err)
	s.NotNil(err)
}

// ===============================================================
// PART OF PARTICIPANT TEST CASE
// ===============================================================
func (s *tixSQLRepositoryTestSuite) Test_CountParticipant_ShouldSuccess() {
	s.T().Run("COUNT APPROVED", func(t *testing.T) {
		now := time.Now().Unix()
		start := time.Unix(now, 0).Format(time.RFC3339)
		end := time.Unix(now, 0).Format(time.RFC3339)
		count := s.mock.
			NewRows([]string{"total"}).
			AddRow(1)
		query := "SELECT COUNT(*) AS total FROM participants WHERE event_id = $1 AND approved_at IS NOT NULL"
		query += fmt.Sprintf(" AND to_timestamp(created_at) >= '%s' AND to_timestamp(created_at) <= '%s'", start, end)
		expectedQuery := regexp.QuoteMeta(query)
		s.mock.ExpectQuery(expectedQuery).WillReturnRows(count)
		res := s.repo.CountParticipants(context.TODO(), 1, common.ParticipantRequestApproved, now, now)
		s.NotZero(res)
		s.Equal(1, res)
	})
	s.T().Run("COUNT DECLINED", func(t *testing.T) {
		count := s.mock.
			NewRows([]string{"total"}).
			AddRow(1)
		query := "SELECT COUNT(*) AS total FROM participants WHERE event_id = $1 AND approved_at IS NULL AND declined_at IS NOT NULL"
		expectedQuery := regexp.QuoteMeta(query)
		s.mock.ExpectQuery(expectedQuery).WillReturnRows(count)
		res := s.repo.CountParticipants(context.TODO(), 1, common.ParticipantRequestDeclined, 0, 0)
		s.NotZero(res)
		s.Equal(1, res)
	})
	s.T().Run("COUNT WAITING", func(t *testing.T) {
		count := s.mock.
			NewRows([]string{"total"}).
			AddRow(1)
		query := "SELECT COUNT(*) AS total FROM participants WHERE event_id = $1 AND approved_at IS NULL AND declined_at IS NULL"
		expectedQuery := regexp.QuoteMeta(query)
		s.mock.ExpectQuery(expectedQuery).WillReturnRows(count)
		res := s.repo.CountParticipants(context.TODO(), 1, common.ParticipantRequestWaiting, 0, 0)
		s.NotZero(res)
		s.Equal(1, res)
	})
}
func (s *tixSQLRepositoryTestSuite) Test_CountParticipant_ShouldError() {
	query := "SELECT COUNT(*) AS total FROM participants WHERE event_id = $1 AND approved_at IS NULL AND declined_at IS NULL"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New("lorem"))
	res := s.repo.CountParticipants(context.TODO(), 1, common.ParticipantRequestWaiting, 0, 0)
	s.Zero(res)
	s.Equal(0, res)
}

func (s *tixSQLRepositoryTestSuite) Test_GetAllParticipant_ShouldSuccess() {
	dataMock := s.mock.
		NewRows([]string{"id", "event_id", "name", "email", "phone", "job", "pop", "dob", "approved_at", "declined_at", "declined_reason"}).
		AddRow(1, 1, "tix", "hellO@tix.id", "082271119900", "SE", "http://bukti.id/123", "1990-12-12", nil, nil, nil)
	query := `
	SELECT id, event_id, name, email, phone, job, pop, 
	       dob, approved_at, declined_at, declined_reason 
	FROM participants WHERE event_id = $1`
	query += fmt.Sprintf(" AND (name LIKE '%%%s%%' OR email LIKE '%%%s%%' OR phone LIKE '%%%s%%')", "tix", "tix", "tix")
	now := time.Now().Unix()
	start := time.Unix(now, 0).Format(time.RFC3339)
	end := time.Unix(now, 0).Format(time.RFC3339)
	query += fmt.Sprintf(" AND to_timestamp(created_at) >= '%s' AND to_timestamp(created_at) <= '%s'", start, end)
	query += fmt.Sprintf(" ORDER BY %s %s", "name", "asc")
	query += fmt.Sprintf(" LIMIT %d", 10)
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
	res, err := s.repo.GetAllParticipants(context.TODO(), 1, "tix", time.Now().Unix(), time.Now().Unix(), 10, "name", "asc")
	s.Nil(err)
	s.NoError(err)
	s.NotNil(res)
}
func (s *tixSQLRepositoryTestSuite) Test_GetAllParticipant_ShouldError() {
	s.T().Run("ERROR FROM QUERY", func(t *testing.T) {
		query := `
		SELECT id, event_id, name, email, phone, job, pop, 
			   dob, approved_at, declined_at, declined_reason 
		FROM participants WHERE event_id = $1`
		expectedQuery := regexp.QuoteMeta(query)
		s.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New("hello"))
		res, err := s.repo.GetAllParticipants(context.TODO(), 1, "", 0, 0, 0, "", "")
		s.NotNil(err)
		s.Error(err)
		s.Nil(res)
	})
	s.T().Run("ERROR FROM SCAN", func(t *testing.T) {
		dataMock := s.mock.
			NewRows([]string{"id", "event_id", "name", "email", "phone", "job", "pop", "dob", "approved_at", "declined_at", "declined_reason"}).
			AddRow(1, 1, nil, nil, "082271119900", "SE", "http://bukti.id/123", "1990-12-12", nil, nil, nil)
		query := `
		SELECT id, event_id, name, email, phone, job, pop, 
			   dob, approved_at, declined_at, declined_reason 
		FROM participants WHERE event_id = $1`
		expectedQuery := regexp.QuoteMeta(query)
		s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
		res, err := s.repo.GetAllParticipants(context.TODO(), 1, "", 0, 0, 0, "", "")
		s.NotNil(err)
		s.Error(err)
		s.Nil(res)
	})
}

func (s *tixSQLRepositoryTestSuite) Test_GetParticipantByEmailAndEventID_ShouldSuccess() {
	dataMock := s.mock.NewRows([]string{"id"}).AddRow(1)
	query := "SELECT id  FROM participants WHERE email = $1 AND event_id = $2 LIMIT 1"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
	data, err := s.repo.GetParticipantByEmailAndEventID(context.TODO(), "hello@tix.id", 1)
	s.NotNil(data)
	s.NoError(err)
	s.Equal(data.ID, int32(1))
}
func (s *tixSQLRepositoryTestSuite) Test_GetParticipantByEmailAndEventID_ShouldError() {
	dataMock := s.mock.NewRows([]string{"id"}).AddRow(nil)
	query := "SELECT id  FROM participants WHERE email = $1 AND event_id = $2 LIMIT 1"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
	data, err := s.repo.GetParticipantByEmailAndEventID(context.TODO(), "hello@tix.id", 1)
	s.Nil(data)
	s.Error(err)
}

func (s *tixSQLRepositoryTestSuite) Test_GetParticipantByParticipantIDAndEventID_ShouldSuccess() {
	dataMock := s.mock.NewRows([]string{"id", "name", "email"}).AddRow(1, "lorem", "lorem@lorem.id")
	query := "SELECT  id, name, email FROM participants WHERE id = $1 AND event_id = $2 LIMIT 1"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
	data, err := s.repo.GetParticipantByIDAndEventID(context.TODO(), 1, 1)
	s.NotNil(data)
	s.NoError(err)
	s.Equal(data.ID, int32(1))
}
func (s *tixSQLRepositoryTestSuite) Test_GetParticipantByParticipantIDAndEventID_ShouldError() {
	dataMock := s.mock.NewRows([]string{"id", "name", "email"}).AddRow(nil, nil, nil)
	query := "SELECT id, name, email FROM participants WHERE id = $1 AND event_id = $2 LIMIT 1"
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).WillReturnRows(dataMock)
	data, err := s.repo.GetParticipantByIDAndEventID(context.TODO(), 1, 1)
	s.Nil(data)
	s.Error(err)
}

func (s *tixSQLRepositoryTestSuite) Test_InsertManyParticipants_ShouldSuccess() {
	s.mock.ExpectBegin()
	s.mock.ExpectPrepare(`.*INSERT INTO participants \(event_id, name, email, phone, job, pop, dob, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\).*`)
	s.mock.ExpectExec(`.*INSERT INTO participants \(event_id, name, email, phone, job, pop, dob, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\).*`).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.repo.InsertManyParticipants(context.Background(), []*entity.Participant{{
		EventID: 1,
		Name:    "tix",
		Email:   "hello@tix.id",
		Phone:   "08272229292",
		Job:     "lorem",
		PoP:     "http://bukti.id/1312312",
		DoB:     "1990-12-12",
	}}, time.Now().Unix())
	s.Nil(err)
	s.NoError(err)
}
func (s *tixSQLRepositoryTestSuite) Test_InsertManyParticipants_ShouldError() {
	s.T().Run("ERROR BEGIN TX", func(t *testing.T) {
		s.mock.ExpectBegin().WillReturnError(errors.New("lorem"))
		err := s.repo.InsertManyParticipants(context.Background(), []*entity.Participant{{
			EventID: 1,
			Name:    "tix",
			Email:   "hello@tix.id",
			Phone:   "08272229292",
			Job:     "lorem",
			PoP:     "http://bukti.id/1312312",
			DoB:     "1990-12-12",
		}}, time.Now().Unix())
		s.NotNil(err)
		s.Error(err)
	})
	s.T().Run("ERROR PREPARE TX", func(t *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectPrepare(`.*INSERT INTO participants \(event_id, name, email, phone, job, pop, dob, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\).*`).WillReturnError(errors.New("lorem"))
		err := s.repo.InsertManyParticipants(context.Background(), []*entity.Participant{{
			EventID: 1,
			Name:    "tix",
			Email:   "hello@tix.id",
			Phone:   "08272229292",
			Job:     "lorem",
			PoP:     "http://bukti.id/1312312",
			DoB:     "1990-12-12",
		}}, time.Now().Unix())
		s.NotNil(err)
		s.Error(err)
	})
	s.T().Run("ERROR EXEC TX", func(t *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectPrepare(`.*INSERT INTO participants \(event_id, name, email, phone, job, pop, dob, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\).*`)
		s.mock.ExpectExec(`.*INSERT INTO participants \(event_id, name, email, phone, job, pop, dob, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\).*`).WillReturnError(errors.New("lorem"))
		err := s.repo.InsertManyParticipants(context.Background(), []*entity.Participant{{
			EventID: 1,
			Name:    "tix",
			Email:   "hello@tix.id",
			Phone:   "08272229292",
			Job:     "lorem",
			PoP:     "http://bukti.id/1312312",
			DoB:     "1990-12-12",
		}}, time.Now().Unix())
		s.NotNil(err)
		s.Error(err)
	})
}

func (s *tixSQLRepositoryTestSuite) Test_UpdateParticipants_ShouldSuccess() {
	dataMock := s.mock.NewRows([]string{"id"}).AddRow(1)
	query := `
		UPDATE participants 
		SET approved_at = $1, declined_at = $2, declined_reason = $3
		WHERE id = $4 RETURNING id;
	`
	now := time.Now().Unix()
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).
		WithArgs(now, nil, nil, 1).
		WillReturnRows(dataMock)
	err := s.repo.UpdateParticipants(context.TODO(), &now, nil, nil, 1)
	s.Nil(err)
	s.NoError(err)
}
func (s *tixSQLRepositoryTestSuite) Test_UpdateParticipants_ShouldError() {
	query := `
		UPDATE participants 
		SET approved_at = $1, declined_at = $2, declined_reason = $3
		WHERE id = $4 RETURNING id;
	`
	now := time.Now().Unix()
	expectedQuery := regexp.QuoteMeta(query)
	s.mock.ExpectQuery(expectedQuery).
		WithArgs(now, nil, nil, 1).
		WillReturnError(errors.New("lorem"))
	err := s.repo.UpdateParticipants(context.TODO(), &now, nil, nil, 1)
	s.NotNil(err)
	s.Error(err)
}

func TestTixSQLRepository(t *testing.T) {
	suite.Run(t, new(tixSQLRepositoryTestSuite))
}
