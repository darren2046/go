package golanglibs

type QueueStruct struct {
	db     *DatabaseStruct
	dbpath string
	closed bool
}

func getQueue(dbpath string) (q *QueueStruct) {
	q = &QueueStruct{}

	q.db = Tools.SQLite(dbpath)
	q.dbpath = dbpath

	return
}

func (m *QueueStruct) Close() {
	m.db.Close()
	unlink(m.dbpath)

	m.closed = true
}

type NamedQueueStruct struct {
	db   *DatabaseStruct
	name string
	tq   *QueueStruct // Top Queue
}

// Will not clean the data already exists
func (m *QueueStruct) New(queueName ...string) *NamedQueueStruct {
	q := &NamedQueueStruct{}

	n := "__empty__name__queue__"
	if len(queueName) != 0 {
		n = queueName[0]
	}

	m.db.CreateTable(n).
		AddColumn("data", "text")

	q.db = m.db
	q.name = n
	q.tq = m

	return q
}

func (m *NamedQueueStruct) Size() int64 {
	return m.db.Table(m.name).Count()
}

// 线程不安全
func (m *NamedQueueStruct) Get(nonblock ...bool) string {
	r := m.db.Table(m.name).First()
	if len(r) == 0 {
		if len(nonblock) != 0 && nonblock[0] {
			return ""
		} else {
			r = m.db.Table(m.name).First()
			for len(r) == 0 {
				Time.Sleep(0.1)
			}
		}
	}

	m.db.Table(m.name).Where("id", "=", Int(r["id"])).Delete()

	return Base64.Decode(Str(r["data"]))
}

func (m *NamedQueueStruct) Put(value string) {
	if value == "" {
		Panicerr("value can not be empty")
	}

	m.db.Table(m.name).
		Data(map[string]interface{}{
			"data": Base64.Encode(value),
		}).Insert()
}
