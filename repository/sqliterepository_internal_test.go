package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTable(t *testing.T) {
	r, err := NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	err = r.createTable("Example.guidewire.com")
	assert.NoError(t, err)

	// Create it a second time...
	err = r.createTable("Example.guidewire.com")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "table already exists")

	// Create a table with a reserved name to trigger an error from SQLite
	err = r.createTable("sqlite_Example.guidewire.com")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "object name reserved for internal use")
}
