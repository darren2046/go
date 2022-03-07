package golanglibs

import (
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseStruct struct {
	engin                  *gorose.Engin
	driver                 string
	dsn                    string
	isclose                bool
	networkErrorRetryTimes int // 网络错误重试次数
}

type DatabaseConfig struct {
	Timeout                int
	Charset                string
	NetworkErrorRetryTimes int // 网络错误重试次数
}

// 用来过滤报错的信息， 如果包含有如下的某一个， 就判断为是网络错误
var databaseNetworkErrorStrings = []string{
	"timeout",
	"invalid connection",
	"no such host",
	"connection refused",
	"bad connection",
}

func doDatabaseThingsAndHandleNetworkError(retry int, f func() error) {
	errortimes := 0
	var err error
	for {
		err = f()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if String(err).In(errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < retry {
				errortimes += 1
				sleep(3)
			} else {
				Panicerr(err)
			}
		} else {
			break
		}
	}
}

func getMySQL(host string, port int, user string, password string, db string, cfg ...DatabaseConfig) *DatabaseStruct {
	var timeoutt int
	var chartsett string
	var networkErrorRetryTimess int

	if len(cfg) != 0 {
		if cfg[0].Timeout != 0 {
			timeoutt = cfg[0].Timeout
		} else {
			timeoutt = 10
		}
		if cfg[0].Charset != "" {
			chartsett = cfg[0].Charset
		} else {
			chartsett = "utf8mb4"
		}
		if cfg[0].NetworkErrorRetryTimes != 0 {
			networkErrorRetryTimess = cfg[0].NetworkErrorRetryTimes
		} else {
			networkErrorRetryTimess = 10
		}
	} else {
		timeoutt = 10
		chartsett = "utf8mb4"
		networkErrorRetryTimess = 0
	}

	m := &DatabaseStruct{}
	m.networkErrorRetryTimes = networkErrorRetryTimess
	m.driver = "mysql"
	m.dsn = user + ":" + password + "@tcp(" + host + ":" + Str(port) + ")/" + db + "?timeout=" + Str(timeoutt) + "s&readTimeout=" + Str(timeoutt) + "s&writeTimeout=" + Str(timeoutt) + "s&charset=" + chartsett
	var config = &gorose.Config{
		Driver: m.driver,
		Dsn:    m.dsn,
	}

	var engin *gorose.Engin
	var err error
	doDatabaseThingsAndHandleNetworkError(m.networkErrorRetryTimes, func() error {
		engin, err = gorose.Open(config)
		return err
	})

	m.engin = engin

	// 它会重连，如果连接坏了
	go func(m *DatabaseStruct) {
		for {
			sleep(3)
			Try(func() {
				m.engin.Ping()
			})
			if m.isclose {
				break
			}
		}
	}(m)

	return m
}

func getSQLite(dbpath string) *DatabaseStruct {
	m := &DatabaseStruct{}
	m.driver = "sqlite3"
	m.dsn = dbpath
	var config = &gorose.Config{
		Driver: m.driver,
		Dsn:    m.dsn,
	}

	engin, err := gorose.Open(config)
	Panicerr(err)

	m.engin = engin

	return m
}

func (m *DatabaseStruct) Query(sql string, args ...interface{}) []gorose.Data {
	db := m.engin.NewOrm()
	res, err := db.Query(sql, args...)
	Panicerr(err)

	for idx, d := range res {
		for k, v := range d {
			if v != nil && Typeof(v) == "time.Time" {
				res[idx][k] = Time.Strftime("%Y-%m-%d %H:%M:%S", v.(time.Time).Unix())
			}
		}
	}
	return res
}

func (m *DatabaseStruct) Close() {
	m.isclose = true
	m.engin.GetQueryDB().Close()
}

func (m *DatabaseStruct) Execute(sql string) int64 {
	db := m.engin.NewOrm()
	res, err := db.Execute(sql)
	Panicerr(err)
	return res
}

type databaseOrmStruct struct {
	orm    gorose.IOrm
	db     *DatabaseStruct
	driver string
	table  string
}

func (m *DatabaseStruct) Table(tbname string) *databaseOrmStruct {
	orm := m.engin.NewOrm()
	return &databaseOrmStruct{
		orm:    orm.Table("`" + tbname + "`"),
		driver: m.driver,
		table:  tbname,
		db:     m,
	}
}

func (m *DatabaseStruct) RenameTable(oldTableName string, newTableNname string) {
	if m.driver == "mysql" {
		m.Execute("RENAME TABLE `" + oldTableName + "` TO `" + newTableNname + "`;")
	} else if m.driver == "sqlite3" {
		m.Execute("ALTER TABLE `" + oldTableName + "` RENAME TO `" + newTableNname + "`;")
	}
}

func (m *databaseOrmStruct) Fields(items ...string) *databaseOrmStruct {
	var i []string
	for _, v := range items {
		i = append(i, "`"+v+"`")
	}
	m.orm.Fields(i...)
	return m
}

func (m *databaseOrmStruct) Where(key string, operator string, value interface{}) *databaseOrmStruct {
	m.orm.Where(key, operator, value)
	return m
}

func (m *databaseOrmStruct) WhereIn(key string, value []interface{}) *databaseOrmStruct {
	m.orm.WhereIn(key, value)
	return m
}

func (m *databaseOrmStruct) WhereNotIn(key string, value []interface{}) *databaseOrmStruct {
	m.orm.WhereNotIn(key, value)
	return m
}

func (m *databaseOrmStruct) WhereNull(columnName string) *databaseOrmStruct {
	m.orm.WhereNull(columnName)
	return m
}

func (m *databaseOrmStruct) WhereNotNull(columnName string) *databaseOrmStruct {
	m.orm.WhereNotNull(columnName)
	return m
}

func (m *databaseOrmStruct) OrWhere(key string, operator string, value interface{}) *databaseOrmStruct {
	m.orm.OrWhere(key, operator, value)
	return m
}

func (m *databaseOrmStruct) OrWhereIn(key string, value []interface{}) *databaseOrmStruct {
	m.orm.OrWhereIn(key, value)
	return m
}

func (m *databaseOrmStruct) Orderby(item ...string) *databaseOrmStruct {
	m.orm.OrderBy(String(" ").Join(item).Get())
	return m
}

func (m *databaseOrmStruct) Limit(number int) *databaseOrmStruct {
	m.orm.Limit(number)
	return m
}

func (m *databaseOrmStruct) Get() (res []gorose.Data) {
	var err error
	doDatabaseThingsAndHandleNetworkError(m.db.networkErrorRetryTimes, func() error {
		res, err = m.orm.Get()
		return err
	})

	m.orm = m.db.Table(m.table).orm

	//print(m)

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
	return m.Limit(pagesize).Offset(offset).Get()
}

func (m *databaseOrmStruct) First() (res gorose.Data) {
	var err error
	doDatabaseThingsAndHandleNetworkError(m.db.networkErrorRetryTimes, func() error {
		res, err = m.orm.First()
		return err
	})

	m.orm = m.db.Table(m.table).orm

	for k, v := range res {
		if v != nil && Typeof(v) == "time.Time" {
			res[k] = Time.Strftime("%Y-%m-%d %H:%M:%S", v.(time.Time).Unix())
		}
	}

	return
}

func (m *databaseOrmStruct) Find(id int) gorose.Data {
	return m.Where("id", "=", id).First()
}

func (m *databaseOrmStruct) Count() (res int64) {
	var err error
	doDatabaseThingsAndHandleNetworkError(m.db.networkErrorRetryTimes, func() error {
		res, err = m.orm.Count()
		return err
	})

	m.orm = m.db.Table(m.table).orm

	return
}

func (m *databaseOrmStruct) Exists() (res bool) {
	if m.Count() == 0 {
		res = false
	} else {
		res = true
	}

	return
}

func (m *databaseOrmStruct) Chunk(limit int, callback func([]gorose.Data) error) {
	err := m.orm.Chunk(limit, callback)
	Panicerr(err)
}

func (m *databaseOrmStruct) BuildSQL() (string, []interface{}) {
	sql, param, err := m.orm.BuildSql()
	Panicerr(err)

	return sql, param
}

func (m *databaseOrmStruct) Data(data interface{}) *databaseOrmStruct {
	m.orm.Data(data)
	return m
}

func (m *databaseOrmStruct) Offset(offset int) *databaseOrmStruct {
	m.orm.Offset(offset)
	return m
}

func (m *databaseOrmStruct) InsertGetID() (num int64) {
	var err error
	doDatabaseThingsAndHandleNetworkError(m.db.networkErrorRetryTimes, func() error {
		num, err = m.orm.InsertGetId()
		return err
	})

	m.orm = m.db.Table(m.table).orm

	return
}

func (m *databaseOrmStruct) Insert() (num int64) {
	var err error
	doDatabaseThingsAndHandleNetworkError(m.db.networkErrorRetryTimes, func() error {
		num, err = m.orm.Insert()
		return err
	})

	m.orm = m.db.Table(m.table).orm

	return
}

func (m *databaseOrmStruct) Update(data ...interface{}) (num int64) {
	var err error
	doDatabaseThingsAndHandleNetworkError(m.db.networkErrorRetryTimes, func() error {
		num, err = m.orm.Update(data...)
		return err
	})

	m.orm = m.db.Table(m.table).orm

	return
}

func (m *databaseOrmStruct) Delete() (num int64) {
	var err error
	doDatabaseThingsAndHandleNetworkError(m.db.networkErrorRetryTimes, func() error {
		num, err = m.orm.Delete()
		return err
	})

	m.orm = m.db.Table(m.table).orm

	return
}

func (m *DatabaseStruct) Tables() (res []string) {
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

func (m *DatabaseStruct) CreateTable(tableName string, engineName ...string) *databaseOrmStruct {
	if !Array(m.Tables()).Has(tableName) {
		if len(engineName) != 0 && m.driver != "mysql" {
			Panicerr("SQLite不支持设定存储引擎")
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
			Panicerr(err)
		}

		columnTypeMaps := map[string]map[string]string{
			"mysql": {
				"int":      "bigint",
				"float":    "double",
				"string":   "VARCHAR(256)",
				"text":     "LONGTEXT",
				"datetime": "DATETIME",
			},
			"sqlite3": {
				"int":      "INTEGER",
				"float":    "REAL",
				"string":   "VARCHAR",
				"text":     "LONGTEXT",
				"datetime": "DATETIME",
				"blob":     "BLOB",
			},
		}

		if m.driver == "mysql" && columnType == "blob" {
			Panicerr("MySQL暂不支持blob")
		}
		columnType = columnTypeMaps[m.driver][columnType]

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
