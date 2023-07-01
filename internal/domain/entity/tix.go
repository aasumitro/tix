package entity

import "database/sql"

type (
	User struct {
		ID              int32
		UUID            string
		Username        string
		Email           string
		EmailVerifiedAt sql.NullInt32
		CreatedAt       sql.NullInt32
		UpdatedAt       sql.NullInt32
	}

	Event struct {
		ID                int32
		GoogleFormID      string
		Name              string
		Location          string
		PreregisterDate   int32
		EventDate         int32
		TotalParticipants int32
		CreatedAt         sql.NullInt32
		UpdatedAt         sql.NullInt32
	}

	Participant struct {
		ID             int32
		EventID        int32
		Name           string
		Email          string
		Phone          string
		Job            string
		PoP            string
		DoB            string
		ApprovedAt     sql.NullInt32
		DeclinedAt     sql.NullInt32
		DeclinedReason sql.NullString
		CreatedAt      sql.NullInt32
		UpdatedAt      sql.NullInt32
	}
)
