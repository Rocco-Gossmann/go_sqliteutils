package go_sqliteutils

import (
	"database/sql"
)

type DBRessource interface {

	// Closes the Database connection
	Close() error

	// Stores a Key=>Value Pair in the _meta Table
	SetValue(key, value string) error
	// retreives the value of a Key=>Value Pair in the _meta Table
	GetValue(key string) (val string, err error)
	// removes the value of a Key=>Value Pair in the _meta Table
	DropValue(key string) error

	// Starts an SQL-Transaction, and allows for manipulating Tables
	Begin() (*sql.Tx, error)
}
