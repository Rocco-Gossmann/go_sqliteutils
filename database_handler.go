package go_sqliteutils

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"

	_ "modernc.org/sqlite"
)

type DatabaseRessource struct {
	db           *sql.DB
	hasMetaTable bool
	err          error
	dbfile       string
}

// Sets up a sqlite db, if none existed before
// @param filename - path to the sqlite-database file
// @param version - version number that is expected from the existing DB File
// @param updateHandler - function, that will be called if the current database version does not match the expected versxion
func (db *DatabaseRessource) InitFromFile(
	filename string,
	version uint,
	updateHandler func(tx *sql.Tx, isVersion uint, shouldBeVersion uint) error,
) (err error) {

	defer func() {
		if err != nil {
			log.Println("Error: ", err)
		}
	}()

	if db == nil {
		panic("can't initialize nil")
	}

	db.dbfile = filename

	_, fierr := os.Stat(filename)

	db.db, err = sql.Open("sqlite", filename)
	db.db.SetMaxOpenConns(1)

	var tx *sql.Tx
	tx, err = db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err != nil {
		return err
	}

	if os.IsNotExist(fierr) {
		err = db.createMetaTable(tx) // <- creates the _meta table and required entries

		var isVersion int64
		if isVersion, err = db.getMetaVersion(tx); err == nil {
			err = updateHandler(tx, uint(isVersion), version)
		}

		if err == nil {
			err = db.updateMetaVersion(tx, version)
		}

		if err != nil {
			return
		}

	} else {

		// Check for meta table existence
		var table sql.NullString
		var isVersion int64

		err = tx.QueryRow("select name from sqlite_master WHERE type='table' and name='_meta'").Scan(&table)
		if err != nil {
			return err
		}

		if !table.Valid {
			if err = db.createMetaTable(tx); err != nil {
				return
			}

			err = tx.QueryRow("select name from sqlite_master WHERE type='table' and name='_meta'").Scan(&table)
			if err != nil {
				return
			}

		}

		if !table.Valid {
			return errors.New("could find find or create _meta table")
		}

		db.hasMetaTable = true

		isVersion, err = db.getMetaVersion(tx)

		if err == nil && uint(isVersion) != version {
			err = updateHandler(tx, uint(isVersion), version)
			if err == nil {
				err = db.updateMetaVersion(tx, version)
			}
			if err != nil {
				return
			}
		}
	}

	err = tx.Commit()
	return
}

func (db *DatabaseRessource) createMetaTable(tx *sql.Tx) (err error) {
	if db.hasMetaTable {
		return
	}

	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS _meta (key TEXT UNIQUE, value TEXT)")

	if err == nil {
		db.hasMetaTable = true
	}

	return

}

func (db *DatabaseRessource) getMetaVersion(tx *sql.Tx) (v int64, err error) {

	var sVersion sql.NullString

	if db.createMetaTable(tx) != nil {
		return
	}

	err = tx.
		QueryRow("SELECT value FROM _meta WHERE key=?", VERSION_KEY).
		Scan(&sVersion)

	if !sVersion.Valid {
		sVersion.String = "0"
		sVersion.Valid = true

		_, err = tx.Exec("INSERT INTO _meta(key, value) VALUES (?, 0)", VERSION_KEY)
		if err != nil {
			return
		}
	}

	v, err = strconv.ParseInt(sVersion.String, 10, 64)

	return
}

func (db *DatabaseRessource) updateMetaVersion(tx *sql.Tx, newVersion uint) error {
	_, err := tx.Exec("UPDATE _meta SET value=? WHERE key=?", newVersion, VERSION_KEY)
	return err
}
