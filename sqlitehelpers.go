package go_sqliteutils

import (
	"database/sql"
	"strings"
	"time"
)

const SQL_DATETIME_FORMAT = "2006-01-02 15:04:05"
const SQL_OUTPUT_DATETIMEFORMAT = "2006-01-02T15:04:05Z"

// Update Handler for when you just want to use the Key=>ValueStore and don't need
// any further Tables
func NoUpdate(tx *sql.Tx, isVersion uint, shouldBeVersion uint) error {
	return nil
}

// Sets up a sqlite db, if none existed before
// @param filename - path to the sqlite-database file
// @param version - version number that is expected from the existing DB File
// @param updateHandler - function, that will be called if the current database version does not match the expected versxion
func Open(
	filename string,
	version uint,
	updateHandler func(tx *sql.Tx, isVersion uint, shouldBeVersion uint) error,
) (DBRessource, error) {

	if version == 0 {
		panic("cant initialize db with version 0")
	}

	db := DatabaseRessource{}
	err := db.InitFromFile(filename, version, updateHandler)

	return &db, err
}

// Helpers
// ==============================================================================
func SQLDateTimePrint(sDateTime string) string {
	oTime, err := time.Parse(SQL_OUTPUT_DATETIMEFORMAT, sDateTime)
	if err != nil {
		return ""
	}

	return oTime.Format("02. Jan. 2006 15:04")
}

// Checks if an error is related to a Unique Contraint (aka. If the element inserted already exists)
func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "constraint failed: UNIQUE constraint failed")
}
