package go_sqliteutils

import "database/sql"

func (db *DatabaseRessource) Begin() (*sql.Tx, error) {

	if db == nil {
		panic("can't call Begin on nil")
	}

	return db.db.Begin()

}
