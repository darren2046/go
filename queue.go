package golanglibs

type QueueStruct struct {
	db *DatabaseStruct
}

func getQueue(db *DatabaseStruct) (q *QueueStruct) {
	q = &QueueStruct{}

	q.db = db

	return
}

type NamedQueueStruct struct {
	db   *DatabaseStruct
	name string
	tq   *QueueStruct // Top Queue
}

// Will not clean the data already exists
func (m *QueueStruct) New(queueName ...string) *NamedQueueStruct {
	q := &NamedQueueStruct{}

	n := "__queue__empty__name__"
	if len(queueName) != 0 {
		n = "__queue__name__" + queueName[0]
	}

	if !Array(m.db.Tables()).Has(n) {
		m.db.CreateTable(n).
			AddColumn("data", "text")
	}

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
