package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenDB_Success(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	ping := db.Ping()

	assert.Nil(t, ping)
}
