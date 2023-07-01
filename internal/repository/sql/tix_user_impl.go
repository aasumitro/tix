package sql

import (
	"context"
	"database/sql"
	"github.com/aasumitro/tix/internal/domain/entity"
	"time"
)

func (repository *tixPostgreSQLRepository) CountUsers(
	ctx context.Context,
) int {
	var total int
	query := "SELECT COUNT(*) AS total FROM users"
	if err := repository.db.QueryRowContext(ctx, query).Scan(&total); err != nil {
		total = 0
	}
	return total
}

func (repository *tixPostgreSQLRepository) GetAllUsers(
	ctx context.Context,
	email string,
) (
	users []*entity.User,
	err error,
) {
	query := "SELECT id, uuid, username, email, email_verified_at FROM users"
	if email != "" {
		query += " WHERE email LIKE %$1%"
	}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) { _ = rows.Close() }(rows)
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.ID, &user.UUID, &user.Username,
			&user.Email, &user.EmailVerifiedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (repository *tixPostgreSQLRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (
	user *entity.User,
	err error,
) {
	query := "SELECT id, uuid, username, email, email_verified_at FROM users WHERE email = $1 LIMIT 1"
	row := repository.db.QueryRowContext(ctx, query, email)
	user = &entity.User{}
	if err := row.Scan(
		&user.ID, &user.UUID, &user.Username,
		&user.Email, &user.EmailVerifiedAt,
	); err != nil {
		return nil, err
	}
	return user, err
}

func (repository *tixPostgreSQLRepository) UpdateUserVerifiedTime(
	ctx context.Context, email string,
) error {
	now := time.Now().Unix()
	query := "UPDATE users SET email_verified_at = $1, updated_at = $2 WHERE email = $3 RETURNING id"
	row := repository.db.QueryRowContext(ctx, query, now, now, email)
	var user entity.User
	return row.Scan(&user.ID)
}

func (repository *tixPostgreSQLRepository) DeleteUser(
	ctx context.Context, uuid string,
) error {
	query := "DELETE FROM users WHERE uuid = $1"
	_, err := repository.db.ExecContext(ctx, query, uuid)
	return err
}
