package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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

	container, db, err := createPostgresContainer(ctx, "testdb")
	suite.NoError(err)

	suite.container = container
	suite.db = db

	mig, err := newMigrator(db)
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

// createPostgresContainer creates a new postgres container and returns a database connection.
func createPostgresContainer(ctx context.Context, dbname string) (testcontainers.Container, *sql.DB, error) {
	env := map[string]string{
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_DB":       dbname,
	}

	port := "5432/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:password@%s:%s/%s?sslmode=disable", host, port.Port(), dbname)
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17.0-alpine",
			ExposedPorts: []string{port},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env:          env,
			WaitingFor: wait.ForAll(
				wait.ForSQL(nat.Port(port), "postgres", dbURL).WithStartupTimeout(time.Second*5),
				wait.ForLog("database system is ready to accept connections"),
			),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, fmt.Errorf("failed to start container: %s", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return container, nil, fmt.Errorf("failed to get container external port: %s", err)
	}

	log.Println("postgres container ready and running at port: ", mappedPort)

	url := fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", mappedPort.Port(), dbname)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return container, db, fmt.Errorf("failed to establish database connection: %s", err)
	}

	return container, db, nil
}

// newMigrator creates a new migrator instance for the given database connection.
func newMigrator(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := migrate_postgres.WithInstance(db, &migrate_postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create migrator driver: %s", err)
	}

	return migrate.NewWithDatabaseInstance("file://../../../../../migrations", "postgres", driver)
}
