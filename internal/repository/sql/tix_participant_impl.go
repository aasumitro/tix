package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain/entity"
	"time"
)

func (repository *tixPostgreSQLRepository) CountParticipants(
	ctx context.Context,
	eventID int32,
	participantStatus common.EventParticipantStatus,
	startBetween, endBetween int64,
) int {
	var total int
	query := "SELECT COUNT(*) AS total FROM participants WHERE event_id = $1"
	if participantStatus != common.ParticipantStatusNone {
		switch participantStatus {
		case common.ParticipantRequestApproved:
			query += " AND approved_at IS NOT NULL"
		case common.ParticipantRequestDeclined:
			query += " AND approved_at IS NULL AND declined_at IS NOT NULL"
		case common.ParticipantRequestWaiting:
			query += " AND approved_at IS NULL AND declined_at IS NULL"
		}
	}
	if startBetween != 0 && endBetween != 0 {
		start := time.Unix(startBetween, 0).Format(time.RFC3339)
		end := time.Unix(endBetween, 0).Format(time.RFC3339)
		query += fmt.Sprintf(" AND to_timestamp(created_at) >= '%s' AND to_timestamp(created_at) <= '%s'", start, end)
	}
	if err := repository.db.QueryRowContext(ctx, query, eventID).Scan(&total); err != nil {
		total = 0
	}
	return total
}

func (repository *tixPostgreSQLRepository) GetAllParticipants(
	ctx context.Context,
	eventID int32, filter string,
	startBetween, endBetween int64, limit int32,
	sortKey, sortDir string,
) (
	participants []*entity.Participant,
	err error,
) {
	query := `
	SELECT id, event_id, name, email, phone, job, pop, 
	       dob, approved_at, declined_at, declined_reason 
	FROM participants WHERE event_id = $1
	`
	if filter != "" {
		query += fmt.Sprintf(" AND (name LIKE '%%%s%%' OR email LIKE '%%%s%%' OR phone LIKE '%%%s%%')", filter, filter, filter)
	}
	if startBetween != 0 && endBetween != 0 {
		start := time.Unix(startBetween, 0).Format(time.RFC3339)
		end := time.Unix(endBetween, 0).Format(time.RFC3339)
		query += fmt.Sprintf(" AND to_timestamp(created_at) >= '%s' AND to_timestamp(created_at) <= '%s'", start, end)
	}
	if sortKey != "" && sortDir != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", sortKey, sortDir)
	}
	if limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := repository.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) { _ = rows.Close() }(rows)
	for rows.Next() {
		var participant entity.Participant
		if err := rows.Scan(
			&participant.ID, &participant.EventID,
			&participant.Name, &participant.Email,
			&participant.Phone, &participant.Job,
			&participant.PoP, &participant.DoB,
			&participant.ApprovedAt, &participant.DeclinedAt,
			&participant.DeclinedReason,
		); err != nil {
			return nil, err
		}
		participants = append(participants, &participant)
	}
	return participants, nil
}

func (repository *tixPostgreSQLRepository) GetParticipantByEmailAndEventID(
	ctx context.Context,
	email string, eventID int32,
) (
	participant *entity.Participant,
	err error,
) {
	query := "SELECT id  FROM participants WHERE email = $1 AND event_id = $2 LIMIT 1"
	row := repository.db.QueryRowContext(ctx, query, email, eventID)
	participant = &entity.Participant{}
	if err := row.Scan(&participant.ID); err != nil {
		return nil, err
	}
	return participant, err
}

func (repository *tixPostgreSQLRepository) GetParticipantByIDAndEventID(
	ctx context.Context,
	participantID, eventID int32,
) (
	participant *entity.Participant,
	err error,
) {
	query := "SELECT id, name, email FROM participants WHERE id = $1 AND event_id = $2 LIMIT 1"
	row := repository.db.QueryRowContext(ctx, query, participantID, eventID)
	participant = &entity.Participant{}
	if err := row.Scan(&participant.ID, &participant.Name, &participant.Email); err != nil {
		return nil, err
	}
	return participant, err
}

func (repository *tixPostgreSQLRepository) InsertManyParticipants(
	ctx context.Context,
	participants []*entity.Participant,
	createdAt int64,
) error {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO participants (event_id, name, email, phone, job, pop, dob, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()
	for _, p := range participants {
		if _, err = stmt.ExecContext(
			ctx, p.EventID, p.Name, p.Email,
			p.Phone, p.Job, p.PoP, p.DoB,
			createdAt,
		); err != nil {
			return err
		}
	}
	return nil
}

func (repository *tixPostgreSQLRepository) UpdateParticipants(
	ctx context.Context,
	approvedAt, declinedAt *int64,
	declinedReason *string,
	id int32,
) error {
	query := `
		UPDATE participants 
		SET approved_at = $1, declined_at = $2, declined_reason = $3
		WHERE id = $4 RETURNING id;
	`
	row := repository.db.QueryRowContext(ctx, query, approvedAt, declinedAt, declinedReason, id)
	data := entity.Participant{}
	return row.Scan(&data.ID)
}
