package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	mpg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Suite is a test suite that manages a PostgreSQL testcontainer and handles database migrations.
// Use it as an embedded struct in your test suites.
type Suite struct {
	suite.Suite

	Container testcontainers.Container
	DB        *sql.DB
	Migrator  *migrate.Migrate
}

var (
	startOnce  sync.Once
	startErr   error
	sharedCont testcontainers.Container
	sharedDB   *sql.DB
	sharedMigr *migrate.Migrate
)

func (s *Suite) SetupSuite() {
	ctx := s.T().Context()

	startOnce.Do(func() {
		c, db, err := CreatePostgresContainer(ctx)
		if err != nil {
			startErr = err
			return
		}
		m, err := NewMigrator(db)
		if err != nil {
			_ = c.Terminate(ctx)
			startErr = err
			return
		}
		sharedCont = c
		sharedDB = db
		sharedMigr = m
	})
	s.Require().NoError(startErr)

	s.Container = sharedCont
	s.DB = sharedDB
	s.Migrator = sharedMigr
}

func (s *Suite) SetupTest() {
	if s.Migrator == nil {
		s.T().Fatal("migrator is nil in SetupTest")
	}
	if err := s.Migrator.Up(); err != nil && err != migrate.ErrNoChange {
		s.T().Fatalf("migration up failed: %v", err)
	}
}

func (s *Suite) TearDownTest() {
	if s.Migrator == nil {
		return
	}
	if err := s.Migrator.Down(); err != nil && err != migrate.ErrNoChange {
		s.T().Fatalf("migration down failed: %v", err)
	}
}

func CreatePostgresContainer(ctx context.Context) (testcontainers.Container, *sql.DB, error) {
	var (
		dbname = "postgres"
		port   = nat.Port("5432/tcp")
		dbURL  = func(host string, port nat.Port) string {
			return fmt.Sprintf("postgres://postgres:password@%s:%s/%s?sslmode=disable", host, port.Port(), dbname)
		}
	)

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17.0-alpine",
			ExposedPorts: []string{string(port)},
			Env: map[string]string{
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_DB":       dbname,
			},
			WaitingFor: wait.ForAll(
				wait.ForSQL(port, "postgres", dbURL).WithStartupTimeout(2*time.Minute),
				wait.ForLog("database system is ready to accept connections"),
			),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("create a generic container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("retrieve host from container: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, port)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("retrieve mapped port from container: %w", err)
	}

	db, err := sql.Open("postgres", dbURL(host, mappedPort))
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("open database connection: %w", err)
	}

	return container, db, nil
}

func NewMigrator(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := mpg.WithInstance(db, &mpg.Config{})
	if err != nil {
		return nil, fmt.Errorf("create migrator driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://../../../migrations", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("create a new migrator: %w", err)
	}

	return migrator, nil
}
