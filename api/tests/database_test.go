package tests

import (
	"testing"

	"github.com/chadzink/skills-api/database"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	dbConnectError := database.ConnectTestDb()
	if dbConnectError != nil {
		t.Error(dbConnectError)
	}

	// Perform a basic database operation (e.g., create a table)
	err := database.DAL.Db.Exec("CREATE TABLE test_table (id SERIAL PRIMARY KEY, name VARCHAR);").Error
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Perform a query to ensure the database connection is working
	var count int64
	database.DAL.Db.Table("test_table").Count(&count)
	assert.Equal(t, int64(0), count, "Expected count to be 0")

	// Optionally, clean up by dropping the test table
	database.DAL.Db.Exec("DROP TABLE test_table;")
}
