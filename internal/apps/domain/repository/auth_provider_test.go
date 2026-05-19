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

func TestAuthProviderRepo_FindByUserID(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	userID := uuid.NewString()
	authProviderID := uuid.NewString()

	_, err := db.ExecContext(ctx, `
		INSERT INTO users(id, email, password)
		VALUES($1, $2, $3)
	`, userID, "auth-provider-find-by-user-id@test.com", "password")
	require.NoError(t, err)
	defer db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)

	_, err = db.ExecContext(ctx, `
		INSERT INTO auth_providers(id, user_id, provider_name)
		VALUES($1, $2, $3)
	`, authProviderID, userID, "local")
	require.NoError(t, err)

	repo := NewAuthProviderRepo()
	authProvider := &entity.AuthProvider{}
	err = repo.FindByUserID(db, ctx, userID, authProvider)
	require.NoError(t, err)

	assert.Equal(t, authProviderID, authProvider.ID)
	assert.Equal(t, userID, authProvider.UserID)
	assert.Equal(t, "local", authProvider.ProviderName)
}

func TestAuthProviderRepo_FindByUserID_NotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewAuthProviderRepo()
	ctx := context.Background()

	authProvider := &entity.AuthProvider{}
	err := repo.FindByUserID(db, ctx, uuid.NewString(), authProvider)

	assert.ErrorIs(t, err, errs.ErrDataNotFound)
	assert.Empty(t, authProvider.ID)
	assert.Empty(t, authProvider.UserID)
	assert.Empty(t, authProvider.ProviderName)
}
