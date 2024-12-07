package go_sqliteutils

import "database/sql"

func (db *DatabaseRessource) Begin() (*sql.Tx, error) {

	if db == nil {
		panic("can't call Begin on nil")
	}

	if db.db == nil {
		panic("can't Begin a transaction for an already closed database")
	}

	return db.db.Begin()

}
