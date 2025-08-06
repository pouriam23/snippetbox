package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	// Connect to the test DB.
	// Use "multiStatements=true" to allow running multiple SQL statements in one Exec() call.
	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// Read and run the setup.sql script.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Register a cleanup function — it will automatically run after the test ends.
	t.Cleanup(func() {
		// Read and run the teardown.sql script.
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		// Close the DB connection.
		db.Close()
	})

	// Return the test DB connection pool.
	return db
}
