package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/yourorg/todo-app/services"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB creates a test database using testcontainers
// Returns the database connection and a cleanup function
func SetupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	// Create PostgreSQL container
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		cleanup()
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Connect using GORM
	db, err := gorm.Open(postgresdriver.Open(connStr), &gorm.Config{})
	if err != nil {
		cleanup()
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := services.AutoMigrate(db); err != nil {
		cleanup()
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db, cleanup
}

// TruncateTables truncates the specified tables in reverse order
func TruncateTables(db *gorm.DB, tables ...string) {
	// Truncate in reverse order (children before parents)
	for i := len(tables) - 1; i >= 0; i-- {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tables[i]))
	}
}