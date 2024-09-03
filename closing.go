package go_sqliteutils

func (db *DatabaseRessource) Close() error {
	if db == nil {
		panic("can't Close nil")
	}

	db.hasMetaTable = false
	err := db.db.Close()

	if err != nil {
		return err
	} else {
		db.db = nil
	}

	return nil
}
