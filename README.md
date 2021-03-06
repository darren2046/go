# golanglibs

This is a toolkit that provide a lot of function or object that make programing easier like Python. 

# Index

* Tools
    * CSV
        * func Reader(fpath string) *csvReader 
            * func (m *csvReader) Read() (res map[string]string) 
            * func (m *csvReader) Readrows() chan map[string]string 
            * func (m *csvReader) Close() 
        * func Writer(fpath string, mode string) *csvWriter 
            * func (m *csvWriter) Flush() 
            * func (m *csvWriter) SetHeaders(headers []string) 
            * func (m *csvWriter) Write(record map[string]string) 
            * func (m *csvWriter) Close() 
    * func LevelDB(datadir string) (l \*LevelDB) 
        * func (m \*LevelDB) Close() 
        * func (m \*LevelDB) Destroy() 
        * func (m \*LevelDB) Exists(key string) bool 
        * func (m \*LevelDB) Get(key string) string 
        * func (m \*LevelDB) Set(key string, value string) 
        * func (m \*LevelDB) Delete(key string) 
    * func Elasticsearch(baseurl string) \*Elasticsearch
        * func (m \*Elasticsearch) Collection(name string) \*ElasticsearchCollection
            * func (m \*ElasticsearchCollection) Index(id interface{}, data map[string]interface{}) 
            * func (m \*ElasticsearchCollection) Search(key string, value string, page int, pagesize int, cfg ...ElasticsearchSearchingConfig) \*ElasticsearchSearchedResult
            * func (m \*ElasticsearchCollection) Delete(id interface{})
            * func (m \*ElasticsearchCollection) Refresh()
        * func(m \*Elasticsearch) Delete(collection string)
    * func Lock() \*lock
        * func (\*lock) Acquire() 
        * func (\*lock) Release() 
    * func RWLock() *RWLock 
        func (m \*RWLock) RAcquire() 
        func (m \*RWLock) RRelease() 
        func (m \*RWLock) WAcquire() 
        func (m \*RWLock) WRelease() 
    * func AliDNS(string, string) \*alidns
        * func (m \*alidns) Total() (TotalCount int64) 
        * func (m \*alidns) List(PageSize int64, PageNumber int64) (res []alidnsDomainInfo) 
        * func (m \*alidns) Domain(domainName string) \*alidnsDomain
            * func (m \*alidnsDomain) List() (res []alidnsRecord)
            * func (m \*alidnsDomain) Add(recordName string, recordType string, recordValue string) (id string)
            * func (m \*alidnsDomain) Delete(name string, dtype string, value string) 
            * func (m \*alidnsDomain) Modify(recordName string, srcRecordType string, srcRecordValue string, dstRecordName string, dstRecordType string, dstRecordValue string)
    * Chart
        * func LineChartWithTimestampAndNumber([]int64, []float64, string, string, string, string) ????????????x???????????????y??????????????????
        * func LineChartWithNumberAndNumber([]float64, []float64, string, string, string, string) ????????????x???y??????????????????
        * func BarChartWithNameAndNumber([]string, []float64, string, string, string) ????????????x???????????????y??????????????????, ????????????????????????????????????????????????
        * func PieChartWithNameAndNumber([]string, []float64, string, string) ???x????????????y?????????????????????????????????
    * func [CloudflareDNS](#toolscloudflaredns)(string, string) \*cloudflare 
        * func (m \*cloudflare) Add(domain string) cloudflare.Zone
        * func (m \*cloudflare) List() (res []cloudflareDomainInfo)
        * func (m \*cloudflare) Domain(domainName string) \*cloudflareDomain
            * func (m \*cloudflareDomain) List() (res []cloudflareRecord)
            * func (m \*cloudflareDomain) Delete(name string)
            * func (m \*cloudflareDomain) Add(recordName string, recordType string, recordValue string, proxied ...bool) \*cloudflare.DNSRecordResponse
            * func (m \*cloudflareDomain) SetProxied(subdomain string, proxied bool)
            * func (m \*cloudflareDomain) Update(recordName string, recordType string, recordValue string, proxied ...bool)
    * Compress
        * func LzmaCompressString(string) string
        * func LzmaDecompressString(string) string
        * func ZlibCompressString(string) string
        * func ZlibDecompressString(string) string
    * func [Crontab](#toolscrontab)() \*crontab ????????????
        * func (m \*crontab) Add(schedule string, fn interface{}, args ...interface{})
        * func (m \*crontab) Destory()
    * func [GodaddyDNS](#toolsgodaddydns)(string, string) \*godaddy
        * func (m \*godaddy) List() (res []godaddyDomainInfo)
        * func (m \*godaddy) Domain(domainName string) \*godaddyDomain
            * func (m \*godaddyDomain) List() (res []godaddyRecord)
            * func (m \*godaddyDomain) Delete(name string, dtype string, value string)
            * func (m \*godaddyDomain) Modify(recordName string, srcRecordType string, srcRecordValue string, dstRecordType string, dstRecordValue string)
            * func (m \*godaddyDomain) Add(recordName string, recordType string, recordValue string)
    * func [Ini](#toolsini)(...string) \*ini
        * func (m \*ini) Get(SectionKeyDefaultComment ...string) (res string)
        * func (m \*ini) GetInt(key ...string) int
        * func (m \*ini) GetInt64(key ...string) int64
        * func (m \*ini) GetFloat64(key ...string) float64
        * func (m \*ini) Set(SectionKeyValueComment ...string)
        * func (m \*ini) Save(fpath ...string) (exist bool)
    * func [JavascriptVM](#toolsjavascriptvm)() \*javascriptVM
        * func (m \*javascriptVM) Run(javascript string) \*javascriptVM
        * func (m \*javascriptVM) Get(variableName string) string
        * func (m \*javascriptVM) Set(variableName string, variableValue interface{})
        * func (m \*javascriptVM) Isdefined(variableName string) bool
    * func [Matrix](#toolsmatrix)(string) \*matrix
        * func (c \*matrix) Login(username string, password string) string
        * func (c \*matrix) SetToken(userID string, token string) \*matrix
        * func (c \*matrix) SetRoomID(roomID string) \*matrix
        * func (c \*matrix) Send(text string)
    * func [Nats](#toolsnats)(string) \*nats
        * func (m \*nats) Subject(subject string) \*subjectNats
            * func (m \*subjectNats) Publish(message string)
            * func (m \*subjectNats) Subscribe() chan string
            * func (m \*subjectNats) Flush()
    * func Totp(string) \*totp
        * func (m \*totp) Validate(pass string) bool
        * func (m \*totp) Password() string
    * func [Pexpect](#toolspexpect)(string) \*pexpect
        * func (m \*pexpect) Sendline(msg string)
        * func (m \*pexpect) Close()
        * func (m \*pexpect) IsAlive() bool
        * func (m \*pexpect) LogToStdout(enable ...bool)
        * func (m \*pexpect) ExitCode() int
        * func (m \*pexpect) GetLog() string
        * func (m \*pexpect) ClearLog()
    * func [ProgressBar](#toolsprogressbar)(string, int64, ...bool) \*progressBar
        * func (m \*progressBar) Add(num int64)
        * func (m \*progressBar) Set(num int64)
        * func (m \*progressBar) SetTotal(total int64)
        * func (m \*progressBar) Clear()
    * func [PrometheusClient](#toolsprometheusclient)(string) \*prometheusClient
        * func (m \*prometheusClient) Query(query string, time ...float64) (res []prometheusResult) 
    * func [PrometheusMetricServer](#toolsprometheusmetricserver)(listenAddr string, path ...string) \*prometheusMetricServer
        * func (m \*prometheusMetricServer) NewCounter(name string, help string) \*PrometheusCounter 
            * func (m \*PrometheusCounter) Add(num float64) 
        * func (m \*prometheusMetricServer) NewCounterWithLabel(name string, label []string, help string) \*PrometheusCounterVec 
            * func (m \*PrometheusCounterVec) Add(num float64, label map[string]string) 
        * func (m \*prometheusMetricServer) NewGauge(name string, help string) \*PrometheusGauge
            * func (m \*PrometheusGauge) Set(num float64) 
        * func (m \*prometheusMetricServer) NewGaugeWithLabel(name string, label []string, help string) \*PrometheusGaugeVec 
            * func (m \*PrometheusGaugeVec) Set(num float64, label map[string]string) 
    * func [MySQL](#toolsmysql)(string, int, string, string, string, ...DatabaseConfig) \*database
    * func SQLite(string) \*database
        * func (m \*database) Query(sql string, args ...interface{}) []gorose.Data
        * func (m \*database) Close()
        * func (m \*database) Execute(sql string) int64
        * func (m \*database) RenameTable(oldTableName string, newTableNname string)
        * func (m \*database) Tables() (res []string)
        * func (m \*database) CreateTable(tableName string, engineName ...string) \*databaseOrm
        * func (m \*database) Table(tbname string) \*databaseOrm
            * func (m \*databaseOrm) Fields(items ...string) \*databaseOrm
            * func (m \*databaseOrm) Where(key string, operator string, value interface{}) \*databaseOrm
            * func (m \*databaseOrm) WhereIn(key string, value []interface{}) \*databaseOrm
            * func (m \*databaseOrm) WhereNotIn(key string, value []interface{}) \*databaseOrm
            * func (m \*databaseOrm) WhereNull(columnName string) \*databaseOrm
            * func (m \*databaseOrm) WhereNotNull(columnName string) \*databaseOrm
            * func (m \*databaseOrm) OrWhere(key string, operator string, value interface{}) \*databaseOrm
            * func (m \*databaseOrm) OrWhereIn(key string, value []interface{}) \*databaseOrm
            * func (m \*databaseOrm) Orderby(item ...string) \*databaseOrm
            * func (m \*databaseOrm) Limit(number int) \*databaseOrm
            * func (m \*databaseOrm) Get() (res []gorose.Data)
            * func (m \*databaseOrm) Paginate(pagesize int, page int) []gorose.Data
            * func (m \*databaseOrm) First() (res gorose.Data)
            * func (m \*databaseOrm) Find(id int) gorose.Data
            * func (m \*databaseOrm) Count() (res int64)
            * func (m \*databaseOrm) Exists() (res bool)
            * func (m \*databaseOrm) Chunk(limit int, callback func([]gorose.Data) error)
            * func (m \*databaseOrm) BuildSQL() (string, []interface{})
            * func (m \*databaseOrm) Data(data interface{}) \*databaseOrm
            * func (m \*databaseOrm) Offset(offset int) \*databaseOrm
            * func (m \*databaseOrm) InsertGetID() (num int64)
            * func (m \*databaseOrm) Insert() (num int64)
            * func (m \*databaseOrm) Update(data ...interface{}) (num int64)
            * func (m \*databaseOrm) Delete() (num int64)
            * func (m \*databaseOrm) DropTable() int64
            * func (m \*databaseOrm) TruncateTable() (status int64)
            * func (m \*databaseOrm) AddColumn(columnName string, columnType string, defaultValue ...string) \*databaseOrm
            * func (m \*databaseOrm) DropColumn(columnName string) \*databaseOrm
            * func (m \*databaseOrm) AddIndex(columnName ...string) \*databaseOrm
            * func (m \*databaseOrm) IndexExists(columnName ...string) (exists bool)
            * func (m \*databaseOrm) DropIndex(columnName ...string) \*databaseOrm
            * func (m \*databaseOrm) Columns() (res map[string]string)
    * func [RabbitMQ](#toolsrabbitmq)(string, string) \*rabbitConnection
        * func (m \*rabbitConnection) Send(data map[string]string)
        * func (m \*rabbitConnection) Recv() chan map[string]string
    * func RateLimit(int) \*rateLimit
        * func (m \*rateLimit) Take()
    * func [Redis](#toolsredis)(string, int, string, int, ...redisConfig) \*Redis
        * func (m \*Redis) Ping() string
        * func (m \*Redis) Del(key string)
        * func (m \*Redis) Set(key string, value string, ttl ...interface{})
        * func (m \*Redis) Get(key string) \*string
        * func (m \*Redis) GetLock(key string, timeoutsec int) \*RedisLock
            * func (m \*RedisLock) Acquire()
            * func (m \*RedisLock) Release()
    * func [SeleniumLocal](#toolsselenium)() \*selenium 
    * func SeleniumRemote (serverURL string) \*selenium
        * func (c \*Selenium) GetSession() string 
        * func (c \*Selenium) SetSession(SessionID string) 
        * func (c \*Selenium) Get(url string) \*Selenium 
        * func (c \*Selenium) Refresh() \*Selenium 
        * func (c \*Selenium) Title() string 
        * func (c \*selenium) PageSource() string
        * func (c \*selenium) Close()
        * func (c \*selenium) Cookie() (co string)
        * func (c \*selenium) Url() string
        * func (c \*selenium) ScrollRight(pixel int)
        * func (c \*selenium) ScrollLeft(pixel int)
        * func (c \*selenium) ScrollUp(pixel int)
        * func (c \*selenium) ScrollDown(pixel int)
        * func (c \*selenium) ResizeWindow(width int, height int) \*selenium
        * func (c \*selenium) Find(xpath string, nowait ...bool) \*seleniumElement
            * func (c \*seleniumElement) Clear() \*seleniumElement
            * func (c \*seleniumElement) Click() \*seleniumElement
            * func (c \*seleniumElement) Text() string
            * func (c \*seleniumElement) Input(s string) \*seleniumElement
            * func (c \*seleniumElement) Submit() \*seleniumElement
            * func (c \*seleniumElement) PressEnter() \*seleniumElement
            * func (c \*seleniumElement) Attribute(name string) \*String 
    * func [SSH](#toolsssh)(string, string, string, int) \*ssh
        * func (m \*ssh) Close()
        * func (m \*ssh) Exec(cmd string) (output string, status int)
        * func (m \*ssh) PushFile(local string, remote string)
        * func (m \*ssh) PullFile(remote string, local string)
    * func StatikOpen(string) \*statikFile
        * func (m \*statikFile) Readlines() chan string
        * func (m \*statikFile) Readline() string
        * func (m \*statikFile) Close()
        * func (m \*statikFile) Read(num ...int) string
        * func (m \*statikFile) Seek(num int64)
    * func Table(...string) \*table
        * func (m \*table) SetMaxCellWidth(width ...int)
        * func (m \*table) AddRow(row ...interface{})
        * func (m \*table) Render() string
    * func Telegram(AppID int32, AppHash string, config ...TelegramConfig) \*Telegram
        * func (m \*Telegram) Chats() []\*TelegramChat 
            * func (m \*TelegramChat) History(limit int32, offset ...int32) (resmsgs []\*tgMessage)
            * func (m \*TelegramChat) Send(text string)
        * func (m \*Telegram) ResolvePeerByUsername(username string) \*TelegramPeerResolved
            * func (m \*TelegramPeerResolved) History(limit int32, offset ...int32) (resmsgs []\*TelegramMessage)
    * func TelegramBot(string) \*telegramBot
        * func (m \*telegramBot) SetChatID(chatid int64) \*telegramBot
        * func (m \*telegramBot) SendFile(path string) tgbotapi.Message
        * func (m \*telegramBot) SendImage(path string) tgbotapi.Message
        * func (m \*telegramBot) SendVideo(path string) tgbotapi.Message
        * func (m \*telegramBot) SendAudio(path string) tgbotapi.Message
        * func (m \*telegramBot) Send(text string, cfg ...tgMsgConfig) tgbotapi.Message
    * func Telegraph(string) \*telegraph
        * func (m \*telegraph) Post(title string, content string) \*telegraphPageInfo
    * func URL(string) \*url
        * func (u \*url) Parse() \*urlComponents
        * func (u \*url) Encode() string
        * func (u \*url) Decode() string
    * func TTLCache(intervalInSecond interface {}) \*ttlCache
        * func (m \*ttlCache) Set(key string, value string)
        * func (m \*ttlCache) Get(key string) string
        * func (m \*ttlCache) Exists(key string) bool
        * func (m \*ttlCache) Count() int
    * func [VNC](#toolsvnc)(string, ...VNCCfg) \*vnc
        * func (m \*vnc) Close()
        * func (m \*vnc) Move(x, y int) \*vnc
        * func (m \*vnc) Click() \*vnc
        * func (m \*vnc) RightClick() \*vnc
        * func (m \*vnc) Input(s string) \*vnc
        * func (m \*vnc) Key() \*vncKey
            * func (m \*vncKey) Enter() \*vncKey
            * func (m \*vncKey) CtrlA() \*vncKey
            * func (m \*vncKey) CtrlC() \*vncKey
            * func (m \*vncKey) CtrlV() \*vncKey
            * func (m \*vncKey) CtrlZ() \*vncKey
            * func (m \*vncKey) CtrlX() \*vncKey
            * func (m \*vncKey) CtrlF() \*vncKey
            * func (m \*vncKey) CtrlD() \*vncKey
            * func (m \*vncKey) CtrlS() \*vncKey
            * func (m \*vncKey) CtrlR() \*vncKey
            * func (m \*vncKey) CtrlE() \*vncKey
            * func (m \*vncKey) Delete() \*vncKey
            * func (m \*vncKey) Tab() \*vncKey
    * func WebSocket(string) \*websocket
        * func (c \*websocket) Send(text string)
        * func (c \*websocket) Recv(timeout ...int) string
        * func (c \*websocket) Close()
    * func [Xlsx](#toolsxlsx)(string) \*xlsx
        * func (c \*xlsx) Save()
        * func (c \*xlsx) Close()
        * func (c \*xlsx) GetSheet(name string) \*xlsxSheet
            * func (c \*xlsxSheet) Get(coordinate string) string
            * func (c \*xlsxSheet) Set(coordinate string, value string) \*xlsxSheet
    * func [XPath](#toolsxpath)(string) \*xpath
        * func (m \*xpath) First(expr string) (res *xpath)
        * func (m \*xpath) Find(expr string) (res []*xpath)
        * func (m \*xpath) Text() string
        * func (m \*xpath) GetAttr(attr string) string
        * func (m \*xpath) Html() string
    * Sysinfo
        * Host
            * func Info() types.HostInfo
	        * func Memory() \*types.HostMemoryInfo
	        * func CPUTimes() types.CPUTimes
        * func Process(pid int) \*sysinfoProcess
            * func (p \*sysinfoProcess) Info() types.ProcessInfo {
            * func (p \*sysinfoProcess) Memory() types.MemoryInfo {
            * func (p \*sysinfoProcess) User() types.UserInfo {
            * func (p \*sysinfoProcess) Parent() \*sysinfoProcess {
            * func (p \*sysinfoProcess) CPUTimes() types.CPUTimes {
    * Queue(datadir string) (q \*queue) 
        * func (m \*queue) Close() 
        * func (m \*queue) Destroy() 
        * func (m \*queue) New(queueName ...string) \*namedQueue 
            * func (m \*namedQueue) Size() int64 
            * func (m \*namedQueue) Get(nonblock ...bool) string 
            * func (m \*namedQueue) Put(value string) 
    * Jieba() \*Jieba ????????????
        * func (m \*Jieba) Close()
        * func (m \*Jieba) Cut(s string) []string 
        * func (m \*Jieba) AddWord(text string)
* Random
    * func Int(min, max int64) int64
    * func Choice(array interface{}) interface{}
    * func String(length int, charset ...string) string
* Re
    * func FindAll(pattern string, text string, multiline ...bool) [][]string
    * func Replace(pattern string, newstring string, text string) string
* Socket
    * [KCP](#socketkcp)
        * func Listen(string, int, string, string) \*kcpServerSideListener
            * func (m \*kcpServerSideListener) Accept() chan *kcp.UDPSession
        * func Connect(string, int, string, string) \*kcp.UDPSession
    * Smux
        * func ServerWrapper(io.ReadWriteCloser, ...SmuxConfig) \*smuxServerSideListener
            * func (m *smuxServerSideListener) Accept() chan *smuxServerSideConnection 
                * func (m *smuxServerSideConnection) Send(data map[string]string, timeout ...int) 
                * func (m *smuxServerSideConnection) Recv(timeout ...int) (data map[string]string) 
                * func (m *smuxServerSideConnection) Close()
        * func ClientWrapper(io.ReadWriteCloser, ...SmuxConfig) \*smuxClientSideSession
            * func (m *smuxClientSideSession) Connect() *smuxClientSideConnection 
            * func (m *smuxClientSideSession) Close() 
                * func (m *smuxClientSideConnection) Send(data map[string]string, timeout ...int) 
                * func (m *smuxClientSideConnection) Recv(timeout ...int) (data map[string]string) 
                * func (m *smuxClientSideConnection) Close() 
    * SSL
        * func [Listen](#socketssllisten)(string, int, string, string) \*tcpServerSideListener
            * func (m *tcpServerSideListener) Accept() chan *TcpServerSideConn 
            * func (m *tcpServerSideListener) Close() 
                * func (m *TcpServerSideConn) Close() 
                * func (m *TcpServerSideConn) Send(str string) 
                * func (m *TcpServerSideConn) Recv(buffersize int) string 
        * func [ServerWrapper](#socketsslserverwrapper)(net.Conn, string, string) \*tcpServerSideConn
            * func (m *tcpClientSideConn) Send(str string, timeout ...int) 
            * func (m *tcpClientSideConn) Recv(buffersize int, timeout ...int) string 
            * func (m *tcpClientSideConn) Close() 
        * func Connect(string, int, ...sslCfg) \*sslClientSideConn
        * func ClientWrapper(net.Conn, ...sslCfg) \*sslClientSideConn
            * func (m *sslClientSideConn) Send(str string) 
            * func (m *sslClientSideConn) Recv(buffersize int) string 
            * func (m *sslClientSideConn) Close() 
    * TCP
        * func [Listen](#sockettcplisten)(string, int) \*tcpServerSideListener
            * func (m *tcpServerSideListener) Accept() chan *TcpServerSideConn 
            * func (m *tcpServerSideListener) Close() 
                * func (m *TcpServerSideConn) Close() 
                * func (m *TcpServerSideConn) Send(str string) 
                * func (m *TcpServerSideConn) Recv(buffersize int) string 
        * func [Connect](#sockettcpconnect)(string, int, ...int) \*tcpClientSideConn
            * func (m *tcpClientSideConn) Send(str string, timeout ...int) 
            * func (m *tcpClientSideConn) Recv(buffersize int, timeout ...int) string 
            * func (m *tcpClientSideConn) Close() 
    * Unix
        * func [Listen](#sockettcplisten)(path string) \*tcpServerSideListener
            * func (m *tcpServerSideListener) Accept() chan *TcpServerSideConn 
            * func (m *tcpServerSideListener) Close() 
                * func (m *TcpServerSideConn) Close() 
                * func (m *TcpServerSideConn) Send(str string) 
                * func (m *TcpServerSideConn) Recv(buffersize int) string 
        * func [Connect](#sockettcpconnect)(path string) \*tcpClientSideConn
            * func (m *tcpClientSideConn) Send(str string, timeout ...int) 
            * func (m *tcpClientSideConn) Recv(buffersize int, timeout ...int) string 
            * func (m *tcpClientSideConn) Close() 
    * [UDP](#socketudp)
        * func Listen(string, int) udpServerSideConn
            * func (m *udpClientSideConn) Send(str string) 
            * func (m *udpClientSideConn) Close() 
            * func (m *udpClientSideConn) Recv(buffersize int) string 
        * func Connect(string, int) udpClientSideConn
            * func (m *udpServerSideConn) Recvfrom(buffersize int, timeout ...int) (string, *net.UDPAddr) 
            * func (m *udpServerSideConn) Sendto(data string, address *net.UDPAddr, timeout ...int) 
            * func (m *udpServerSideConn) Close() 
* String
    * func (s \*string) Get() string
    * func (s \*string) Sub(start, end int) \*string
    * func (s \*string) Has(substr string) bool
    * func (s \*string) Len() int
    * func (s \*string) Reverse() string
    * func (s \*string) Chunk(length int) (res []string)
    * func (s \*string) Utf8Len() int
    * func (s \*string) Repeat(count int) \*string
    * func (s \*string) Shuffle() \*string
    * func (s \*string) Index(substr string) int
    * func (s \*string) Replace(old, new string) \*string
    * func (s \*string) Upper() \*string
    * func (s \*string) Lower() \*string
    * func (s \*string) Join(pieces []string) \*string
    * func (s \*string) Strip(characterMask ...string) \*string
    * func (s \*string) Split(sep ...string) []string
    * func (s \*string) Count(substr string) int
    * func (s \*string) EndsWith(substr string) (res bool)
    * func (s \*string) StartsWith(substr string) (res bool)
    * func (s \*string) Splitlines(strip ...bool) []string
    * func (s \*string) In(str string) bool
    * func (s \*string) LStrip(characterMask ...string) \*string
    * func (s \*string) RStrip(characterMask ...string) \*string
    * func (ss\*string) Isdigit() bool
    * func (s \*string) HasChinese() bool
    * func (s \*string) Filter(charts ...string) \*string
    * func (s \*string) RemoveHtmlTag() \*string
    * func (s \*string) RemoveNonUTF8Character() \*string
    * func (s \*string) DetectLanguage() string
    * func (s \*string) IsAscii() bool 
    * func (s \*string) RegexFindAll(pattern string, multiline ...bool) (res [][]\*string) 
    * func (s \*string) RegexReplace(pattern string, newstring string) \*string 
* Time:
    * func Now() float64
    * func TimeDuration(interface {}) time.Duration
    * func FormatDuration(int64) string
    * func Sleep(interface {})
    * func Strptime(string, string) int64
    * func Strftime(string, interface {}) string
* [Argparser](#argparser)(description string) \*argparseIni
    * func (m \*argparseIni) Get(section, key, defaultValue, comment string) (res string)
    * func (m \*argparseIni) GetInt(section, key, defaultValue, comment string) int
    * func (m \*argparseIni) GetInt64(section, key, defaultValue, comment string) int64
    * func (m \*argparseIni) GetFloat64(section, key, defaultValue, comment string) float64
    * func (m \*argparseIni) GetBool(section, key, defaultValue, comment string) bool
    * func (m \*argparseIni) Save(fpath ...string) (exist bool)
    * func (m \*argparseIni) GetHelpString() (h string)
    * func (m \*argparseIni) ParseArgs() \*argparseIni
* Base64
    * func Encode(str string) string
    * func Decode(str string) string
* Binary
    * func Map2bin(m map[string]string) string
    * func Bin2map(s string) (res map[string]string)
* Cmd
    * func GetOutput(command string, timeoutSecond ...interface{}) string
    * func GetStatusOutput(command string, timeoutSecond ...interface{}) (int, string)
    * func GetOutputWithShell(command string, timeoutSecond ...interface{}) string 
    * func GetStatusOutputWithShell(command string, timeoutSecond ...interface{}) (int, string)
    * func Tail(command string) chan string
    * func Exists(cmd string) bool
    * func Which(cmd string) (path string)
* Crypto
    * func Xor(data, key string) string
    * func Aes(key string) \*aes
        * func (a \*aes) Encrypt(plaintext string) string 
        * func (a \*aes) Decrypt(ciphertext string) string 
    * func ChaCha20Poly1305(key string) *chacha20poly1305
        * func (m \*chacha20poly1305) Encrypt(plantext string) (ciphertext string)
        * func (m \*chacha20poly1305) Decrypt(ciphertext string) (plaintext string)
* File(filePath string) \*file
    * func (f \*file) Time() \*fileTime
    * func (f \*file) Stat() os.FileInfo
    * func (f \*file) Size() int64
    * func (f \*file) Touch()
    * func (f \*file) Chmod(mode int) bool
    * func (f \*file) Chown(uid, gid int) bool
    * func (f \*file) Mtime() int64
    * func (f \*file) Unlink()
    * func (f \*file) Copy(dest string)
    * func (f \*file) Move(newPosition string)
* Open(args ...string) \*fileIO
    * func (m \*fileIO) Readlines() chan string
    * func (m \*fileIO) Readline() string
    * func (m \*fileIO) Close()
    * func (m \*fileIO) Write(str interface{}) \*fileIO
    * func (m \*fileIO) Read(num ...int) string
    * func (m \*fileIO) Seek(num int64)
* Funcs
    * func Nslookup(name string, querytype string, dnsService ...string) [][]string
    * func FakeName() string
    * func FileType(fpath string) string
    * func Inotify(path string) chan *fsnotifyFileEvent
    * func IPLocation(ip string, dbpath ...string) \*ipLocationInfo
    * func HightLightHTMLForCode(code string, codeType ...string) (html string)
    * func Markdown2html(md string) string
    * func CPUUsagePerProgress() (res map[int64]progressCPUUsage)
    * func ResizeImg(srcPath string, dstPath string, width int, height ...int)
    * func GetRSS(url string, config ...rssConfig) \*gofeed.Feed
    * func GbkToUtf8(s string) string
    * func Utf8ToGbk(s string) string
    * func GetSnowflakeID(nodeNumber ...int) int64
    * func GetRemoteServerSSLCert(host string, port ...int) []*x509.Certificate
    * func Tailf(path string, startFromEndOfFile ...bool) chan *tail.Line
    * func BaiduTranslateAnyToZH(text string) string
    * func ParseUserAgent(UserAgent string) ua.UserAgent
    * func Wget(url string, cfg ...WgetCfg) (filename string)
    * func Whois(s string, servers ...string) string
    * func IpInNet(ip string, Net string, mask ...string) bool
    * func Int2ip(ipnr int64) string
    * func Ip2int(ipnr string) int64
    * func Zh2PinYin(zh string) (ress []string)
    * func Fmtsize(num uint64) string
    * func Sniffer(interfaceName string, filterString string, promisc ...bool) chan *networkPacket 
	* func ReadPcapFile(pcapFile string) chan *networkPacket
* Hash
    * func Md5sum(str string) string
    * func Md5File(path string) string
    * func Sha1sum(str string) string
    * func Sha1File(path string) string
* Html
    * func Encode(str string) string
    * func Decode(str string) string
* Http
    * func Head(uri string, args ...interface{}) httpResp
    * func PostFile(uri string, filePath string, args ...interface{}) httpResp
    * func PostRaw(uri string, body string, args ...interface{}) httpResp
    * func PostJSON(uri string, json interface{}, args ...interface{}) httpResp
    * func Post(uri string, args ...interface{}) httpResp
    * func Get(uri string, args ...interface{}) httpResp
    * func PutJSON(uri string, json interface{}, args ...interface{}) httpResp
    * func Put(uri string, args ...interface{}) httpResp
    * func PutRaw(uri string, body string, args ...interface{}) httpResp
    * func Delete(uri string, args ...interface{}) httpResp
* Json
    * func Dumps(v interface{}, pretty ...bool) string
    * func Loads(str string) map[string]interface{}
    * func Valid(j string) bool
    * func [Yaml2json](#jsonyaml2json)(y string) string ??????yaml???json
    * func [Json2yaml](#jsonjson2yaml)(j string) string ??????json???yaml, ?????????????????????????????????
    * func Format(js string) string ?????????json?????????, ????????????format???????????????????????????
    * func [XPath](jsonxpath)(string) \*xpathJson
        * func (m \*xpathJson) Exists(expr string) bool
        * func (m \*xpathJson) First(expr string) (res *xpathJson)
        * func (m \*xpathJson) Find(expr string) (res []*xpathJson)
        * func (m \*xpathJson) Text() *string
* Math
    * func Abs(number float64) float64
    * func Sum(array interface{}) (sumresult float64)
    * func Average(array interface{}) (avgresult float64)
    * func DecimalToAny(num, n int64) string
    * func AnyToDecimal(num string, n int64) int64
* Os
    * func Mkdir(filename string)
    * func Getcwd() string
    * func Exit(status int)
    * func Touch(filePath string)
    * func Chmod(filePath string, mode int) bool
    * func Chown(filePath string, uid, gid int) bool
    * func Copy(filePath, dest string)
    * func Rename(filePath, newPosition string)
    * func Move(filePath, newPosition string)
    * Path 
        * func Exists(path string) bool
        * func IsFile(path string) bool
        * func IsDir(path string) bool
        * func Basename(path string) string
        * func Basedir(path string) string
        * func Dirname(path string) string
        * func Join(args ...string) string
        * func Abspath(path string) string
        * func IsSymlink(path string) bool
    * func System(command string, timeoutSecond ...interface{}) int
    * func SystemWithShell(command string, timeoutSecond ...interface{}) int
    * func Hostname() string
    * func Envexists(varname string) bool
    * func Getenv(varname string) string
    * func Walk(path string) chan string
    * func Listdir(path string) (res []string)
    * func SelfDir() string
    * func TempFilePath() string
    * func TempDirPath() string
    * func Getuid() int
    * func ProgressAlive(pid int) bool
    * func GoroutineID() int64
    * func Unlink(filename string)
    * Stdin
        * func Readlines() chan *String
        * func Readline() *String
        * func Read(num ...int) *String
    * Stdout
        * Write(str interface{}) *fileIO
Others

* func [Open](#open)(filePath string) *fileIO
* func [Try](#try)(f func(), trycfg ...TryConfig) exception
* Lg *log 
    * func (m *log) SetLevel(level string) 
    * func (m *log) GetLevel() string 
    * func (m *log) SetLogFile(path string, maxLogFileCount int, logFileSizeInMB ...int) 
    * func (m *log) Error(args ...interface{}) 
    * func (m *log) Warn(args ...interface{}) 
    * func (m *log) Info(args ...interface{}) 
    * func (m *log) Trace(args ...interface{}) 
    * func (m *log) Debug(args ...interface{}) 
* func Chr(ascii int) string 
* func Ord(char string) int 
* func Repr(obj interface{}) string 
* func Print(data ...interface{}) int 
* func Printf(format string, data ...interface{}) int 
* func Sprint(data ...interface{}) string 
* func Range(num ...int) []int 
* func Typeof(v interface{}) string 

# Example

## Json.Yaml2Json

```go
func main() {
	j := `code: 0
mesg: Get Domains Successful
result:
  active: true
  domains:
  - ishomee.com
  - dx2cone1.xyz
  - zhiyunxianghe.com
success: true`
	fmt.Println(Json.Yaml2json(j))
}
```

## Json.Json2yaml

```go
func main() {
	j := `{"code":0,"mesg":"Get Domains Successful","result":{"active":true,"domains":["ishomee.com","dx2cone1.xyz","zhiyunxianghe.com"]},"success":true}`
	Print(Json.Json2yaml(j))
}
```

## Tools.Chart.LineChartWithTimestampAndNumber

```go
func main() {
	x := "2020-04-21,2020-05-09,2020-05-11,2020-05-14,2020-05-15,2020-05-17,2020-05-19,2020-05-20,2020-05-21,2020-05-22,2020-05-24,2020-05-25,2020-05-26,2020-05-27,2020-05-28,2020-05-29,2020-05-30,2020-05-31,2020-06-01,2020-06-02,2020-06-03,2020-06-04,2020-06-05,2020-06-06,2020-06-07,2020-06-08,2020-06-09,2020-06-10,2020-06-11,2020-06-12,2020-06-13,2020-06-14,2020-06-15,2020-06-16,2020-06-17,2020-06-18,2020-06-19,2020-06-20,2020-06-21,2020-06-22,2020-06-23,2020-06-25,2020-06-26,2020-06-27,2020-06-28,2020-06-29,2020-06-30,2020-07-01,2020-07-02,2020-07-03,2020-07-04,2020-07-05,2020-07-06,2020-07-07,2020-07-08,2020-07-09,2020-07-10,2020-07-11,2020-07-12,2020-07-13,2020-07-14,2020-07-15,2020-07-16,2020-07-17,2020-07-18,2020-07-19,2020-07-20,2020-07-21,2020-07-22,2020-07-23,2020-07-24,2020-07-25,2020-07-26,2020-07-27,2020-07-28,2020-07-29,2020-07-30,2020-07-31,2020-08-01,2020-08-02,2020-08-03,2020-08-04,2020-08-05,2020-08-06,2020-08-07,2020-08-08,2020-08-09,2020-08-10,2020-08-11,2020-08-12,2020-08-13"
	var xx []int64
	for _, i := range String(x).Split(",") {
		xx = append(xx, Time.Strptime("%Y-%m-%d", i))
	}

	y := "100,100,500,100,100,100,200,700,200,700,300,400,900,1100,1400,900,3004,908,1460,4400,1500,2000,2950,2150,2750,7150,3850,4050,3900,4800,4200,7400,6700,6150,7400,7250,7550,9800,8900,5300,1700,1000,800,1500,1150,1300,2060,3820,4852,4320,4960,5160,2610,2640,3300,1770,2690,2020,2360,2050,1580,1410,1080,850,1540,1410,1460,1540,1620,1370,3328,3898,2218,2238,2398,2038,1700,750,1100,1700,1650,1340,950,2270,540,890,1390,1900,1580,2450,1680"
	var yy []float64
	for _, i := range String(y).Split(",") {
		yy = append(yy, Float64(i))
	}

	drawLineChartWithTimeSeries(xx, yy, "??????", "??????", "???????????????", "output.png") // ???????????????????????????png????????????
}
```

## Tools.Crontab

```go
func func1(arg string) {
	Print(arg)
}

func main() {
	c := Tools.Crontab()
	
	c.Add("*/1 * * * *", func() {
		Print(now())
	})

	c.Add("*/1 * * * *", func(param1 string, param2 string) {
		Print(Time.Now(), "with param: "+param1+" and "+param2)
	}, "paramValue1", "paramValue2")

	c.Add("00 16 * * *", func1, "args1")

	select {}
}
```

??????????????????

```
*     *     *     *     *        

^     ^     ^     ^     ^
|     |     |     |     |
|     |     |     |     +----- day of week (0-6) (Sunday=0)
|     |     |     +------- month (1-12)
|     |     +--------- day of month (1-31)
|     +----------- hour (0-23)
+------------- min (0-59)
```

?????????

1. `* * * * *` run on every minute
2. `10 * * * *` run at 0:10, 1:10 etc
4. `10 15 * * *` run at 15:10 every day
5. `* * 1 * *` run on every minute on 1st day of month
6. `0 0 1 1 *` Happy new year schedule
7. `0 0 * * 1` Run at midnight on every Monday
8. `* 10,15,19 * * *` run at 10:00, 15:00 and 19:00
9. `1-15 * * * *` run at 1, 2, 3...15 minute of each hour
10. `0 0-5,10 * * *` run on every hour from 0-5 and in 10 oclock
11. `*/2 * * * *` run every two minutes
12. `10 */3 * * *` run every 3 hours on 10th min
13. `0 12 */2 * *` run at noon on every two days
14. `1-59/2 * * *` * run every two minutes, but on odd minutes

## Tools.JavascriptVM

```go
func main() {
  s := "a = 1;console.log(b);"
  vm := Tools.JavascriptVM()
  vm.Set("b", "2")
  vm.Run(s)
  Print(vm.Get("a"))
}
```

## Tools.Matrix

```go
  // ??????homeserver???url, ?????????????????????????????????id, ???????????????????????????????????????
  cli = Tools.Matrix("https://example.com").SetRoomID("!AquOdzAnBLIQvfPkan:example.com")

  // ??????????????????????????????, ???????????????token, ????????????token?????????, ????????????????????????token
  token := cli.Login("bot", "123456") // ???????????????????????????token
  // ?????????????????????token, ??????synapse???token????????????????????????
  cli.SetToken("@bot:example.com", "syt_Ym90_iHallJrSVvDLFCfnvnZZ_4a2WKt")
  
  cli.Send(msg)
}
```

## Tools.Nats

```go
func main() {
	server := "nat://nats.nats.svc.cluster.local"
	subj := Tools.Nats(server).Subject("mysubject")

	go func() {
		for msg := range subj.Subscribe() {
			Lg.trace(msg)
		}
	}()

	for {
		sleeptime := Rand.Int(1, 3)
		Time.Sleep(sleeptime)
		subj.Publish("sleep for " + str(sleeptime) + " second(s) just now")
	}
}
```

## Tools.Pexpect

?????????py????????????????????????

```python
try:
    while True:
        a = input("Enter something: ")
        Print("You Enter: ", a)
except:
	pass
```

```go
func main() {
	i := 0
	p := Tools.Pexpect("python3 test.py")
	defer p.Close()
	p.LogToStdout()
	for p.IsAlive() {
		Time.Sleep(1)
		if strings.Contains(p.GetLog(), "Enter something:") {
			p.ClearLog()
            p.Sendline(Str(i))
		}
		i++
		if i >= 5 {
			p.Close()
			break
		}
	}
	Print("Exit code:", p.ExitCode())
}
```

## Tools.Ini

```go
func main() {
	i := Tools.Ini("c.ini")                               // ????????????????????????????????????ini?????????
	Print(i.Get("section", "key", "value", "comment")) // ?????????????????????????????????????????????????????????????????????????????????????????????????????????????????????
	Print(i.Save())                                    // ???????????????????????????????????????true??????????????????false
}
```

## Tools.ProgressBar

```go
func main() {
	bar := Tools.ProgressBar("example bar", 100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		Time.Sleep(5)
		if i == 80 {
			bar.SetTotal(300) // ??????????????????
		}
	}

	for i := 0; i < 200; i++ {
		bar.Add(1)
		Time.Sleep(5)
	}
}
```

## Tools.PrometheusClient

????????????????????????, ??????label???, ??????????????????(?????????????????????, ????????????)

```go
func main() {
	p := Tools.PrometheusClient("http://localhost:9090")
	pr := p.Query("sum_over_time(channel_register_count_in_5_minutes{channel=\"1\"}[1h]) / sum_over_time(channel_inpour_count_in_5_minutes[1h]) < 100")
	Lg.Debug(pr)
}
```

????????????

```go
[]main.prometheusResult{
  main.prometheusResult{
    Label: map[string]string{
      "instance":  "10.0.0.1:9100",
      "job":       "my-service-svc",
      "namespace": "default",
      "pod":       "my-service-332332234-1231232",
      "service":   "my-service-svc",
      "channel":   "1",
      "endpoint":  "my-endpoint",
    },
    Value: 44.651376,
  },
}
```

## Tools.PrometheusMetricServer

```go
func TestPrometheusMetricServer(t *testing.T) {
	p := getPrometheusMetricServer("0.0.0.0:9301")
	c := p.NewCounter("test_counter", "this is a test Counter metrics")
	g := p.NewGauge("test_gauge", "this is a test Gauge metrics")

	go func() {
		for {
			c.Add(1)
			Time.Sleep(1)
		}
	}()

	for {
		g.Set(Float64(Random.Int(0, 10000)))
		Time.Sleep(1)
	}
}
```

## Tools.MySQL

?????????

```go
func main() {
	db := Tools.MySQL("mysql-svc", 3306, "root", "", "test")
	//db := getSQLite("data.db")
	// ??????
	db.CreateTable("tbName")
	// ?????????
	db.Table("tbName").AddColumn("intType", "int")      // bigint
	db.Table("tbName").AddColumn("floatType", "float")  // double
	db.Table("tbName").AddColumn("vcharType", "string") // VARCHAR(512)
	db.Table("tbName").AddColumn("textType", "text")    // LONGTEXT
	// ?????????
	db.Table("tbName").DropColumn("intType") // SQLite?????????
	// ????????????
	db.Table("tbName").AddIndex("floatType")
	db.Table("tbName").AddIndex("floatType", "vcharType")
	// ????????????
	db.Table("tbName").DropIndex("floatType")
	db.Table("tbName").DropIndex("floatType", "vcharType")

	// ?????????????????????
	db.CreateTable("usercodes").
		AddColumn("usercode", "string").
		AddColumn("start", "int").
		AddColumn("duration", "int").
		AddIndex("usercode")

	// ??????????????????
	pg := getSQLite(":memory:").
		CreateTable("progress").
		AddColumn("pid", "float").
		AddColumn("name", "string").
		AddColumn("cpu", "float").
		AddColumn("cmd", "string").
		AddColumn("start", "int").
		AddColumn("end", "int").
		AddColumn("notified", "string").
		AddIndex("cpu")
}
```

????????????

```go
func main() {
	db := getMySQL("mysql-svc", 3306, "root", "", "test")

	u := db.Table("user") // ????????????u?????????????????????????????????????????????, ???????????????????????????????????????db.Table()

	// select
	reses := db.Table("user").Fields("id", "name", "age").Where("age", ">", 0).Orderby("id desc").limit(2).get()
	fmt.Println(reses) // [map[age:6 id:5 name:cat ] map[age:5 id:4 name:monkey]]
	// ???????????????
	fmt.Println(reses[0]) // map[age:6 id:5 name:cat ]
	// ?????????????????????????????????
	for _, r := range reses {
		fmt.Println(r["name"])
	}

	// ?????????????????????
	res := db.Table("user").Where("age", ">", 0).Orderby("id", "desc").First()
	fmt.Println(res) // map[age:6 id:5 name:cat ]
	Print(len(res)) // 0, ??????????????? 

	count := db.Table("user").Where("age", ">", 0).count()
	fmt.Println(count) // 5

	// ???????????????
	db.Table("user").Fields("id", "name", "age").Where("age", ">", 0).Orderby("id").chunk(2, func(data []gorose.Data) error {
		fmt.Println("In Chunk: ", data)
		// In Chunk:  [map[age:1 id:1 name:cookie] map[age:2 id:2 name:ares]]
		// In Chunk:  [map[age:3 id:3 name:div] map[age:5 id:4 name:monkey]]
		// In Chunk:  [map[age:6 id:5 name:cat ]]
		return nil
	})

	// ??????????????????
	var data = map[string]interface{}{"age": 17, "name": "it3"}
	id := db.Table("user").Data(data).InsertGetID()
	fmt.Println(id) // 6??? ???????????????id

	// ??????????????????
	var multiData = []map[string]interface{}{{"age": 18, "name": "it4"}, {"age": 19, "name": "it5"}}
	re := db.Table("user").Data(multiData).Insert()
	fmt.Println(re) // 2 ??? ?????????????????????

	// ????????????
	re = db.Table("user").Where("id", "=", 1).OrWhere("age", ">", 5).Data(map[string]interface{}{"age": 29, "name": "new Name"}).update()
	fmt.Println(re) // 5, ?????????????????????

	// ????????????
	re = db.Table("user").Where("id", "=", 1).Delete()
	fmt.Println(re) // 1, ?????????????????????

	rese := db.Query("select count(id) as `count`, `age` from `user` group by `age` order by `count` desc")
	fmt.Println(rese) // [map[age:29 count:4] map[age:2 count:1] map[age:3 count:1] map[age:5 count:1]]

	ress := db.Execute("delete from `user` where `age` = 29")
	fmt.Println(ress) // 4

	sql, param := db.Table("user").Fields("id", "name", "age").Where("age", ">", 0).Orderby("id desc").limit(2).buildSQL()
	fmt.Println(sql, param) // SELECT `id`,`name`,`age` FROM `user` WHERE `age` > ? ORDER BY id desc LIMIT 2 [0]
}
```

?????????????????????SQL

```sql
CREATE TABLE `user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `age` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 0 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

INSERT INTO `user` VALUES (1, 'cookie', 1);
INSERT INTO `user` VALUES (2, 'ares', 2);
INSERT INTO `user` VALUES (3, 'div', 3);
INSERT INTO `user` VALUES (4, 'monkey', 5);
INSERT INTO `user` VALUES (5, 'cat ', 6)
```

# jsonDumps() - ??????????????????json?????????

```go
func main() {
	a := jsonMap{
		"a": "b",
		"c": "d",
		"e": jsonMap{"f": "g"},
		"h": jsonArr{1, "k"},
	}
	j := jsonDumps(a) // {"a":"b","c":"d","e":{"f":"g"},"h":[1,"k"]}
	Print(j)
	k := jsonLoads(j)
	Print(k)      // map[a:b c:d e:map[f:g] h:[1 k]]
	Print(k["a"]) // b
}
```

## Tools.RabbitMQ

```go
func main() {
	go func() {
		rb := Tools.RabbitMQ("amqp://guest:guest@rabbitmq-svc:5672/", "default")
		rb.Send(map[string]string{"data": "Test Message"})
	}()

	go func() {
		rb := Tools.RabbitMQ("amqp://guest:guest@rabbitmq-svc:5672/", "default")
		msg := <-rb.Recv()
		Lg.debug(msg)
	}()

	select {}
}
```

## Tools.Redis

```go
func main() {
	rdb := Tools.Redis("redis-svc", 6379, "", 1)
	rdb.Set("key", "value")
	fmt.Println(*rdb.Get("key")) // ??????key???????????????nil, ???????????????value????????????string?????????
	rdb.Set("ttl", "delete after 1 second", 1)
	rdb.Set("ttl2", "delete after 0.5 second", 0.5)
	rdb.Del("key")
}
```

## Tools.Selenium

```go
func main() {
	// ??????chromedriver???PATH?????????????????????????????????, ?????????????????????, ????????????, ?????????????????????
	sn := Tools.Selenium("https://example.com/auth/login")
	defer sn.Close() // ???????????????????????????

	// ??????
	Lg.Trace("?????????")
	// ?????????select??????????????????option???xpath
	sn.First(`/html/body/div/div[1]/div[1]/div[2]/form/div[3]/div/select/option[2]`).click()
	Lg.Trace("???????????????")
	sn.First(`//*[@id="login"]`).clear().input("user") // ??????????????????
	Lg.Trace("????????????")
	sn.First(`//*[@id="password"]`).clear().input("pass")
	Lg.Trace("????????????")
	sn.First(`/html/body/div[1]/div[1]/div[1]/div[2]/form/center/div/input`).click()

	// ????????????????????????
	Lg.Trace("????????????????????????")
	sn.First(`//*[@id="gotomemberinfo"]`).input("uid12345").pressEnter() // ???????????????
	vipLevel := sn.First(`/html/body/div[2]/div/div[1]/div[1]/div[3]/div/div/div/div[1]/fieldset[2]/div/div[1]/div[1]`).text()
	Lg.Trace("VIP??????:", is.VipLevel)

	Lg.Trace("??????????????????")
	url := sn.url() // ???????????????url
	uid := String(url).Split("/")[len(String(url).Split("/"))-1]
	j := Http.Get("https://example.com/player_management/getDetails/"+uid, httpHeader{
		"cookie": sn.Cookie(), // cookie????????????
	}).content
	//Lg.Trace(j)
	jj := Tools.JsonXPath(j)
	length = jj.First("//Details/total").Text()
	Lg.Trace("????????????:", is.length)


	select {}
}
```

## Tools.SSH

```go
func main() {
	s := Tools.SSH("root", "root", "192.168.152.19", 22)
	Print(s.Exec("id"))
	s.PullFile("anaconda-ks.cfg", "tmp.file")
	s.PushFile("main.go", "main.go")
}
```

## Tools.VNC

```go
func main() {
	v := Tools.VNC("192.168.3.49:5900", VNCCfg{Password: "r"}) // vncCfg??????
	v.Move(237, 570).click()
	v.Input("Hello world!\nHHHHHH") // ???????????????, ??????????????????????????????. 

    v.Key().CtrlC()

	// Ctrl-C
	v.VC.KeyEvent(vncNonAsciiKeyMap.Control, true)
	v.VC.KeyEvent(vncAsciiKeyMap["c"][0], true)
	v.VC.KeyEvent(vncAsciiKeyMap["c"][0], false)
	v.VC.KeyEvent(vncNonAsciiKeyMap.Control, false)
}
```

## Tools.Xlsx

```go
func main() {
	// ??????????????????, ??????????????????????????????
	xlsx := Tools.Xlsx("Book1.xlsx")
	// ??????????????????, ????????????????????????????????????
	sheet1 := xlsx.GetSheet("sheet4")
	// ??????B??????14????????????value, ???C??????3????????????key
	sheet1.Set("B14", "value").Set("C3", "key")
	// ?????????
	fmt.Println(sheet1.Get("C3"))
	// ??????
	xlsx.Close()
}
```

## Tools.XPath

?????????xml??????

```xml
<bookstore>

<book category="cooking">
  <title lang="en">Everyday Italian</title>
  <author>Giada De Laurentiis</author>
  <year>2005</year>
  <price>30.00</price>
</book>

<book category="children">
  <title lang="zh-cn">Harry Potter</title>
  <author>J K. Rowling</author>
  <year>2005</year>
  <price>29.99</price>
</book>

<book category="web">
  <title lang="zh-tw">XQuery Kick Start</title>
  <author>James McGovern</author>
  <author>Per Bothner</author>
  <author>Kurt Cagle</author>
  <author>James Linn</author>
  <author>Vaidyanathan Nagarajan</author>
  <year>2003</year>
  <price>49.99</price>
</book>

<book category="web">
  <title lang="zh-hk">Learning XML</title>
  <author>Erik T. Ray</author>
  <year>2003</year>
  <price>39.95</price>
</book>

</bookstore> 
```

????????????

```go
package main

func main() {
	content := Open("i.html").Read()
	doc := Tools.XPath(content)
	for _, title := range doc.Find("//title") {
		Lg.Trace("??????lang??????: " + title.GetAttr("lang") + ". ??????title???????????????: " + title.Text())
	}

	book := doc.Find("//bookstore/book[1]")[0]
	Lg.Trace("?????????????????????html: ", book.ChildHTML())

	Lg.Trace("??????book???????????????html: ", book.Html())

	author := doc.Find("//bookstore/book[1]/author[2]")
	Lg.Trace("????????????book??????????????????author:", author)
}
```

??????

```s
02-10 00:54:10   1 [TRAC] (main.go:7) ??????lang??????: en. ??????title???????????????: Everyday Italian
02-10 00:54:10   1 [TRAC] (main.go:7) ??????lang??????: zh-cn. ??????title???????????????: Harry Potter
02-10 00:54:10   1 [TRAC] (main.go:7) ??????lang??????: zh-tw. ??????title???????????????: XQuery Kick Start
02-10 00:54:10   1 [TRAC] (main.go:7) ??????lang??????: zh-hk. ??????title???????????????: Learning XML
02-10 00:54:10   1 [TRAC] (main.go:11) ?????????????????????html:  
                        <title lang="en">Everyday Italian</title>
                        <author>Giada De Laurentiis</author>
                        <year>2005</year>
                        <price>30.00</price>
02-10 00:54:10   1 [TRAC] (main.go:13) ??????book???????????????html:  <book category="cooking">
                        <title lang="en">Everyday Italian</title>
                        <author>Giada De Laurentiis</author>
                        <year>2005</year>
                        <price>30.00</price>
                      </book>
02-10 00:54:10   1 [TRAC] (main.go:16) ????????????book??????????????????author: []
```

## Socket.KCP

??????: 

1. ??????TCP?????????????????????, ??????????????????????????????. 
  1.1 ????????????Connect?????????????????????????????????????????????, ??????????????????????????????, ???????????????????????????, ?????????TCP???SYN???
  1.2 NAT??????????????????UDP???????????????, ???????????????, ???????????????????????????????????????????????????, ????????????????????????20??????????????????????????????????????????????????????3???20?????????????????????????????????????????????????????????3???20???????????????????????????????????????????????????
  1.3 ???????????????????????? ????????????????????????????????????????????????????????????????????????????????????????????????sleep?????????????????????????????????????????????????????????
  1.4 ???????????????????????????????????????????????????????????????????????????????????????timeout???timeout?????????120???????????????120????????????timeout??????????????????goroutine???????????????
2. ?????????
  1.1 ?????????????????????????????????????????? ??????send???????????????goroutine????????????????????????

```go
package main

var key string = "demo key keykeykeykeykeykeykey"
var salt string = "demo salt saltsaltsaltsaltsaltsalt"

var lg *log

func main() {
	args := Argparser("test kcp")
	side := args.Get("", "side", "s", "\"c\" for client, \"s\" for server")
	addr := args.Get("", "addr", "127.0.0.1", "address for listen or connect to")
	port := args.GetInt("", "port", "12345", "port for listen or connect to")
	args.ParseArgs()

  // ?????????
	if side == "c" {
		c := Socket.KCP.Connect(addr, port, key, salt)
		c.Send("1", "2", "3")
    Time.Sleep(1) // ??????1?????? ???????????????????????? ?????????
  // ?????????
	} else if side == "s" {
		k := Socket.KCP.Listen(addr, port, key, salt)
		c := <-k.Accept()
    Print(c.Recv(10)) // ??????[]string{"1", "2", "3"}, ??????10???????????????, ??????nil
	}
}
```

## Socket.TCP.Connect

```go
func main() {
	c := Socket.TCP.Connect("localhost", 8888)
	defer c.Close()
	c.Send("GET / HTTP/1.1\r\n\r\n")
	fmt.Println(c.Recv(1024))
}
```

## Socket.TCP.Listen

```go
func main() {
	l := Socket.TCP.Listen("0.0.0.0", 8899)
	defer l.Close()

	for c := range l.Accept() {
		fmt.Println(c.Recv(1024))
		c.Send("HTTP/1.1 200 OK\r\n\r\n")
		c.Close()
	}
}
```

## Socket.SSL.Listen

```go
func main() {
	crt := `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----`
	key := `-----BEGIN PRIVATE KEY-----
-----END PRIVATE KEY-----`
	cacrt := `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----`
	if os.Args[1] == "s" {
		Lg.Trace("SSL Server started.")
		sl := Socket.SSL.Listen("0.0.0.0", 7789, key, crt)
		for sc := range sl.Accept() {
			go func(sc *TcpServerSideConn) {
				defer sc.Close()
				try(func() {
					for {
						sc.Send(sc.Recv(1024))
					}
				}).Except(func(e error) {
					Lg.Trace("Error:", e)
				})
			}(sc)
		}
	} else if os.Args[1] == "c" {
		Lg.Trace("SSL Client started")
		sc := Socket.SSL.Connect("127.0.0.1", 7789, SSLCfg{Domain: "k.example.it", AdditionRootCA: []string{cacrt}})
		for {
			sc.Send(Time.Strftime("%Y-%m-%d %H:%M:%S", Time.Now()))
			fmt.Println(sc.Recv(1024))
			Time.Sleep(1)
		}
	}
}
```

## Socket.SSL.ServerWrapper

```go
func main() {
	crt := `-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----`
	key := `-----BEGIN PRIVATE KEY-----
...
-----END PRIVATE KEY-----`
	cacrt := `-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----`

	if os.Args[1] == "s" {
		tl := Socket.TCP.Listen("0.0.0.0", 65432)
		tc := <-tl.Accept()

		buf := tc.Recv(10)
		Lg.Trace("Receive from TCP:", buf)

		sc := Socket.SSL.ServerWrapper(tc.Conn, key, crt)

		buf = sc.Recv(10)
		Lg.Trace("Receive from SSL:", buf)
	} else {
		tc := Socket.TCP.Connect("127.0.0.1", 65432)
		tc.Send("Hello TCP!")

		sc := Socket.SSL.ClientWrapper(
			tc.Conn,
			SSLCfg{
				WithoutSystemRootCA: true,
				AdditionRootCA:      []string{cacrt},
				Domain:              "k.example.it",
			})

		sc.Send("Hello SSL!")

		Time.Sleep(1)
	}
}
```

## Socket.UDP

Connect

```go
func main() {
	c := Socket.UDP.Connect("localhost", 8899)
	defer c.Close()
	c.Send("Hello World!")
	fmt.Println(c.Recv(1024))
}
```

Listen

```go
func main() {
	c := Socket.UDP.Listen("0.0.0.0", 8899)
	defer c.Close()
	data, addr := c.Recvfrom(1024)
	fmt.Println(data)
	c.Sendto("You are welcome!", addr)
}
```

## Argparser

* -b???????????????
* -c??????????????????, ??????????????????, ??????????????????????????????????????????, ??????????????????????????????????????????, ??????????????????????????????????????????app, ?????????app.ini??????
* -h????????????, ??????????????????????????????
* ??????????????????????????????, ????????????????????????????????????
* ?????????????????????, ?????????????????????, ??????????????????, ?????????????????????, ???????????????, ???????????????????????????

```go
type arg  {
	InCluster      bool
	ConfigFile     string
	Namespace      string
	TelegramAPIKey string
	TelegramChatID int64
}

func main() {
	args := new(arg)
	a := Argparser("kubernetes???pod????????????????????????")
	args.InCluster = a.GetBool("", "InCluster", "false", "?????????????????????, ????????????????????????????????????config??????")
	args.ConfigFile = a.Get("", "ConfigFile", "", "????????????????????????, ????????????????????????")
	args.Namespace = a.Get("", "Namespace", "", "?????????????????????namespace, ????????????, ????????????, ???????????????")
	args.TelegramAPIKey = a.Get("", "TelegramAPIKey", "", "telegram bot???api key")
	args.TelegramChatID = a.GetInt64("", "TelegramChatID", "0", "telegram bot????????????????????????group??????channel???id")
	a.ParseArgs()
}
```

## Json.XPath

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	s := `{
		"name": "John",
		"age"      : 26,
		"address"  : {
		  "streetAddress": "naist street",
		  "city"         : "Nara",
		  "postalCode"   : "630-0192"
		},
		"phoneNumbers": [
		  {
			"type"  : "iPhone",
			"number": "0123-4567-8888"
		  },
		  {
			"type"  : "home",
			"number": "0123-4567-8910"
		  }
    	],
		"brothers": [
			"john",
			"jack",
		]
    	"nullvalue": null
	}`
	x := Json.XPath(s)
	name := x.First("//name").Text()
	fmt.Printf("Name: %s\n", name)

	var a []string
	for _, n := range x.Find("//phoneNumbers/*/number") {
		a = append(a, n.Text())
	}
	fmt.Printf("All phone number: %s\n", strings.Join(a, ","))

	if n := x.First("//address/streetAddress"); n != nil {
		fmt.Printf("address: %s\n", n.Text())
	}

	for _, b := range j.Find("//brothers/*") {
		Print(b.Text())
	}

  fmt.Println("First phone number:", x.First("//phoneNumbers[1]/*/number").Text()) // ????????????*, ??????????????????????????????????????????element???xml??????, ?????????
  
  Lg.Debug(x.First("//nullvalue").Text()) // ""
}
```

??????

```
Name: John
All phone number: 0123-4567-8888,0123-4567-8910
address: naist street
First phone number: 0123-4567-8888
```

## Tools.GodaddyDNS

```go
func main() {
	gd := Tools.GodaddyDNS("333", "222")
	Print(gd.List())                                  // ??????????????????
	dm := gd.Domain("yletx.com")                      // ?????????????????????
	dm.Add("googledns", "A", "8.8.8.8")               // ???
	dm.Add("googledns_cname", "CNAME", "twitter.com") // ???
	dm.Add("googledns_txt", "TXT", "by twitter")      // ???
	dm.Delete("googledns")                            // ????????????????????????????????????????????????????????????????????????????????????????????????
	dm.Modify("googledns", "A", "1.1.1.1")            // ???
	Print(dm.List())                                  // ???
}
```

## Tools.CloudflareDNS

????????????, ???????????????godaddy??????????????????, ???????????????????????????, ??????????????????cdn. ????????????????????????, ????????????godaddy??????. 

```go
package main

func main() {
	cf := Tools.CloudflareDNS("ip5lwomzy87ohjuoacfzvqup591ipsqi", "example@gmail.com")

	// ??????cloudflare????????????????????????
	cf.Add("example.com")

	Lg.Trace("??????????????????")
	for _, i := range cf.List() {
		if i.Status == "active" {
			Print(i)
			break
		}
	}

	dm := cf.Domain("example.com")

	Lg.Trace("???????????????????????????")
	for _, i := range dm.List() {
		Print(i)
	}

	Lg.Trace("??????dns??????")
	dm.Add("@", "A", "8.8.8.8")
	dm.Add("arecord", "a", "7.7.7.7")
	dm.Add("cnamerecord", "cname", "google.com")
	dm.Add("txtrecord", "txt", "this is a text")

	Lg.Trace("????????????a??????")
	dm.Delete("", "a", "")

	Lg.Trace("????????????a??????")
	dm.Add("@", "A", "8.8.8.8")
	Time.Sleep(5) // ????????????, ????????????????????????, ????????????????????????, ?????????5??????
	dm.Add("@", "A", "6.6.6.6")
	dm.Add("arecord", "a", "7.7.7.7")
	// ????????????a??????
	dm.Modify("@", "a", "8.8.8.8", "a", "5.5.5.5")

	Lg.Trace("????????????cdn")
	dm.SetProxied("@", true)

	Lg.Trace("??????????????????")
	dm.Delete("", "", "")
}
```

## Open

```go
func main() {
	for line := range Open("/etc/passwd").Readlines() { // ????????????chan, for????????????
		fmt.Println(line) // ?????????????????????close???????????????????????????break
	}

	fd := Open("/etc/passwd") // ?????????????????????r
	defer fd.Close()
	fmt.Println(fd.Readline()) // ??????????????? ??????????????????close
	fmt.Println(fd.Readline()) // ?????????????????? ??????????????????close
	fmt.Println(fd.Read(10))   // ??????10???????????? ??????????????????close
	fmt.Println(fd.Read())     // ????????????, ?????????close

	fd = Open("text.txt", "w") // ????????????????????????
	defer fd.Close()
	fd.Write("this is a test text")
	fd.Close()

	fd = Open("text.txt", "rb") // ??????????????????read????????????????????????
	defer fd.Close()
	fmt.Println(fd.Read(5))

	fd = Open("text.txt", "a") // ??????????????????????????????
	defer fd.Close()
	fd.Write(" append text ??????")
	fd.close()
}
```

# Try

```go
func main() {
	Try(func() {
		Bool("abc")
	}).Except(func(err error) {
		fmt.Println(err)
	})
}
```

# Note

`Func.Sniffer`???`Func.ReadPcapFile`????????????????????????`tag`: `pcap`, ???????????????TCP???UDP?????????.

```bash
$ go build -tags pcap .
```















