package golanglibs

type toolsStruct struct {
	Lock          func() *LockStruct
	AliDNS        func(accessKeyID string, accessKeySecret string) *alidnsStruct
	Chart         *chartStruct
	CloudflareDNS func(key string, email string) *cloudflareStruct
	Compress      *compressStruct
	Crontab       func() *crontabStruct
	GodaddyDNS    func(key string, secret string) *godaddyStruct
	Ini           func(fpath ...string) *iniStruct
	JavascriptVM  func() *javascriptVMStruct
	Matrix        func(homeserverURL string) *MatrixStruct
	Nats          func(server string) *natsStruct
	Totp          func(key string) *totpStruct
	Pexpect       func(command string) *pexpectStruct
	ProgressBar   func(title string, total int64, showBytes ...bool) *progressBarStruct
	Prometheus    func(url string) *prometheusStruct
	MySQL         func(host string, port int, user string, password string, db string, cfg ...DatabaseConfig) *databaseStruct
	SQLite        func(dbpath string) *databaseStruct
	RabbitMQ      func(rabbitMQURL string, queueName string) *rabbitConnectionStruct
	RateLimit     func(rate int) *rateLimitStruct
	Redis         func(host string, port int, password string, db int, cfg ...redisConfig) *RedisStruct
	Selenium      func(url string, browser ...string) *seleniumStruct
	SSH           func(user string, pass string, host string, port int) *sshStruct
	StatikOpen    func(path string) *statikFileStruct
	Table         func(header ...string) *tableStruct
	TelegramBot   func(token string) *telegramBotStruct
	Telegraph     func(AuthorName string) *telegraphStruct
	URL           func(url string) *urlStruct
	TTLCache      func(ttlsecond interface{}) *ttlCacheStruct
	VNC           func(server string, cfg ...VNCCfg) *vncStruct
	WebSocket     func(url string) *websocketStruct
	Xlsx          func(path string) *xlsxStruct
	XPath         func(htmlString string) *xpathStruct
	Sysinfo       *sysinfoStruct
	Queue         func(datadir string) (q *queueStruct)
}

var Tools toolsStruct

func init() {
	Tools = toolsStruct{
		Lock:          getLock,
		AliDNS:        getAlidns,
		Chart:         &chartstruct,
		CloudflareDNS: getCloudflare,
		Compress:      &compressstruct,
		Crontab:       getCrontab,
		GodaddyDNS:    getGodaddy,
		Ini:           getIni,
		JavascriptVM:  getJavascriptVM,
		Matrix:        getMatrix,
		Nats:          getNats,
		Totp:          getTotp,
		Pexpect:       pexpect,
		ProgressBar:   getProgressBar,
		Prometheus:    getPrometheus,
		MySQL:         getMySQL,
		SQLite:        getSQLite,
		RabbitMQ:      getRabbitMQ,
		RateLimit:     getRateLimit,
		Redis:         getRedis,
		Selenium:      getSelenium,
		SSH:           getSSH,
		StatikOpen:    statikOpen,
		Table:         getTable,
		TelegramBot:   getTelegramBot,
		Telegraph:     getTelegraph,
		URL:           getUrl,
		TTLCache:      getTTLCache,
		VNC:           getVNC,
		WebSocket:     getWebSocket,
		Xlsx:          getXlsx,
		XPath:         getXPath,
		Sysinfo:       &sysinfostruct,
		Queue:         getQueue,
	}
}
