package repository

import (
	"context"
	"database/sql"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/entity"
	"github.com/Neavtixs/echainy-api/internal/errs"
	"github.com/lib/pq"
)

type WorkspaceMemberRepo struct {
}

func NewWorkspaceMemberRepo() *WorkspaceMemberRepo {
	return &WorkspaceMemberRepo{}
}

func (r *WorkspaceMemberRepo) Create(db *sql.Tx, ctx context.Context, workspaceMember *entity.WorkspaceMember) error {
	query := `
		INSERT INTO workspace_members(id, workspace_id, user_id, role)
		VALUES($1, $2, $3, $4)
	`

	result, err := db.ExecContext(
		ctx,
		query,
		workspaceMember.ID,
		workspaceMember.WorkspaceID,
		workspaceMember.UserID,
		workspaceMember.Role,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code {
			case "23505":
				return errs.ErrWorkspaceIDUsed
			case "23503":
				if pgErr.Constraint == "fk_workspace_members_workspace" {
					return errs.ErrWorkspaceIDNotFound
				}
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
