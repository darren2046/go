package golanglibs

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type QueueStruct struct {
	db      *leveldb.DB
	datadir string
	closed  bool
	nqrv    map[string]*NamedQueueRuntimeVarsStruct
}

type NamedQueueRuntimeVarsStruct struct {
	head int64
	tail int64
	lock *LockStruct
}

func getQueue(datadir string) (q *QueueStruct) {
	q = &QueueStruct{}

	db, err := leveldb.OpenFile(datadir, nil)
	Panicerr(err)

	q.db = db
	q.datadir = datadir
	q.nqrv = make(map[string]*NamedQueueRuntimeVarsStruct)

	return
}

func (m *QueueStruct) Close() {
	err := m.db.Close()
	Panicerr(err)

	m.closed = true
}

func (m *QueueStruct) Destroy() {
	if !m.closed {
		m.Close()
	}
	Os.Unlink(m.datadir)
}

type NamedQueueStruct struct {
	db   *leveldb.DB
	name string
	tq   *QueueStruct // Top Queue
}

// Will not clean the data already exists
func (m *QueueStruct) New(queueName ...string) *NamedQueueStruct {
	q := &NamedQueueStruct{}

	n := ""
	if len(queueName) != 0 {
		n = queueName[0]
	}

	if !Map(m.nqrv).Has(n) {
		m.nqrv[n] = &NamedQueueRuntimeVarsStruct{}
	}

	status, err := m.db.Has([]byte(n+"_head"), nil)
	Panicerr(err)
	if status {
		head, err := m.db.Get([]byte(n+"_head"), nil)
		Panicerr(err)
		m.nqrv[n].head = Int64(Str(head))
	}

	status, err = m.db.Has([]byte(n+"_tail"), nil)
	Panicerr(err)
	if status {
		tail, err := m.db.Get([]byte(n+"_tail"), nil)
		Panicerr(err)
		m.nqrv[n].tail = Int64(Str(tail))
	}

	m.nqrv[n].lock = Tools.Lock()

	q.db = m.db
	q.name = n
	q.tq = m

	return q
}

func (m *NamedQueueStruct) Size() int64 {
	m.tq.nqrv[m.name].lock.Acquire()
	defer m.tq.nqrv[m.name].lock.Release()

	return m.tq.nqrv[m.name].tail - m.tq.nqrv[m.name].head
}

func (m *NamedQueueStruct) Get(nonblock ...bool) string {
	m.tq.nqrv[m.name].lock.Acquire()
	defer m.tq.nqrv[m.name].lock.Release()

	if m.tq.nqrv[m.name].head == m.tq.nqrv[m.name].tail {
		if len(nonblock) != 0 && nonblock[0] {
			return ""
		} else {
			for m.tq.nqrv[m.name].head == m.tq.nqrv[m.name].tail {
				Time.Sleep(0.1)
			}
		}
	}

	value, err := m.db.Get([]byte(m.name+"_"+Str(m.tq.nqrv[m.name].head)), nil)
	Panicerr(err)

	err = m.db.Delete([]byte(m.name+"_"+Str(m.tq.nqrv[m.name].head)), nil)
	Panicerr(err)

	m.tq.nqrv[m.name].head += 1

	err = m.db.Put([]byte(m.name+"_head"), []byte(Str(m.tq.nqrv[m.name].head)), nil)
	Panicerr(err)

	return Str(value)
}

func (m *NamedQueueStruct) Put(value string) {
	if value == "" {
		Panicerr("value can not be empty")
	}

	m.tq.nqrv[m.name].lock.Acquire()
	defer m.tq.nqrv[m.name].lock.Release()

	err := m.db.Put([]byte(m.name+"_"+Str(m.tq.nqrv[m.name].tail)), []byte(value), nil)
	Panicerr(err)

	m.tq.nqrv[m.name].tail += 1

	err = m.db.Put([]byte(m.name+"_tail"), []byte(Str(m.tq.nqrv[m.name].tail)), nil)
	Panicerr(err)
}
