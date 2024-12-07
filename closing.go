package go_sqliteutils

func (me *DatabaseRessource) Close() (err error) {
	if me == nil {
		panic("can't Close nil")
	}

	me.hasMetaTable = false
	err = me.db.Close()

	if err == nil {
		me.db = nil
	}

	return
}

func (me *DatabaseRessource) Closed() bool {
	return me.db == nil
}
