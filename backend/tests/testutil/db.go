package testutil

import (
	"database/sql"
	"os"
	"testing"

	"github.com/flagmanagment/backend/internal/repository"
	_ "github.com/lib/pq"
)

// GetTestDB connects to the test database and returns the *sql.DB and repository.Store.
// If the database is not accessible, the test is skipped gracefully unless REQUIRE_TEST_DB=true.
func GetTestDB(t *testing.T) (*sql.DB, repository.Store) {
	t.Helper()

	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/flagmanagment?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		if os.Getenv("REQUIRE_TEST_DB") == "true" {
			t.Fatalf("failed to open database: %v", err)
		}
		t.Skipf("skipping test requiring database: %v", err)
		return nil, nil
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		if os.Getenv("REQUIRE_TEST_DB") == "true" {
			t.Fatalf("failed to ping database: %v", err)
		}
		t.Skipf("skipping test: database unreachable at %s: %v", connStr, err)
		return nil, nil
	}

	store := repository.NewStore(db)
	return db, store
}

