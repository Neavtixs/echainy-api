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

func TestAuthProviderRepo_Create(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepo()
	authProviderRepo := NewAuthProviderRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	user := &entity.User{
		ID:       uuid.NewString(),
		Email:    "auth-provider-create@test.com",
		Password: "password",
	}

	err = userRepo.Create(tx, ctx, user)
	assert.NoError(t, err)

	authProvider := &entity.AuthProvider{
		ID:           uuid.NewString(),
		UserID:       user.ID,
		ProviderName: "local",
	}

	err = authProviderRepo.Create(tx, ctx, authProvider)
	assert.NoError(t, err)
}

func TestAuthProviderRepo_Create_DuplicateUserID(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepo()
	authProviderRepo := NewAuthProviderRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	user := &entity.User{
		ID:       uuid.NewString(),
		Email:    "auth-provider-duplicate-user-id@test.com",
		Password: "password",
	}

	err = userRepo.Create(tx, ctx, user)
	assert.NoError(t, err)

	authProvider1 := &entity.AuthProvider{
		ID:           uuid.NewString(),
		UserID:       user.ID,
		ProviderName: "local",
	}

	authProvider2 := &entity.AuthProvider{
		ID:           uuid.NewString(),
		UserID:       user.ID,
		ProviderName: "google",
	}

	err = authProviderRepo.Create(tx, ctx, authProvider1)
	assert.NoError(t, err)

	err = authProviderRepo.Create(tx, ctx, authProvider2)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrUserIDUsed.Error())
}

func TestAuthProviderRepo_Create_UserIDNotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	authProviderRepo := NewAuthProviderRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	authProvider := &entity.AuthProvider{
		ID:           uuid.NewString(),
		UserID:       uuid.NewString(),
		ProviderName: "local",
	}

	err = authProviderRepo.Create(tx, ctx, authProvider)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrUserIDNotFound.Error())
}
