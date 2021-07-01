package postgres

import (
	"database/sql"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"

	"github.com/ditointernet/go-dito/lib/env"
	"github.com/ditointernet/go-dito/lib/errors"

	// Local Postgres driver
	_ "github.com/lib/pq"
)

// postgresDriver is a type of DB driver for postgres.
type postgresDriver int

const (
	// postgresDriverDefault is the generic Postgres driver.
	postgresDriverDefault postgresDriver = iota
	// postgresDriverDefault is the CloudSQL's Postgres driver.
	postgresDriverCloudSQL
)

var validPostgresDrivers = []string{
	"postgres",
	"cloudsqlpostgres",
}

// String is a string representation of the PostgresDriver.
func (d postgresDriver) String() string {
	return validPostgresDrivers[d]
}

// ClientParams encapsulates all Client's dependencies.
type ClientParams struct {
	Tracer         trace.Tracer
	Environment    env.Environment
	MaxConnections int
	DatabaseURI    string
}

// Client is a custom postgres client.
type Client struct {
	tracer      trace.Tracer
	maxConn     int
	driver      postgresDriver
	databaseURI string
}

// NewClient creates a new Client instance.
// It uses postgres generic driver by default. If in production environment, uses cloudsqlpostgres driver.
// It uses 5 max connections by default.
func NewClient(params ClientParams) (Client, error) {
	if params.DatabaseURI == "" {
		return Client{}, errors.NewMissingRequiredDependency("DatabaseURI")
	}

	if params.MaxConnections < 1 {
		params.MaxConnections = 5
	}

	var driver postgresDriver
	if params.Environment == env.ProductionEnvironment {
		driver = postgresDriverCloudSQL
	}

	return Client{
		tracer:      params.Tracer,
		maxConn:     params.MaxConnections,
		driver:      driver,
		databaseURI: params.DatabaseURI,
	}, nil
}

// MustNewClient creates a new Client instance.
// It uses postgres generic driver by default. If in production environment, uses cloudsqlpostgres driver.
// It uses 5 max connections by default.
// It panics if any error is found.
func MustNewClient(params ClientParams) Client {
	cli, err := NewClient(params)
	if err != nil {
		panic(err)
	}

	return cli
}

// Connect connects the client with the database.
func (c Client) Connect(ctx context.Context) (*sql.DB, error) {
	if c.tracer != nil {
		var span trace.Span
		ctx, span = c.tracer.Start(ctx, "postgres.Client.Connect")
		defer span.End()
	}

	dbConn, err := sql.Open(c.driver.String(), c.databaseURI)
	if err != nil {
		return nil, ErrCantOpenConnection
	}

	if err = dbConn.Ping(); err != nil {
		return nil, ErrDatabaseNotReached
	}

	dbConn.SetMaxOpenConns(c.maxConn)
	return dbConn, nil
}
