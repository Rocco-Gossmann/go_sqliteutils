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

	// Defines if an error was returned, because a GetValue was requested for
	// an item not in the Database
	IsNoResultForKey(err error) bool

	// Starts an SQL-Transaction, and allows for manipulating Tables
	Begin() (*sql.Tx, error)
}
