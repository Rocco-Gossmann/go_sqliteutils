package go_sqliteutils

import (
	"errors"
	"log"
)

const VERSION_KEY = "2af72f100c356273d46284f6fd1dfc08"
const SQLITE_NO_ROWS = "sql: no rows in result set"

var keyBlackList = map[string]struct{}{
	VERSION_KEY: {},
}

var MetaKeyAccessViolation = errors.New("Tried to access a reserved key")
var noResultsForKey = errors.New("No results for this Key")

func (db *DatabaseRessource) IsNoResultForKey(err error) bool {
	return err == noResultsForKey
}

func (db *DatabaseRessource) SetValue(key, value string) error {

	if _, ok := keyBlackList[key]; ok {
		return MetaKeyAccessViolation
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM _meta WHERE key=?", key)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.Exec("INSERT INTO _meta(key, value) VALUES (?, ?)", key, value)

	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	log.Println("go_sqliteutils: COMMIT:", err)
	return err

}

func (db *DatabaseRessource) GetValue(key string) (val string, err error) {
	if _, ok := keyBlackList[key]; ok {
		return "", MetaKeyAccessViolation
	}

	err = db.db.QueryRow("SELECT value FROM _meta WHERE key = ?", key).Scan(&val)
	if err != nil && err.Error() == SQLITE_NO_ROWS {
		err = noResultsForKey
	}

	return
}

func (db *DatabaseRessource) DropValue(key string) error {
	if _, ok := keyBlackList[key]; ok {
		return MetaKeyAccessViolation
	}

	_, err := db.db.Exec("DELETE FROM _meta WHERE mame = ?", key)
	if err != nil {
		return err
	}

	return nil
}
