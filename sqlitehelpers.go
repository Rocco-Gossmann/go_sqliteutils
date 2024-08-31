package go_sqliteutils

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	utils "github.com/rocco-gossmann/go_utils"
	_ "modernc.org/sqlite"
)

const SQL_DATETIME_FORMAT = "2006-01-02 15:04:05"
const SQL_OUTPUT_DATETIMEFORMAT = "2006-01-02T15:04:05Z"

var (
	db           *sql.DB
	hasMetaTable bool = false
)

// Sets up a sqlite db, if none existed before
// @param filename - path to the sqlite-database file
// @param version - version number that is expected from the existing DB File
// @param updateHandler - function, that will be called if the current database version does not match the expected versxion
func InitDBFile(filename string, version uint, updateHandler func(db *sql.DB, isVersion uint, shouldBeVersion uint)) {

	if version == 0 {
		panic("cant initialize db with version 0")
	}

	var err error
	if db != nil {
		return
	}

	_, fierr := os.Stat(filename)
	db, err = sql.Open("sqlite", filename)
	utils.Err(err)

	if os.IsNotExist(fierr) {
		getMetaVersion() // <- creates the _meta table and required entries
		updateHandler(db, 0, version)
		updateMetaVersion(version)
	} else {
		// Check for meta table existence
		tables, err := RowQueryStatement("select name from sqlite_master WHERE type='table' and name='_meta'")
		utils.Err(err)
		var table sql.NullString
		tables.Scan(&table)
		hasMetaTable = table.Valid

		// update version if needed
		isVersion := uint(getMetaVersion())

		if isVersion != version {
			updateHandler(db, isVersion, version)
			updateMetaVersion(version)
		}

	}
}

func DeInitDB() {
	if db != nil {
		db.Close()
		db = nil
		hasMetaTable = false
	}
}

// Runs a prepared statement on the database. Requires the DB to be initilaized first
func ExecStatement(statement string, args ...any) (sql.Result, error) {
	stmt, err := db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}

// Runs a prepared statement on the database. Requires the DB to be initilaized first
func QueryStatement(statement string, args ...any) (*sql.Rows, error) {
	stmt, err := db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(args...)
}

// Runs a prepared statement on the database. Requires the DB to be initilaized first
func RowQueryStatement(statement string, args ...any) (*sql.Row, error) {
	stmt, err := db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args...), nil
}

// @depricated use IsUniqueConstraintError instead.
// (this was a spelling error, just left it in to not break existing things for now)
// will be removed in the future
func IsUniqueContraintError(err error) bool {
	return IsUniqueConstraintError(err)
}

// Checks if an error is related to a Unique Contraint (aka. If the element inserted already exists)
func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "constraint failed: UNIQUE constraint failed")
}

// Helpers
// ==============================================================================

func createMetaTable() {
	if hasMetaTable {
		return
	}

	_, err := ExecStatement("CREATE TABLE IF NOT EXISTS _meta (key TEXT UNIQUE, value TEXT)")
	utils.Err(err)

	hasMetaTable = true
}

func getMetaVersion() int64 {

	createMetaTable()
	row, err := RowQueryStatement("SELECT value FROM _meta WHERE key='version'")
	utils.Err(err)

	var sVersion sql.NullString
	row.Scan(&sVersion)

	if !sVersion.Valid {
		sVersion.String = "0"
		sVersion.Valid = true

		_, err = ExecStatement("INSERT INTO _meta(key, value) VALUES ('version', 0)")
		utils.Err(err)
	}

	i, err := strconv.ParseInt(sVersion.String, 10, 64)
	utils.Err(err)

	return i
}

func updateMetaVersion(newVersion uint) {
	fmt.Println("update version ", newVersion)
	_, err := ExecStatement("UPDATE _meta SET value=? WHERE key='version'", newVersion)
	utils.Err(err)
}

func SQLDateTimePrint(sDateTime string) string {
	oTime, err := time.Parse(SQL_OUTPUT_DATETIMEFORMAT, sDateTime)
	if err != nil {
		return ""
	}

	return oTime.Format(utils.DATETIME_PRINT)
}
