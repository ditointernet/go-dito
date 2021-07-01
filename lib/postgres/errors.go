package postgres

import (
	"github.com/ditointernet/go-dito/lib/errors"
)

var (
	// ErrCantOpenConnection indicates an error on openning connection to the postgres database.
	ErrCantOpenConnection = errors.New("error on openning connection to the postgres database").WithCode("CANT_OPEN_POSTGRES_CONNECTION")
	// ErrDatabaseNotReached indicates that the database can not be reached.
	ErrDatabaseNotReached = errors.New("the database can not be reached").WithCode("CANT_REACH_POSTGRES_DATABASE")
)
