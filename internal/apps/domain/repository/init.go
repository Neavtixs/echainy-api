package repository

import (
	"database/sql"
	"testing"

	"github.com/Neavtixs/go-backend-template/configs"
)

func SetupTestDB(t *testing.T) *sql.DB {
	configs.LoadEnv("../../../../.env")
	db := configs.GetConnection()

	return db
}
