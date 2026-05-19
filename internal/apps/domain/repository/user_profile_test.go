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

func TestUserProfileRepo_Create(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepo()
	userProfileRepo := NewUserProfileRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	user := &entity.User{
		ID:       uuid.NewString(),
		Email:    "profile-create@test.com",
		Password: "password",
	}

	err = userRepo.Create(tx, ctx, user)
	assert.NoError(t, err)

	userProfile := &entity.UserProfile{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Name:      "Test User",
		AvatarURL: "",
	}

	err = userProfileRepo.Create(tx, ctx, userProfile)
	assert.NoError(t, err)
}

func TestUserProfileRepo_Create_DuplicateUserID(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepo()
	userProfileRepo := NewUserProfileRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	user := &entity.User{
		ID:       uuid.NewString(),
		Email:    "profile-duplicate-user-id@test.com",
		Password: "password",
	}

	err = userRepo.Create(tx, ctx, user)
	assert.NoError(t, err)

	userProfile1 := &entity.UserProfile{
		ID:     uuid.NewString(),
		UserID: user.ID,
		Name:   "Test User 1",
	}

	userProfile2 := &entity.UserProfile{
		ID:     uuid.NewString(),
		UserID: user.ID,
		Name:   "Test User 2",
	}

	err = userProfileRepo.Create(tx, ctx, userProfile1)
	assert.NoError(t, err)

	err = userProfileRepo.Create(tx, ctx, userProfile2)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrUserIDUsed.Error())
}

func TestUserProfileRepo_Create_UserIDNotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	userProfileRepo := NewUserProfileRepo()
	ctx := context.Background()

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	userProfile := &entity.UserProfile{
		ID:     uuid.NewString(),
		UserID: uuid.NewString(),
		Name:   "Test User",
	}

	err = userProfileRepo.Create(tx, ctx, userProfile)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrUserIDNotFound.Error())
}

func TestUserProfileRepo_FindByUserID(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	userID := uuid.NewString()
	profileID := uuid.NewString()

	_, err := db.ExecContext(ctx, `
		INSERT INTO users(id, email, password)
		VALUES($1, $2, $3)
	`, userID, "profile-find-by-user-id@test.com", "password")
	require.NoError(t, err)
	defer db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)

	_, err = db.ExecContext(ctx, `
		INSERT INTO user_profiles(id, user_id, name, avatar_url)
		VALUES($1, $2, $3, $4)
	`, profileID, userID, "Find Profile", "https://example.com/avatar.png")
	require.NoError(t, err)

	repo := NewUserProfileRepo()
	userProfile := &entity.UserProfile{}
	err = repo.FindByUserID(db, ctx, userID, userProfile)
	require.NoError(t, err)

	assert.Equal(t, profileID, userProfile.ID)
	assert.Equal(t, userID, userProfile.UserID)
	assert.Equal(t, "Find Profile", userProfile.Name)
	assert.Equal(t, "https://example.com/avatar.png", userProfile.AvatarURL)
}

func TestUserProfileRepo_FindByUserID_NotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserProfileRepo()
	ctx := context.Background()

	userProfile := &entity.UserProfile{}
	err := repo.FindByUserID(db, ctx, uuid.NewString(), userProfile)

	assert.ErrorIs(t, err, errs.ErrDataNotFound)
	assert.Empty(t, userProfile.ID)
	assert.Empty(t, userProfile.UserID)
	assert.Empty(t, userProfile.Name)
	assert.Empty(t, userProfile.AvatarURL)
}
