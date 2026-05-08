package repository

import (
	"context"
	"database/sql"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/entity"
	"github.com/Neavtixs/echainy-api/internal/errs"
	"github.com/lib/pq"
)

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
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
