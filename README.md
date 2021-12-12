* Tools
    * func Lock() \*lock
        * func (*lock) Acquire() 
        * func (*lock) Release() 
    * func AliDNS(string, string) \*alidns
        * func (m \*alidns) Total() (TotalCount int64) 
        * func (m \*alidns) List(PageSize int64, PageNumber int64) (res []alidnsDomainInfo) 
        * func (m \*alidns) Domain(domainName string) \*alidnsDomain
            * func (m \*alidnsDomain) List() (res []alidnsRecord)
            * func (m \*alidnsDomain) Add(recordName string, recordType string, recordValue string) (id string)
            * func (m \*alidnsDomain) Delete(name string, dtype string, value string) 
            * func (m \*alidnsDomain) modify(recordName string, srcRecordType string, srcRecordValue string, dstRecordName string, dstRecordType string, dstRecordValue string)
    * Chart
        * func LineChartWithTimestampAndNumber([]int64, []float64, string, string, string, string)
        * func LineChartWithNumberAndNumber([]float64, []float64, string, string, string, string)
        * func BarChartWithNameAndNumber([]string, []float64, string, string, string)
        * func PieChartWithNameAndNumber([]string, []float64, string, string)
    * func CloudflareDNS(string, string) \*cloudflare 
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
    * func Crontab() \*crontab
        * func (m \*crontab) Add(schedule string, fn interface{}, args ...interface{})
        * func (m \*crontab) Destory()
    * func GodaddyDNS(string, string) \*godaddy
        * func (m \*godaddy) List() (res []godaddyDomainInfo)
        * func (m \*godaddy) Domain(domainName string) \*godaddyDomain
            * func (m \*godaddyDomain) List() (res []godaddyRecord)
            * func (m \*godaddyDomain) Delete(name string, dtype string, value string)
            * func (m \*godaddyDomain) Modify(recordName string, srcRecordType string, srcRecordValue string, dstRecordType string, dstRecordValue string)
            * func (m \*godaddyDomain) Add(recordName string, recordType string, recordValue string)
    * func Ini(...string) \*ini
        * func (m \*ini) Get(SectionKeyDefaultComment ...string) (res string)
        * func (m \*ini) GetInt(key ...string) int
        * func (m \*ini) GetInt64(key ...string) int64
        * func (m \*ini) getFloat64(key ...string) float64
        * func (m \*ini) Set(SectionKeyValueComment ...string)
        * func (m \*ini) Save(fpath ...string) (exist bool)
    * func JavascriptVM() \*javascriptVM
        * func (m \*javascriptVM) Run(javascript string) \*javascriptVM
        * func (m \*javascriptVM) Get(variableName string) string
        * func (m \*javascriptVM) Set(variableName string, variableValue interface{})
        * func (m \*javascriptVM) Isdefined(variableName string) bool
    * func Matrix(string) \*matrix
        * func (c \*matrix) Login(username string, password string) string
        * func (c \*matrix) SetToken(userID string, token string) \*matrix
        * func (c \*matrix) SetRoomID(roomID string) \*matrix
        * func (c \*matrix) Send(text string)
    * func Nats(string) \*nats
        * func (m \*nats) Subject(subject string) \*subjectNats
            * func (m \*subjectNats) Publish(message string)
            * func (m \*subjectNats) Subscribe() chan string
            * func (m \*subjectNats) Flush()
    * func Totp(string) \*totp
        * func (m \*totp) Validate(pass string) bool
        * func (m \*totp) Password() string
    * func Pexpect(string) \*pexpect
        * func (m \*pexpect) Sendline(msg string)
        * func (m \*pexpect) Close()
    * func ProgressBar(string, int64, ...bool) \*progressBar
        * func (m \*progressBar) Add(num int64)
        * func (m \*progressBar) Set(num int64)
        * func (m \*progressBar) SetTotal(total int64)
        * func (m \*progressBar) Clear()
    * func Prometheus(string) \*prometheus
        * func (m \*prometheus) Query(query string, time ...float64) (res []prometheusResult) 
    * func MySQL(string, int, string, string, string, ...DatabaseConfig) \*database
    * func SQLite(string) \*database
        * func (m \*database) Query(sql string, args ...interface{}) []gorose.Data
        * func (m \*database) Close()
        * func (m \*database) Execute(sql string) int64
        * func (m \*database) RenameTable(oldTableName string, newTableNname string)
        * func (m \*database) tables() (res []string)
        * func (m \*database) createTable(tableName string, engineName ...string) \*databaseOrm
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
            * func (m \*databaseOrm) find(id int) gorose.Data
            * func (m \*databaseOrm) Count() (res int64)
            * func (m \*databaseOrm) exists() (res bool)
            * func (m \*databaseOrm) chunk(limit int, callback func([]gorose.Data) error)
            * func (m \*databaseOrm) buildSQL() (string, []interface{})
            * func (m \*databaseOrm) data(data interface{}) \*databaseOrm
            * func (m \*databaseOrm) offset(offset int) \*databaseOrm
            * func (m \*databaseOrm) insertGetID() (num int64)
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
    * func RabbitMQ(string, string) \*rabbitConnection
        * func (m \*rabbitConnection) Send(data map[string]string)
        * func (m \*rabbitConnection) Recv() chan map[string]string
    * func RateLimit(int) \*rateLimit
        * func (m \*rateLimit) Take()
    * func Redis(string, int, string, int, ...redisConfig) \*Redis
        * func (m \*Redis) Ping() string
        * func (m \*Redis) Del(key string)
        * func (m \*Redis) Set(key string, value string, ttl ...interface{})
        * func (m \*Redis) Get(key string) \*string
        * func (m \*Redis) GetLock(key string, timeoutsec int) \*RedisLock
            * func (m \*RedisLock) acquire()
            * func (m \*RedisLock) Release()
    * func Selenium(string) \*selenium
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
    * func SSH(string, string, string, int) \*ssh
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
    * func TelegramBot(string) \*telegramBot
        * func (m \*telegramBot) SetChatID(chatid int64) \*telegramBot
        * func (m \*telegramBot) SendFile(path string) tgbotapi.Message
        * func (m \*telegramBot) SendImage(path string) tgbotapi.Message
        * func (m \*telegramBot) Send(text string, cfg ...tgMsgConfig) tgbotapi.Message
    * func Telegraph(string) \*telegraph
        * func (m \*telegraph) Post(title string, content string) \*telegraphPageInfo
    * func URL(string) \*url
        * func (u \*url) Parse() \*urlComponents
        * func (u \*url) Encode() string
        * func (u \*url) Decode() string
    * func TTLCache(interface {}) \*ttlCache
        * func (m \*ttlCache) Set(key string, value string)
        * func (m \*ttlCache) Get(key string) string
        * func (m \*ttlCache) Exists(key string) bool
        * func (m \*ttlCache) Count() int
    * func VNC(string, ...VNCCfg) \*vnc
        * func (m \*vnc) Close()
        * func (m \*vnc) Move(x, y int) \*vnc
        * func (m \*vnc) Click() \*vnc
        * func (m \*vnc) RightClick() \*vnc
        * func (m \*vnc) Input(s string) \*vnc
        * func (m \*vnc) Key() \*vncKey
            * func (m \*vncKey) Enter() \*vncKey
            * func (m \*vncKey) Ctrl_a() \*vncKey
            * func (m \*vncKey) Ctrl_c() \*vncKey
            * func (m \*vncKey) Ctrl_v() \*vncKey
            * func (m \*vncKey) Ctrl_z() \*vncKey
            * func (m \*vncKey) Ctrl_x() \*vncKey
            * func (m \*vncKey) Ctrl_f() \*vncKey
            * func (m \*vncKey) Ctrl_d() \*vncKey
            * func (m \*vncKey) Ctrl_s() \*vncKey
            * func (m \*vncKey) Ctrl_r() \*vncKey
            * func (m \*vncKey) Ctrl_e() \*vncKey
            * func (m \*vncKey) delete() \*vncKey
            * func (m \*vncKey) tab() \*vncKey
    * func WebSocket(string) \*websocket
        * func (c \*websocket) Send(text string)
        * func (c \*websocket) Recv(timeout ...int) string
        * func (c \*websocket) Close()
    * func Xlsx(string) \*xlsx
        * func (c \*xlsx) Save()
        * func (c \*xlsx) Close()
        * func (c \*xlsx) GetSheet(name string) \*xlsxSheet
            * func (c \*xlsxSheet) Get(coordinate string) string
            * func (c \*xlsxSheet) Set(coordinate string, value string) \*xlsxSheet
    * func XPath(string) \*xpath
        * func (m \*xpath) First(expr string) (res *xpath)
        * func (m \*xpath) Find(expr string) (res []*xpath)
        * func (m \*xpath) Text() string
        * func (m \*xpath) GetAttr(attr string) string
        * func (m \*xpath) Html() string
    * func JsonXPath(string) \*xpathJson
        * func (m \*xpathJson) Exists(expr string) bool
        * func (m \*xpathJson) First(expr string) (res *xpathJson)
        * func (m \*xpathJson) Find(expr string) (res []*xpathJson)
        * func (m \*xpathJson) Text() string
* Random
    * func Int(min, max int64) int64
    * func Choice(array interface{}) interface{}
    * func String(length int, charset ...string) string
* Re
    * func FindAll(pattern string, text string, multiline ...bool) [][]string
    * func Replace(pattern string, newstring string, text string) string
* Socket
    * KCP: kcp{
        * func Listen(string, int, string, string) \*kcpServerSideListener
        * func Connect(string, int, string, string) \*kcpClientSideConn
        * func RawListen(string, int, string, string) \*kcpRawServerSideListener
        * func RawConnect(string, int, string, string) \*kcp.UDPSession
    * Smux: smux{
        * func ServerWrapper(io.ReadWriteCloser, ...SmuxConfig) \*smuxServerSideListener
        * func ClientWrapper(io.ReadWriteCloser, ...SmuxConfig) \*smuxClientSideSession
    * SSL: ssl{
        * func Listen(string, int, string, string) \*tcpServerSideListener
        * func ServerWrapper(net.Conn, string, string) \*tcpServerSideConn
        * func Connect(string, int, ...sslCfg) \*sslClientSideConn
        * func ClientWrapper(net.Conn, ...sslCfg) \*sslClientSideConn
    * TCP: tcp{
        * func Listen(string, int) \*tcpServerSideListener
        * func Connect(string, int, ...int) \*tcpClientSideConn
    * UDP: udp{
        * func Listen(string, int) udpServerSideConn
        * func Connect(string, int) udpClientSideConn
* String
    * func (s \*string) Get() string
    * func (s \*string) Sub(start, end int) \*string
    * func (s \*string) Has(substr string) bool
    * func (s \*string) sub(start, end int) string
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
    * func (s \*string) hasChinese() bool
    * func (s \*string) Filter(charts ...string) \*string
    * func (s \*string) RemoveHtmlTag() \*string
    * func (s \*string) RemoveNonUTF8Character() \*string
    * func (s \*string) DetectLanguage() string
* Time:
    * func Now() float64
    * func TimeDuration(interface {}) time.Duration
    * func FormatDuration(int64) string
    * func Sleep(interface {})
    * func Strptime(string, string) int64
    * func Strftime(string, interface {}) string
* Argparser(description string) \*argparseIni
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
    * Map2bin func(m map[string]string) string
    * Bin2map func(s string) (res map[string]string)
* Cmd
    * func GetOutput(command string, timeoutSecond ...interface{}) string
    * func GetStatusOutput(command string, timeoutSecond ...interface{}) (int, string)
    * func GetOutputWithShell(command string, timeoutSecond ...interface{}) string
    * func GetStatusOutputWithShell(command string, timeoutSecond ...interface{}) (int, string)
    * func Tail(command string) chan string
    * func Exists(cmd string) bool
* Crypto
    * func Xor(data, key string) string
    * func Aes(key string) \*aes
* File(filePath string) \*file
    * func (f \*file) Time() \*fileTime
    * func (f \*file) Stat() os.FileInfo
    * func (f \*file) Size() int64
    * func (f \*file) Touch()
    * func (f \*file) Chmod(mode os.FileMode) bool
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
    * Markdown2html          func(md string) string
    * func CPUUsagePerProgress() (res map[int64]progressCPUUsage)
    * func ResizeImg(srcPath string, dstPath string, width int, height ...int)
    * func GetRSS(url string, config ...rssConfig) \*gofeed.Feed
    * GbkToUtf8              func(s string) string
    * Utf8ToGbk              func(s string) string
    * func GetSnowflakeID(nodeNumber ...int) int64
    * func GetRemoteServerSSLCert(host string, port ...int) []*x509.Certificate
    * func Tailf(path string, startFromEndOfFile ...bool) chan *tail.Line
    * func BaiduTranslateAnyToZH(text string) string
    * func ParseUserAgent(UserAgent string) ua.UserAgent
    * func Wget(url string, cfg ...WgetCfg) (filename string)
    * func Whois(s string, servers ...string) string
    * func IpInNet(ip string, Net string, mask ...string) bool
    * Int2ip                 func(ipnr int64) string
    * Ip2int                 func(ipnr string) int64
    * Zh2PinYin              func(zh string) (ress []string)
* Hash
    * Md5sum   func(str string) string
    * Md5File  func(path string) string
    * Sha1sum  func(str string) string
    * Sha1File func(path string) string
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
* Json
    * func Dumps(v interface{}, pretty ...bool) string
    * func Loads(str string) map[string]interface{}
    * func Valid(j string) bool
    * Yaml2json func(y string) string
    * Json2yaml func(j string) string
    * func Format(js string) string
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
    * func Chmod(filePath string, mode os.FileMode) bool
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

