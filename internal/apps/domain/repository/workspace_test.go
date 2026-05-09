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

func TestWorkspaceRepo_Create(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepo()
	workspaceRepo := NewWorkspaceRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	user := &entity.User{
		ID:       uuid.NewString(),
		Email:    "workspace-create@test.com",
		Password: "password",
	}

	err = userRepo.Create(tx, ctx, user)
	assert.NoError(t, err)

	workspace := &entity.Workspace{
		ID:          uuid.NewString(),
		OwnerUserID: user.ID,
		Name:        "Test Workspace",
		Slug:        "test-workspace",
		AvatarURL:   "",
	}

	err = workspaceRepo.Create(tx, ctx, workspace)
	assert.NoError(t, err)
}

func TestWorkspaceRepo_Create_SameSlugDifferentOwner(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepo()
	workspaceRepo := NewWorkspaceRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	user1 := &entity.User{
		ID:       uuid.NewString(),
		Email:    "workspace-duplicate-slug-1@test.com",
		Password: "password",
	}

	user2 := &entity.User{
		ID:       uuid.NewString(),
		Email:    "workspace-duplicate-slug-2@test.com",
		Password: "password",
	}

	err = userRepo.Create(tx, ctx, user1)
	assert.NoError(t, err)

	err = userRepo.Create(tx, ctx, user2)
	assert.NoError(t, err)

	workspace1 := &entity.Workspace{
		ID:          uuid.NewString(),
		OwnerUserID: user1.ID,
		Name:        "Test Workspace 1",
		Slug:        "duplicate-workspace",
	}

	workspace2 := &entity.Workspace{
		ID:          uuid.NewString(),
		OwnerUserID: user2.ID,
		Name:        "Test Workspace 2",
		Slug:        "duplicate-workspace",
	}

	err = workspaceRepo.Create(tx, ctx, workspace1)
	assert.NoError(t, err)

	err = workspaceRepo.Create(tx, ctx, workspace2)
	assert.NoError(t, err)
}

func TestWorkspaceRepo_Create_SameSlugSameOwner(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepo()
	workspaceRepo := NewWorkspaceRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	user := &entity.User{
		ID:       uuid.NewString(),
		Email:    "workspace-same-slug-same-owner@test.com",
		Password: "password",
	}

	err = userRepo.Create(tx, ctx, user)
	assert.NoError(t, err)

	workspace1 := &entity.Workspace{
		ID:          uuid.NewString(),
		OwnerUserID: user.ID,
		Name:        "Test Workspace 1",
		Slug:        "same-owner-duplicate-workspace",
	}

	workspace2 := &entity.Workspace{
		ID:          uuid.NewString(),
		OwnerUserID: user.ID,
		Name:        "Test Workspace 2",
		Slug:        "same-owner-duplicate-workspace",
	}

	err = workspaceRepo.Create(tx, ctx, workspace1)
	assert.NoError(t, err)

	err = workspaceRepo.Create(tx, ctx, workspace2)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrSlugUsed.Error())
}

func TestWorkspaceRepo_Create_UserIDNotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	workspaceRepo := NewWorkspaceRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	workspace := &entity.Workspace{
		ID:          uuid.NewString(),
		OwnerUserID: uuid.NewString(),
		Name:        "Test Workspace",
		Slug:        "workspace-without-owner",
	}

	err = workspaceRepo.Create(tx, ctx, workspace)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrUserIDNotFound.Error())
}
