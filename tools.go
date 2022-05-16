package golanglibs

type toolsStruct struct {
	Lock                   func() *LockStruct
	RWLock                 func() *RWLockStruct
	AliDNS                 func(accessKeyID string, accessKeySecret string) *alidnsStruct
	Chart                  *chartStruct
	CloudflareDNS          func(key string, email string) *cloudflareStruct
	Compress               *compressStruct
	Crontab                func() *crontabStruct
	GodaddyDNS             func(key string, secret string) *godaddyStruct
	Ini                    func(fpath ...string) *iniStruct
	JavascriptVM           func() *javascriptVMStruct
	Matrix                 func(homeserverURL string) *MatrixStruct
	Nats                   func(server string) *natsStruct
	Totp                   func(key string) *totpStruct
	Pexpect                func(command string) *PexpectStruct
	ProgressBar            func(title string, total int64, showBytes ...bool) *progressBarStruct
	PrometheusClient       func(url string) *prometheusClientStruct
	PrometheusMetricServer func(listenAddr string, path ...string) *prometheusMetricServerStruct
	MySQL                  func(host string, port int, user string, password string, db string, cfg ...DatabaseConfig) *DatabaseStruct
	SQLite                 func(dbpath string) *DatabaseStruct
	RabbitMQ               func(rabbitMQURL string, queueName string) *rabbitConnectionStruct
	RateLimit              func(rate int) *rateLimitStruct
	Redis                  func(host string, port int, cfg ...RedisConfig) *RedisStruct
	SeleniumLocal          func() *SeleniumStruct
	SeleniumRemote         func(serverURL string) *SeleniumStruct
	SSH                    func(user string, pass string, host string, port int) *sshStruct
	StatikOpen             func(path string) *statikFileStruct
	Table                  func(header ...string) *tableStruct
	TelegramBot            func(token string) *TelegramBotStruct
	Telegraph              func(AuthorName string) *telegraphStruct
	URL                    func(url string) *urlStruct
	TTLCache               func(ttlsecond interface{}) *ttlCacheStruct
	VNC                    func(server string, cfg ...VNCCfg) *VNCStruct
	WebSocket              func(url string) *websocketStruct
	Xlsx                   func(path string) *xlsxStruct
	XPath                  func(htmlString string) *xpathStruct
	Sysinfo                *sysinfoStruct
	Queue                  func(datadir string) (q *QueueStruct)
	Jieba                  func() *JiebaStruct
	Telegram               func(AppID int32, AppHash string, config ...TelegramConfig) *TelegramStruct
	Elasticsearch          func(baseurl string) *ElasticsearchStruct
	LevelDB                func(datadir string) (l *LevelDBStruct)
}

var Tools toolsStruct

func init() {
	Tools = toolsStruct{
		RWLock:                 getRWLock,
		Lock:                   getLock,
		AliDNS:                 getAlidns,
		Chart:                  &chartstruct,
		CloudflareDNS:          getCloudflare,
		Compress:               &compressstruct,
		Crontab:                getCrontab,
		GodaddyDNS:             getGodaddy,
		Ini:                    getIni,
		JavascriptVM:           getJavascriptVM,
		Matrix:                 getMatrix,
		Nats:                   getNats,
		Totp:                   getTotp,
		Pexpect:                pexpect,
		ProgressBar:            getProgressBar,
		PrometheusClient:       getPrometheusClient,
		PrometheusMetricServer: getPrometheusMetricServer,
		MySQL:                  getMySQL,
		SQLite:                 getSQLite,
		RabbitMQ:               getRabbitMQ,
		RateLimit:              getRateLimit,
		Redis:                  getRedis,
		SeleniumLocal:          getSeleniumLocal,
		SeleniumRemote:         getSeleniumRemote,
		SSH:                    getSSH,
		StatikOpen:             statikOpen,
		Table:                  getTable,
		TelegramBot:            getTelegramBot,
		Telegraph:              getTelegraph,
		URL:                    getUrl,
		TTLCache:               getTTLCache,
		VNC:                    getVNC,
		WebSocket:              getWebSocket,
		Xlsx:                   getXlsx,
		XPath:                  getXPath,
		Sysinfo:                &sysinfostruct,
		Queue:                  getQueue,
		Telegram:               getTelegram,
		Jieba:                  getJieba,
		Elasticsearch:          getElasticsearch,
		LevelDB:                getLevelDB,
	}
}
