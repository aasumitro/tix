package sql

import (
	"context"
	"database/sql"
	"github.com/aasumitro/tix/internal/domain/entity"
	"github.com/aasumitro/tix/internal/domain/request"
)

func (repository *tixPostgreSQLRepository) GetAllEvents(
	ctx context.Context,
) (
	events []*entity.Event,
	err error,
) {
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
		GROUP BY events.id ORDER BY events.id DESC;
	`
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) { _ = rows.Close() }(rows)
	for rows.Next() {
		var event entity.Event
		if err := rows.Scan(
			&event.ID, &event.GoogleFormID,
			&event.Name, &event.Location,
			&event.PreregisterDate,
			&event.EventDate,
			&event.TotalParticipants,
		); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func (repository *tixPostgreSQLRepository) GetEventByGoogleFormID(
	ctx context.Context,
	googleFormID string,
) (
	event *entity.Event,
	err error,
) {
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
		LIMIT 1;
	`
	row := repository.db.QueryRowContext(ctx, query, googleFormID)
	event = &entity.Event{}
	if err := row.Scan(
		&event.ID, &event.GoogleFormID,
		&event.Name, &event.Location,
		&event.PreregisterDate,
		&event.EventDate,
		&event.TotalParticipants,
	); err != nil {
		return nil, err
	}
	return event, nil
}

func (repository *tixPostgreSQLRepository) InsertNewEvent(
	ctx context.Context,
	param *request.EventRequestMakeNew,
) (
	event *entity.Event,
	err error,
) {
	query := `
		INSERT INTO events (google_form_id, name, location, preregister_date, event_date) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id, google_form_id, name, location, preregister_date, event_date
	`
	row := repository.db.QueryRowContext(
		ctx, query, param.GoogleFormID, param.Name,
		param.Location, param.PreregisterDate, param.EventDate)
	event = &entity.Event{}
	if err := row.Scan(&event.ID, &event.GoogleFormID,
		&event.Name, &event.Location,
		&event.PreregisterDate,
		&event.EventDate,
	); err != nil {
		return nil, err
	}
	return event, nil
}
