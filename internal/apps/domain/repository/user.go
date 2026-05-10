package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/entity"
	"github.com/Neavtixs/echainy-api/internal/errs"
	"github.com/lib/pq"
)

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) FindByID(db *sql.DB, ctx context.Context, id string, user *entity.User) error {
	query := `
		SELECT id, email, password
		FROM users
		WHERE id = $1
	`

	result := db.QueryRowContext(ctx, query, id)
	if err := result.Err(); err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "22P02" {
				return errs.ErrInvalidType
			}
		}

		return err
	}

	if err := result.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrDataNotFound
		}
		return err
	}

	return nil
}

func (r *UserRepo) FindByEmail(db *sql.DB, ctx context.Context, email string, user *entity.User) error {
	query := `
		SELECT id, email, password
		FROM users
		WHERE email = $1
	`

	result := db.QueryRowContext(ctx, query, email)
	if err := result.Err(); err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "22P02" {
				return errs.ErrInvalidType
			}
		}

		return err
	}

	if err := result.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrDataNotFound
		}
		return err
	}

	return nil
}

// langsung validasi email yang duplikat
func (r *UserRepo) Create(db *sql.Tx, ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users(id, email, password)
		VALUES($1, $2, $3)
	`

	result, err := db.ExecContext(ctx, query, user.ID, user.Email, user.Password)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return errs.ErrEmailUsed
			}
		}

		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrFailedCreateData
	}

	return err
}
