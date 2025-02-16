package auth

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

// PostgresSuite is a test suite to be used with postgres integration tests.
type PostgresSuite struct {
	suite.Suite

	container testcontainers.Container
	db        *sql.DB
	migrator  *migrate.Migrate
}

// SetupSuite runs before all the tests in the suite.
func (s *PostgresSuite) SetupSuite() {
	ctx := context.Background()

	container, db, err := postgres.CreatePostgresContainer(ctx, "testdb")
	s.NoError(err)

	s.container = container
	s.db = db

	mig, err := postgres.NewMigrator(db)
	s.NoError(err)

	s.migrator = mig
}

// TearDownSuite runs after all the tests in the suite.
func (s *PostgresSuite) TearDownSuite() {
	err := s.db.Close()
	s.NoError(err)

	err = s.container.Terminate(context.Background())
	s.NoError(err)
}

// SetupTest runs before each test in the suite.
func (s *PostgresSuite) SetupTest() {
	s.NoError(s.migrator.Up())
}

// TearDownTest runs after each test in the suite.
func (s *PostgresSuite) TearDownTest() {
	s.NoError(s.migrator.Down())
}

// TestPostgresSuite runs all test in the suite.
func TestPostgresSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Parallel()

	suite.Run(t, new(PostgresSuite))
}
