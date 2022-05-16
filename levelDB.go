package golanglibs

import "github.com/syndtr/goleveldb/leveldb"

type LevelDBStruct struct {
	closed  bool
	db      *leveldb.DB
	datadir string
}

func getLevelDB(datadir string) (l *LevelDBStruct) {
	l = &LevelDBStruct{}

	db, err := leveldb.OpenFile(datadir, nil)
	Panicerr(err)

	l.db = db
	l.datadir = datadir

	return
}

func (m *LevelDBStruct) Close() {
	if !m.closed {
		err := m.db.Close()
		Panicerr(err)

		m.closed = true
	}
}

func (m *LevelDBStruct) Destroy() {
	if !m.closed {
		m.Close()
	}
	Os.Unlink(m.datadir)
}

func (m *LevelDBStruct) Exists(key string) bool {
	if m.closed {
		Panicerr("Database is closed.")
	}

	status, err := m.db.Has([]byte(key), nil)
	Panicerr(err)
	return status
}

func (m *LevelDBStruct) Get(key string) string {
	if m.closed {
		Panicerr("Database is closed.")
	}
	if m.Exists(key) {
		value, err := m.db.Get([]byte(key), nil)
		Panicerr(err)
		return Str(value)
	} else {
		Panicerr("Key \"" + key + "\" not exists.")
	}
	return ""
}

func (m *LevelDBStruct) Set(key string, value string) {
	if m.closed {
		Panicerr("Database is closed.")
	}

	err := m.db.Put([]byte(key), []byte(value), nil)
	Panicerr(err)
}

func (m *LevelDBStruct) Delete(key string) {
	if m.closed {
		Panicerr("Database is closed.")
	}

	if m.Exists(key) {
		err := m.db.Delete([]byte(key), nil)
		Panicerr(err)
	}
}
