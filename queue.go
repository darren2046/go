package golanglibs

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type queueStruct struct {
	db      *leveldb.DB
	datadir string
	closed  bool
}

func getQueue(datadir string) (q *queueStruct) {
	q = &queueStruct{}

	db, err := leveldb.OpenFile(datadir, nil)
	Panicerr(err)

	q.db = db
	q.datadir = datadir

	return
}

func (m *queueStruct) Close() {
	err := m.db.Close()
	Panicerr(err)

	m.closed = true
}

func (m *queueStruct) Destroy() {
	if !m.closed {
		m.Close()
	}
	Os.Unlink(m.datadir)
}

type namedQueueStruct struct {
	head  int64
	tail  int64
	db    *leveldb.DB
	glock *LockStruct
	plock *LockStruct
	name  string
}

// Will not clean the data already exists
func (m *queueStruct) New(queueName ...string) *namedQueueStruct {
	q := &namedQueueStruct{}

	n := ""
	if len(queueName) != 0 {
		n = queueName[0]
	}

	status, err := m.db.Has([]byte(n+"_head"), nil)
	Panicerr(err)
	if status {
		head, err := m.db.Get([]byte(n+"_head"), nil)
		Panicerr(err)
		q.head = Int64(Str(head))
	}

	status, err = m.db.Has([]byte(n+"_tail"), nil)
	Panicerr(err)
	if status {
		tail, err := m.db.Get([]byte(n+"_tail"), nil)
		Panicerr(err)
		q.tail = Int64(Str(tail))
	}

	q.db = m.db
	q.glock = Tools.Lock()
	q.plock = Tools.Lock()
	q.name = n

	return q
}

func (m *namedQueueStruct) Size() int64 {
	return m.tail - m.head
}

func (m *namedQueueStruct) Get(nonblock ...bool) string {
	m.glock.Acquire()
	defer m.glock.Release()

	if len(nonblock) != 0 && nonblock[0] {
		return ""
	}

	for m.head == m.tail {
		Time.Sleep(0.1)
	}
	value, err := m.db.Get([]byte(Str(m.head)), nil)
	Panicerr(err)

	err = m.db.Delete([]byte(Str(m.head)), nil)
	Panicerr(err)

	m.head += 1

	err = m.db.Put([]byte(m.name+"_head"), []byte(Str(m.head)), nil)
	Panicerr(err)

	return Str(value)
}

func (m *namedQueueStruct) Put(value string) {
	if value == "" {
		Panicerr("value can not be empty")
	}

	m.plock.Acquire()
	defer m.plock.Release()

	err := m.db.Put([]byte(Str(m.tail)), []byte(value), nil)
	Panicerr(err)

	m.tail += 1

	err = m.db.Put([]byte(m.name+"_tail"), []byte(Str(m.tail)), nil)
	Panicerr(err)
}
