package tests

import (
	"github.com/chadzink/skills-api/database"
	"github.com/stretchr/testify/assert"
)

func (suite *TestWithDbSuite) TestDatabaseConnection() {
	// Used to count results for assertions
	var count int64

	// Check if the seed data is loaded for skills
	database.DAL.Db.Table("skills").Count(&count)
	assert.Equal(suite.T(), int64(4), count, "Expected count to be 4")

	// Check if the seed data is loaded for categories
	database.DAL.Db.Table("categories").Count(&count)
	assert.Equal(suite.T(), int64(4), count, "Expected count to be 4")

	// Check if the seed data is loaded for people
	database.DAL.Db.Table("people").Count(&count)
	assert.Equal(suite.T(), int64(3), count, "Expected count to be 3")

	// Perform a basic database operation (e.g., create a table)
	err := database.DAL.Db.Exec("CREATE TABLE test_table (id SERIAL PRIMARY KEY, name VARCHAR);").Error
	if err != nil {
		suite.T().Fatalf("Failed to create table: %v", err)
	}

	// Perform a query to ensure the database connection is working
	database.DAL.Db.Table("test_table").Count(&count)
	assert.Equal(suite.T(), int64(0), count, "Expected count to be 0")

	// Optionally, clean up by dropping the test table
	database.DAL.Db.Exec("DROP TABLE test_table;")
}
