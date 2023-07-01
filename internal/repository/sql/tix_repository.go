package sql

import (
	"database/sql"
	"github.com/aasumitro/tix/internal/domain"
)

type tixPostgreSQLRepository struct {
	db *sql.DB
}

func NewTixPostgreSQLRepository(
	db *sql.DB,
) domain.IPostgreSQLRepository {
	return &tixPostgreSQLRepository{db: db}
}
