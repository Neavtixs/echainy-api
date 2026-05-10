package repository

import (
	"context"
	"testing"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/entity"
	"github.com/Neavtixs/echainy-api/internal/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceMemberRepo_Create(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepo()
	workspaceRepo := NewWorkspaceRepo()
	workspaceMemberRepo := NewWorkspaceMemberRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	user := &entity.User{
		ID:       uuid.NewString(),
		Email:    "workspace-member-create@test.com",
		Password: "password",
	}
	err = userRepo.Create(tx, ctx, user)
	require.NoError(t, err)

	workspace := &entity.Workspace{
		ID:          uuid.NewString(),
		OwnerUserID: user.ID,
		Name:        "Workspace Member Test",
		Slug:        "workspace-member-test",
		AvatarURL:   "",
	}
	err = workspaceRepo.Create(tx, ctx, workspace)
	require.NoError(t, err)

	workspaceMember := &entity.WorkspaceMember{
		ID:          uuid.NewString(),
		WorkspaceID: workspace.ID,
		UserID:      user.ID,
		Role:        "ADMIN",
	}

	err = workspaceMemberRepo.Create(tx, ctx, workspaceMember)
	assert.NoError(t, err)
}

func TestWorkspaceMemberRepo_Create_DuplicateWorkspaceUser(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepo()
	workspaceRepo := NewWorkspaceRepo()
	workspaceMemberRepo := NewWorkspaceMemberRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	user := &entity.User{
		ID:       uuid.NewString(),
		Email:    "workspace-member-duplicate@test.com",
		Password: "password",
	}
	err = userRepo.Create(tx, ctx, user)
	require.NoError(t, err)

	workspace := &entity.Workspace{
		ID:          uuid.NewString(),
		OwnerUserID: user.ID,
		Name:        "Workspace Member Duplicate Test",
		Slug:        "workspace-member-duplicate-test",
		AvatarURL:   "",
	}
	err = workspaceRepo.Create(tx, ctx, workspace)
	require.NoError(t, err)

	workspaceMember1 := &entity.WorkspaceMember{
		ID:          uuid.NewString(),
		WorkspaceID: workspace.ID,
		UserID:      user.ID,
		Role:        "ADMIN",
	}
	workspaceMember2 := &entity.WorkspaceMember{
		ID:          uuid.NewString(),
		WorkspaceID: workspace.ID,
		UserID:      user.ID,
		Role:        "MEMBER",
	}

	err = workspaceMemberRepo.Create(tx, ctx, workspaceMember1)
	require.NoError(t, err)

	err = workspaceMemberRepo.Create(tx, ctx, workspaceMember2)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrWorkspaceIDUsed.Error())
}

func TestWorkspaceMemberRepo_Create_WorkspaceNotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepo()
	workspaceMemberRepo := NewWorkspaceMemberRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	user := &entity.User{
		ID:       uuid.NewString(),
		Email:    "workspace-member-workspace-not-found@test.com",
		Password: "password",
	}
	err = userRepo.Create(tx, ctx, user)
	require.NoError(t, err)

	workspaceMember := &entity.WorkspaceMember{
		ID:          uuid.NewString(),
		WorkspaceID: uuid.NewString(),
		UserID:      user.ID,
		Role:        "ADMIN",
	}

	err = workspaceMemberRepo.Create(tx, ctx, workspaceMember)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrWorkspaceIDNotFound.Error())
}

func TestWorkspaceMemberRepo_Create_UserNotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepo()
	workspaceRepo := NewWorkspaceRepo()
	workspaceMemberRepo := NewWorkspaceMemberRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	user := &entity.User{
		ID:       uuid.NewString(),
		Email:    "workspace-member-user-not-found@test.com",
		Password: "password",
	}
	err = userRepo.Create(tx, ctx, user)
	require.NoError(t, err)

	workspace := &entity.Workspace{
		ID:          uuid.NewString(),
		OwnerUserID: user.ID,
		Name:        "Workspace Member User Not Found Test",
		Slug:        "workspace-member-user-not-found-test",
		AvatarURL:   "",
	}
	err = workspaceRepo.Create(tx, ctx, workspace)
	require.NoError(t, err)

	workspaceMember := &entity.WorkspaceMember{
		ID:          uuid.NewString(),
		WorkspaceID: workspace.ID,
		UserID:      uuid.NewString(),
		Role:        "ADMIN",
	}

	err = workspaceMemberRepo.Create(tx, ctx, workspaceMember)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrUserIDNotFound.Error())
}
