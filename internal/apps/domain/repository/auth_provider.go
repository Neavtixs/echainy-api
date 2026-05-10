package repository

import (
	"context"
	"database/sql"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/entity"
	"github.com/Neavtixs/echainy-api/internal/errs"
	"github.com/lib/pq"
)

type AuthProviderRepo struct {
}

func NewAuthProviderRepo() *AuthProviderRepo {
	return &AuthProviderRepo{}
}

func (r *AuthProviderRepo) FindByUserID(db *sql.DB, ctx context.Context, userID string, authProvider *entity.AuthProvider) error {
	query := `
		SELECT id, user_id, provider_name
		FROM auth_providers
		WHERE user_id = $1
	`

	result := db.QueryRowContext(ctx, query, userID)
	if err := result.Scan(&authProvider.ID, &authProvider.UserID, &authProvider.ProviderName); err != nil {
		if err == sql.ErrNoRows {
			return errs.ErrDataNotFound
		}
		return err
	}

	return nil
}

func (r *AuthProviderRepo) Create(db *sql.Tx, ctx context.Context, authProvider *entity.AuthProvider) error {
	query := `
		INSERT INTO auth_providers(id, user_id, provider_name)
		VALUES($1, $2, $3)
	`

	result, err := db.ExecContext(
		ctx,
		query,
		authProvider.ID,
		authProvider.UserID,
		authProvider.ProviderName,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code {
			case "23505":
				return errs.ErrUserIDUsed
			case "23503":
				return errs.ErrUserIDNotFound
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

	return nil
}
