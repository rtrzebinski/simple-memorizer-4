package web

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
func (suite *PostgresSuite) SetupSuite() {
	ctx := context.Background()

	container, db, err := postgres.CreatePostgresContainer(ctx, "testdb")
	suite.NoError(err)

	suite.container = container
	suite.db = db

	mig, err := postgres.NewMigrator(db)
	suite.NoError(err)

	suite.migrator = mig
}

// TearDownSuite runs after all the tests in the suite.
func (suite *PostgresSuite) TearDownSuite() {
	err := suite.db.Close()
	suite.NoError(err)

	err = suite.container.Terminate(context.Background())
	suite.NoError(err)
}

// SetupTest runs before each test in the suite.
func (suite *PostgresSuite) SetupTest() {
	suite.NoError(suite.migrator.Up())
}

// TearDownTest runs after each test in the suite.
func (suite *PostgresSuite) TearDownTest() {
	suite.NoError(suite.migrator.Down())
}

// TestPostgresSuite runs all test in the suite.
func TestPostgresSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Parallel()

	suite.Run(t, new(PostgresSuite))
}
