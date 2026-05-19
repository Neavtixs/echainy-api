package workspace

import (
	"database/sql"
	"fmt"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/entity"
	"github.com/Neavtixs/echainy-api/internal/apps/domain/repository"
	"github.com/Neavtixs/echainy-api/internal/dto"
	"github.com/Neavtixs/echainy-api/internal/errs"
	"github.com/Neavtixs/echainy-api/internal/helper"
	"github.com/google/uuid"
)

type Service struct {
	DB                  *sql.DB
	WorkspaceRepo       *repository.WorkspaceRepo
	WorkspaceMemberRepo *repository.WorkspaceMemberRepo
}

func NewService(db *sql.DB, workspaceRepo *repository.WorkspaceRepo, workspaceMemberRepo *repository.WorkspaceMemberRepo) *Service {
	return &Service{
		DB:                  db,
		WorkspaceRepo:       workspaceRepo,
		WorkspaceMemberRepo: workspaceMemberRepo,
	}
}

func (s *Service) New(input *dto.InputNewWorkspace) (*dto.ResultNewWorkspace, error) {
	userID, ok := input.Ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, errs.ErrInvalidAccessToken
	}

	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	workspace := &entity.Workspace{
		ID:          uuid.NewString(),
		OwnerUserID: userID,
		Name:        input.Name,
		Slug:        helper.GenerateSlug(input.Name),
		AvatarURL:   input.AvatarURL,
	}

	if err := s.WorkspaceRepo.Create(tx, input.Ctx, workspace); err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}

	workspaceMember := &entity.WorkspaceMember{
		ID:          uuid.NewString(),
		WorkspaceID: workspace.ID,
		UserID:      userID,
		Role:        "OWNER",
	}

	if err := s.WorkspaceMemberRepo.Create(tx, input.Ctx, workspaceMember); err != nil {
		return nil, fmt.Errorf("create workspace member: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &dto.ResultNewWorkspace{
		ID:          workspace.ID,
		OwnerUserID: workspace.OwnerUserID,
		Name:        workspace.Name,
		Slug:        workspace.Slug,
		AvatarURL:   workspace.AvatarURL,
		Role:        workspaceMember.Role,
	}, nil
}

func (s *Service) List(input *dto.InputListWorkspace) (*dto.ResultListWorkspace, error) {
	userID, ok := input.Ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, errs.ErrInvalidAccessToken
	}

	workspaces, err := s.WorkspaceMemberRepo.FindWorkspacesByUserID(s.DB, input.Ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find workspaces by user id: %w", err)
	}

	resultWorkspaces := make([]dto.ResultListWorkspaceItem, 0, len(workspaces))
	for _, workspace := range workspaces {
		resultWorkspaces = append(resultWorkspaces, dto.ResultListWorkspaceItem{
			ID:          workspace.ID,
			OwnerUserID: workspace.OwnerUserID,
			Name:        workspace.Name,
			Slug:        workspace.Slug,
			AvatarURL:   workspace.AvatarURL,
			Role:        workspace.Role,
		})
	}

	return &dto.ResultListWorkspace{
		Workspaces: resultWorkspaces,
	}, nil
}
