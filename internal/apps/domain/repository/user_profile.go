package repository

import (
	"context"
	"database/sql"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/entity"
	"github.com/Neavtixs/echainy-api/internal/errs"
	"github.com/lib/pq"
)

type UserProfileRepo struct {
}

func NewUserProfileRepo() *UserProfileRepo {
	return &UserProfileRepo{}
}

func (r *UserProfileRepo) Create(db *sql.Tx, ctx context.Context, userProfile *entity.UserProfile) error {
	query := `
		INSERT INTO user_profiles(id, user_id, name, avatar_url)
		VALUES($1, $2, $3, $4)
	`

	result, err := db.ExecContext(
		ctx,
		query,
		userProfile.ID,
		userProfile.UserID,
		userProfile.Name,
		userProfile.AvatarURL,
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
