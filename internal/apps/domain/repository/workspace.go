package repository

import (
	"context"
	"database/sql"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/entity"
	"github.com/Neavtixs/echainy-api/internal/errs"
	"github.com/lib/pq"
)

type WorkspaceRepo struct {
}

func NewWorkspaceRepo() *WorkspaceRepo {
	return &WorkspaceRepo{}
}

func (r *WorkspaceRepo) Create(db *sql.Tx, ctx context.Context, workspace *entity.Workspace) error {
	query := `
		INSERT INTO workspaces(id, owner_user_id, name, slug, avatar_url)
		VALUES($1, $2, $3, $4, $5)
	`

	result, err := db.ExecContext(
		ctx,
		query,
		workspace.ID,
		workspace.OwnerUserID,
		workspace.Name,
		workspace.Slug,
		workspace.AvatarURL,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code {
			case "23505":
				return errs.ErrSlugUsed
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
