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

func TestUserRepo_Create(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepo()
	ctx := context.Background()

	userData := &entity.User{
		ID:       uuid.NewString(),
		Email:    "duplicate@test.com",
		Password: "",
	}

	tx1, err := db.Begin()
	require.NoError(t, err)
	defer tx1.Rollback()

	err = repo.Create(tx1, ctx, userData)
	assert.NoError(t, err)
}

func TestUserRepo_Create_DuplicateEmail(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepo()
	ctx := context.Background()

	user1 := &entity.User{
		ID:       uuid.NewString(),
		Email:    "duplicate2@test.com",
		Password: "",
	}

	user2 := &entity.User{
		ID:       uuid.NewString(),
		Email:    "duplicate2@test.com",
		Password: "",
	}

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	err = repo.Create(tx, ctx, user1)
	assert.NoError(t, err)

	err = repo.Create(tx, ctx, user2)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrEmailUsed.Error())
}

func TestUserRepo_FindByEmail(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepo()
	ctx := context.Background()

	userData := &entity.User{
		ID:       uuid.NewString(),
		Email:    "find-by-email@test.com",
		Password: "password",
	}

	_, err := db.ExecContext(ctx, `
		INSERT INTO users(id, email, password)
		VALUES($1, $2, $3)
	`, userData.ID, userData.Email, userData.Password)
	require.NoError(t, err)
	defer db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userData.ID)

	user := &entity.User{}
	err = repo.FindByEmail(db, ctx, userData.Email, user)
	require.NoError(t, err)

	assert.Equal(t, userData.ID, user.ID)
	assert.Equal(t, userData.Email, user.Email)
	assert.Equal(t, userData.Password, user.Password)
}

func TestUserRepo_FindByEmail_NotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepo()
	ctx := context.Background()

	user := &entity.User{}
	err := repo.FindByEmail(db, ctx, "find-by-email-not-found@test.com", user)

	assert.ErrorIs(t, err, errs.ErrDataNotFound)
	assert.Empty(t, user.ID)
	assert.Empty(t, user.Email)
	assert.Empty(t, user.Password)
}
