package golanglibs

import (
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
	_ "github.com/mattn/go-sqlite3"
)

type databaseStruct struct {
	engin                  *gorose.Engin
	driver                 string
	dsn                    string
	isclose                bool
	networkErrorRetryTimes int // 网络错误重试次数
}

type DatabaseConfig struct {
	timeout                int
	charset                string
	networkErrorRetryTimes int // 网络错误重试次数
}

// 用来过滤报错的信息， 如果包含有如下的某一个， 就判断为是网络错误
var databaseNetworkErrorStrings = []string{
	"timeout",
	"invalid connection",
	"no such host",
	"connection refused",
	"bad connection",
}

func getMySQL(host string, port int, user string, password string, db string, cfg ...DatabaseConfig) *databaseStruct {
	var timeoutt int
	var chartsett string
	var networkErrorRetryTimess int

	if len(cfg) != 0 {
		if cfg[0].timeout != 0 {
			timeoutt = cfg[0].timeout
		} else {
			timeoutt = 10
		}
		if cfg[0].charset != "" {
			chartsett = cfg[0].charset
		} else {
			chartsett = "utf8mb4"
		}
		if cfg[0].networkErrorRetryTimes != 0 {
			networkErrorRetryTimess = cfg[0].networkErrorRetryTimes
		} else {
			networkErrorRetryTimess = 10
		}
	} else {
		timeoutt = 10
		chartsett = "utf8mb4"
		networkErrorRetryTimess = 0
	}

	m := &databaseStruct{}
	m.networkErrorRetryTimes = networkErrorRetryTimess
	m.driver = "mysql"
	m.dsn = user + ":" + password + "@tcp(" + host + ":" + Str(port) + ")/" + db + "?timeout=" + Str(timeoutt) + "s&readTimeout=" + Str(timeoutt) + "s&writeTimeout=" + Str(timeoutt) + "s&charset=" + chartsett
	var config = &gorose.Config{
		Driver: m.driver,
		Dsn:    m.dsn,
	}

	errortimes := 0
	var err error
	var engin *gorose.Engin
	for {
		engin, err = gorose.Open(config)
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if String(err).In(errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}

	m.engin = engin

	// 它会重连，如果连接坏了
	go func(m *databaseStruct) {
		for {
			sleep(3)
			try(func() {
				m.engin.Ping()
			})
			if m.isclose {
				break
			}
		}
	}(m)

	return m
}

func getSQLite(dbpath string) *databaseStruct {
	m := &databaseStruct{}
	m.driver = "sqlite3"
	m.dsn = dbpath
	var config = &gorose.Config{
		Driver: m.driver,
		Dsn:    m.dsn,
	}

	engin, err := gorose.Open(config)
	panicerr(err)

	m.engin = engin

	return m
}

func (m *databaseStruct) Query(sql string, args ...interface{}) []gorose.Data {
	db := m.engin.NewOrm()
	res, err := db.Query(sql, args...)
	panicerr(err)

	for idx, d := range res {
		for k, v := range d {
			if v != nil && Typeof(v) == "time.Time" {
				res[idx][k] = Time.Strftime("%Y-%m-%d %H:%M:%S", v.(time.Time).Unix())
			}
		}
	}
	return res
}

func (m *databaseStruct) Close() {
	m.isclose = true
	m.engin.GetQueryDB().Close()
}

func (m *databaseStruct) Execute(sql string) int64 {
	db := m.engin.NewOrm()
	res, err := db.Execute(sql)
	panicerr(err)
	return res
}

type databaseOrmStruct struct {
	orm    gorose.IOrm
	db     *databaseStruct
	driver string
	table  string
	lock   *lockStruct
	lockby int64
}

func (m *databaseStruct) Table(tbname string) *databaseOrmStruct {
	orm := m.engin.NewOrm()
	return &databaseOrmStruct{
		orm:    orm.Table("`" + tbname + "`"),
		driver: m.driver,
		table:  tbname,
		db:     m,
		lock:   getLock(), // 为了保证线程安全，链式操作当中要上锁，返回数据解锁。为了保证在多个线程中复用同一个databaseStruct的时候报错。
		lockby: -1,
	}
}

func (m *databaseStruct) RenameTable(oldTableName string, newTableNname string) {
	if m.driver == "mysql" {
		m.Execute("RENAME TABLE `" + oldTableName + "` TO `" + newTableNname + "`;")
	} else if m.driver == "sqlite3" {
		m.Execute("ALTER TABLE `" + oldTableName + "` RENAME TO `" + newTableNname + "`;")
	}
}

func (m *databaseOrmStruct) Fields(items ...string) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	var i []string
	for _, v := range items {
		i = append(i, "`"+v+"`")
	}
	m.orm.Fields(i...)
	return m
}

func (m *databaseOrmStruct) Where(key string, operator string, value interface{}) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	m.orm.Where(key, operator, value)
	return m
}

func (m *databaseOrmStruct) WhereIn(key string, value []interface{}) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	m.orm.WhereIn(key, value)
	return m
}

func (m *databaseOrmStruct) WhereNotIn(key string, value []interface{}) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	m.orm.WhereNotIn(key, value)
	return m
}

func (m *databaseOrmStruct) WhereNull(columnName string) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	m.orm.WhereNull(columnName)
	return m
}

func (m *databaseOrmStruct) WhereNotNull(columnName string) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	m.orm.WhereNotNull(columnName)
	return m
}

func (m *databaseOrmStruct) OrWhere(key string, operator string, value interface{}) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	m.orm.OrWhere(key, operator, value)
	return m
}

func (m *databaseOrmStruct) OrWhereIn(key string, value []interface{}) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	m.orm.OrWhereIn(key, value)
	return m
}

func (m *databaseOrmStruct) Orderby(item ...string) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	m.orm.OrderBy(String(" ").Join(item).Get())
	return m
}

func (m *databaseOrmStruct) Limit(number int) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	m.orm.Limit(number)
	return m
}

func (m *databaseOrmStruct) Get() (res []gorose.Data) {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}

	errortimes := 0
	var err error
	for {
		res, err = m.orm.Get()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if String(err).In(errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}

	m.orm = m.db.Table(m.table).orm

	//print(m)

	m.lock.Release()
	m.lockby = -1

	for idx, d := range res {
		for k, v := range d {
			if v != nil && Typeof(v) == "time.Time" {
				res[idx][k] = Time.Strftime("%Y-%m-%d %H:%M:%S", v.(time.Time).Unix())
			}
		}
	}

	return
}

func (m *databaseOrmStruct) Paginate(pagesize int, page int) []gorose.Data {
	offset := pagesize * (page - 1)
	return m.Limit(pagesize).offset(offset).Get()
}

func (m *databaseOrmStruct) First() (res gorose.Data) {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}

	errortimes := 0
	var err error
	for {
		res, err = m.orm.First()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if String(err).In(errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}

	m.orm = m.db.Table(m.table).orm

	m.lock.Release()
	m.lockby = -1

	for k, v := range res {
		if v != nil && Typeof(v) == "time.Time" {
			res[k] = Time.Strftime("%Y-%m-%d %H:%M:%S", v.(time.Time).Unix())
		}
	}

	return
}

func (m *databaseOrmStruct) find(id int) gorose.Data {
	return m.Where("id", "=", id).First()
}

func (m *databaseOrmStruct) Count() (res int64) {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}

	errortimes := 0
	var err error
	for {
		res, err = m.orm.Count()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if String(err).In(errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}

	m.orm = m.db.Table(m.table).orm

	m.lock.Release()
	m.lockby = -1

	return
}

func (m *databaseOrmStruct) exists() (res bool) {
	if m.Count() == 0 {
		res = false
	} else {
		res = true
	}

	return
}

func (m *databaseOrmStruct) chunk(limit int, callback func([]gorose.Data) error) {
	err := m.orm.Chunk(limit, callback)
	panicerr(err)
}

func (m *databaseOrmStruct) buildSQL() (string, []interface{}) {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	sql, param, err := m.orm.BuildSql()
	panicerr(err)

	m.lock.Release()
	m.lockby = -1

	return sql, param
}

func (m *databaseOrmStruct) data(data interface{}) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	m.orm.Data(data)
	return m
}

func (m *databaseOrmStruct) offset(offset int) *databaseOrmStruct {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	m.orm.Offset(offset)
	return m
}

func (m *databaseOrmStruct) insertGetID() (num int64) {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	errortimes := 0
	var err error
	for {
		num, err = m.orm.InsertGetId()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if String(err).In(errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}

	m.orm = m.db.Table(m.table).orm

	m.lock.Release()
	m.lockby = -1

	return
}

func (m *databaseOrmStruct) Insert() (num int64) {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}

	errortimes := 0
	var err error
	for {
		num, err = m.orm.Insert()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if String(err).In(errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}
	m.orm = m.db.Table(m.table).orm

	m.lock.Release()
	m.lockby = -1

	return
}

func (m *databaseOrmStruct) Update(data ...interface{}) (num int64) {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}
	errortimes := 0
	var err error
	for {
		num, err = m.orm.Update(data...)
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if String(err).In(errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}
	m.orm = m.db.Table(m.table).orm

	m.lock.Release()
	m.lockby = -1

	return
}

func (m *databaseOrmStruct) Delete() (num int64) {
	if m.lockby != Os.GoroutineID() {
		m.lock.Acquire()
		m.lockby = Os.GoroutineID()
	}

	errortimes := 0
	var err error
	for {
		num, err = m.orm.Delete()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if String(err).In(errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}
	m.orm = m.db.Table(m.table).orm

	m.lock.Release()
	m.lockby = -1

	return
}

func (m *databaseStruct) tables() (res []string) {
	if m.driver == "mysql" {
		for _, v := range m.Query("show tables;") {
			for _, i := range v {
				res = append(res, Str(i))
			}
		}
	} else if m.driver == "sqlite3" {
		for _, i := range m.Query("SELECT `name` FROM sqlite_master WHERE type='table';") {
			res = append(res, Str(i["name"]))
		}
	}
	return
}

func (m *databaseStruct) createTable(tableName string, engineName ...string) *databaseOrmStruct {
	if !Array(m.tables()).Has(tableName) {
		if len(engineName) != 0 && m.driver != "mysql" {
			panicerr("SQLite不支持设定存储引擎")
		}
		if m.driver == "mysql" {
			if len(engineName) != 0 {
				m.Execute("CREATE TABLE `" + tableName + "`(`id` BIGINT UNSIGNED AUTO_INCREMENT,PRIMARY KEY ( `id` ))ENGINE=" + engineName[0] + " DEFAULT CHARSET=utf8mb4;")
			} else {
				m.Execute("CREATE TABLE `" + tableName + "`(`id` BIGINT UNSIGNED AUTO_INCREMENT,PRIMARY KEY ( `id` ))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
			}
		} else {
			m.Execute("CREATE TABLE `" + tableName + "` (`id` INTEGER PRIMARY KEY AUTOINCREMENT)")
		}
	}
	return m.Table(tableName)
}

func (m *databaseOrmStruct) DropTable() int64 {
	return m.db.Execute("DROP TABLE `" + m.table + "`")
}

func (m *databaseOrmStruct) TruncateTable() (status int64) {
	if m.driver == "mysql" {
		status = m.db.Execute("TRUNCATE TABLE `" + m.table + "`")
	} else {
		status = m.db.Execute("DELETE FROM TABLE `" + m.table + "`")
	}
	return
}

func (m *databaseOrmStruct) AddColumn(columnName string, columnType string, defaultValue ...string) *databaseOrmStruct {
	// lg.debug(columnName, m.columns())
	if !Map(m.Columns()).Has(columnName) {
		if !Array([]string{"int", "float", "string", "text", "datetime", "blob"}).Has(columnType) {
			err := errors.New("只支持以下数据类型：\"int\", \"float\", \"string\", \"text\", \"datetime\", \"blob\"(SQLite支持, MySql不支持)")
			panicerr(err)
		}
		if m.driver == "mysql" {
			if columnType == "int" {
				columnType = "bigint"
			} else if columnType == "float" {
				columnType = "double"
			} else if columnType == "string" {
				columnType = "VARCHAR(256)"
			} else if columnType == "text" {
				columnType = "LONGTEXT"
			} else if columnType == "datetime" {
				columnType = "DATETIME"
			} else if columnType == "blob" {
				columnType = "LONGBLOB"
				panicerr("MySQL暂不支持blob")
			}
		} else {
			if columnType == "int" {
				columnType = "INTEGER"
			} else if columnType == "float" {
				columnType = "REAL"
			} else if columnType == "string" {
				columnType = "VARCHAR"
			} else if columnType == "text" {
				columnType = "LONGTEXT"
			} else if columnType == "datetime" {
				columnType = "DATETIME"
			} else if columnType == "blob" {
				columnType = "BLOB"
			}
		}

		var sql string
		if len(defaultValue) == 0 {
			sql = "ALTER TABLE `" + m.table + "` ADD `" + columnName + "` " + columnType + ";"
		} else {
			sql = "ALTER TABLE `" + m.table + "` ADD `" + columnName + "` " + columnType + " DEFAULT \"" + defaultValue[0] + "\";"
		}

		m.db.Execute(sql)

	}
	return m
}

func (m *databaseOrmStruct) DropColumn(columnName string) *databaseOrmStruct {
	if Map(m.Columns()).Has(columnName) {
		if m.driver == "mysql" {
			m.db.Execute("ALTER TABLE `" + m.table + "`  DROP " + columnName)
		} else {
			panic(newerr("SQLite does not support drop column"))
		}
	}
	return m
}

func (m *databaseOrmStruct) AddIndex(columnName ...string) *databaseOrmStruct {
	if !m.IndexExists(columnName...) {
		columns := "`" + String("`,`").Join(columnName).Get() + "`"
		indexName := "idx_" + String("_").Join(columnName).Get()
		if m.driver == "mysql" {
			m.db.Execute("ALTER TABLE `" + m.table + "` ADD INDEX " + indexName + "(" + columns + ")")
		} else {
			m.db.Execute("CREATE INDEX " + indexName + " on `" + m.table + "` (" + columns + ");")
		}
	}
	return m
}

func (m *databaseOrmStruct) IndexExists(columnName ...string) (exists bool) {
	indexName := "idx_" + String("_").Join(columnName).Get()
	if m.driver == "mysql" {
		for _, v := range m.db.Query("SHOW INDEX FROM `" + m.table + "`") {
			if v["Key_name"] == indexName {
				exists = true
			}
		}
	} else if m.driver == "sqlite3" {
		if Int(m.db.Query("SELECT count(name) as `count` FROM sqlite_master WHERE type='index' AND name='" + indexName + "';")[0]["count"]) == 1 {
			exists = true
		}
	}
	return
}

func (m *databaseOrmStruct) DropIndex(columnName ...string) *databaseOrmStruct {
	indexName := "idx_" + String("_").Join(columnName).Get()
	if m.driver == "mysql" {
		m.db.Execute("ALTER TABLE `" + m.table + "` DROP INDEX " + indexName)
	} else {
		m.db.Execute("DROP INDEX " + indexName)
	}
	return m
}

func (m *databaseOrmStruct) Columns() (res map[string]string) {
	res = make(map[string]string)

	if m.driver == "mysql" {
		for _, i := range m.db.Query("SHOW COLUMNS FROM `" + m.table + "`;") {
			// lg.debug(i)
			if String(Str(i["Type"])).Lower().Get() == "bigint(20)" {
				res[Str(i["Field"])] = "int"
			} else if String(Str(i["Type"])).Lower().Get() == "double" {
				res[Str(i["Field"])] = "float"
			} else if String(Str(i["Type"])).Lower().Get() == "varchar(512)" {
				res[Str(i["Field"])] = "string"
			} else if String(Str(i["Type"])).Lower().Get() == "longtext" {
				res[Str(i["Field"])] = "text"
			} else if String(Str(i["Type"])).Lower().Get() == "datetime" {
				res[Str(i["Field"])] = "datetime"
			} else if String(Str(i["Type"])).Lower().Get() == "longblob" {
				res[Str(i["Field"])] = "blob"
			} else {
				res[Str(i["Field"])] = Str(i["Type"])
			}
		}
	} else {
		for _, i := range m.db.Query("PRAGMA table_info(`" + m.table + "`);") {
			if String(Str(i["type"])).Upper().Get() == "INTEGER" {
				res[Str(i["name"])] = "int"
			} else if String(Str(i["type"])).Upper().Get() == "REAL" {
				res[Str(i["name"])] = "float"
			} else if String(Str(i["type"])).Upper().Get() == "VARCHAR" {
				res[Str(i["name"])] = "string"
			} else if String(Str(i["type"])).Upper().Get() == "LONGTEXT" {
				res[Str(i["name"])] = "text"
			} else if String(Str(i["type"])).Upper().Get() == "DATETIME" {
				res[Str(i["name"])] = "datetime"
			} else if String(Str(i["type"])).Upper().Get() == "BLOB" {
				res[Str(i["name"])] = "blob"
			} else {
				res[Str(i["name"])] = Str(i["type"])
			}
		}
	}
	return
}
