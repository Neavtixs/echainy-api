package repository

import (
	"context"
	"testing"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/entity"
	"github.com/Neavtixs/echainy-api/internal/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

	tx1, _ := db.Begin()
	defer tx1.Rollback()

	err := repo.Create(tx1, ctx, userData)
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

	tx, _ := db.Begin()
	defer tx.Rollback()

	err := repo.Create(tx, ctx, user1)
	assert.NoError(t, err)

	err = repo.Create(tx, ctx, user2)
	assert.Error(t, err)
	assert.EqualError(t, err, errs.ErrEmailUsed.Error())
}
