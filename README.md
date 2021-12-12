Tools

* func Lock() *lockStruct
    * func (*lockStruct) Acquire() 
    * func (*lockStruct) Release() 
* AliDNS: func(string, string) *golanglibs.alidnsStruct
    * func (m *alidnsStruct) Total() (TotalCount int64) 
    * func (m *alidnsStruct) List(PageSize int64, PageNumber int64) (res []alidnsDomainInfoStruct) 
    * func (m *alidnsStruct) Domain(domainName string) *alidnsDomainStruct
        * func (m *alidnsDomainStruct) List() (res []alidnsRecord)
        * func (m *alidnsDomainStruct) Add(recordName string, recordType string, recordValue string) (id string)
        * func (m *alidnsDomainStruct) Delete(name string, dtype string, value string) 
        * func (m *alidnsDomainStruct) modify(recordName string, srcRecordType string, srcRecordValue string, dstRecordName string, dstRecordType string, dstRecordValue string)
* Chart
    * LineChartWithTimestampAndNumber: func([]int64, []float64, string, string, string, string) {...},
    * LineChartWithNumberAndNumber:    func([]float64, []float64, string, string, string, string) {...},
    * BarChartWithNameAndNumber:       func([]string, []float64, string, string, string) {...},
    * PieChartWithNameAndNumber:       func([]string, []float64, string, string) {...},
CloudflareDNS: func(string, string) *golanglibs.cloudflareStruct 
    func (m *cloudflareStruct) Add(domain string) cloudflare.Zone {
    func (m *cloudflareStruct) List() (res []cloudflareDomainInfoStruct) {
    func (m *cloudflareStruct) Domain(domainName string) *cloudflareDomainStruct {
        func (m *cloudflareDomainStruct) List() (res []cloudflareRecord) {
        func (m *cloudflareDomainStruct) Delete(name string) {
        func (m *cloudflareDomainStruct) Add(recordName string, recordType string, recordValue string, proxied ...bool) *cloudflare.DNSRecordResponse {
        func (m *cloudflareDomainStruct) SetProxied(subdomain string, proxied bool) {
        func (m *cloudflareDomainStruct) Update(recordName string, recordType string, recordValue string, proxied ...bool) {
Compress:
    LzmaCompressString:   func(string) string {...},
    LzmaDecompressString: func(string) string {...},
    ZlibCompressString:   func(string) string {...},
    ZlibDecompressString: func(string) string {...},
Crontab:      func() *golanglibs.crontabStruct {...},
    func (m *crontabStruct) Add(schedule string, fn interface{}, args ...interface{})
    func (m *crontabStruct) Destory()
GodaddyDNS:   func(string, string) *golanglibs.godaddyStruct {...},
    func (m *godaddyStruct) List() (res []godaddyDomainInfoStruct) {
    func (m *godaddyStruct) Domain(domainName string) *godaddyDomainStruct {
        func (m *godaddyDomainStruct) List() (res []godaddyRecord) {
        func (m *godaddyDomainStruct) Delete(name string, dtype string, value string) {
        func (m *godaddyDomainStruct) Modify(recordName string, srcRecordType string, srcRecordValue string, dstRecordType string, dstRecordValue string) {
        func (m *godaddyDomainStruct) Add(recordName string, recordType string, recordValue string) {
Ini:          func(...string) *golanglibs.iniStruct {...},
    func (m *iniStruct) Get(SectionKeyDefaultComment ...string) (res string) {
    func (m *iniStruct) GetInt(key ...string) int {
    func (m *iniStruct) GetInt64(key ...string) int64 {
    func (m *iniStruct) getFloat64(key ...string) float64 {
    func (m *iniStruct) Set(SectionKeyValueComment ...string) {
    func (m *iniStruct) Save(fpath ...string) (exist bool) {
JavascriptVM: func() *golanglibs.javascriptVMStruct {...},
    func (m *javascriptVMStruct) Run(javascript string) *javascriptVMStruct {
    func (m *javascriptVMStruct) Get(variableName string) string {
    func (m *javascriptVMStruct) Set(variableName string, variableValue interface{}) {
    func (m *javascriptVMStruct) Isdefined(variableName string) bool {
Matrix:       func(string) *golanglibs.matrixStruct {...},
    func (c *matrixStruct) Login(username string, password string) string {
    func (c *matrixStruct) SetToken(userID string, token string) *matrixStruct {
    func (c *matrixStruct) SetRoomID(roomID string) *matrixStruct {
    func (c *matrixStruct) Send(text string) {
Nats:         func(string) *golanglibs.natsStruct {...},
    func (m *natsStruct) Subject(subject string) *subjectNatsStruct {
        func (m *subjectNatsStruct) Publish(message string) {
        func (m *subjectNatsStruct) Subscribe() chan string {
        func (m *subjectNatsStruct) Flush() {
Totp:         func(string) *golanglibs.totpStruct {...},
    func (m *totpStruct) Validate(pass string) bool {
    func (m *totpStruct) Password() string {
Pexpect:      func(string) *golanglibs.pexpectStruct {...},
    func (m *pexpectStruct) Sendline(msg string) {
    func (m *pexpectStruct) Close() {
ProgressBar:  func(string, int64, ...bool) *golanglibs.progressBarStruct {...},
    func (m *progressBarStruct) Add(num int64) {
    func (m *progressBarStruct) Set(num int64) {
    func (m *progressBarStruct) SetTotal(total int64) {
    func (m *progressBarStruct) Clear() {
Prometheus:   func(string) *golanglibs.prometheusStruct {...},
    func (m *prometheusStruct) Query(query string, time ...float64) (res []prometheusResultStruct) 
MySQL:        func(string, int, string, string, string, ...golanglibs.DatabaseConfig) *golanglibs.databaseStruct {...},
SQLite:       func(string) *golanglibs.databaseStruct {...},
    func (m *databaseStruct) Query(sql string, args ...interface{}) []gorose.Data {
    func (m *databaseStruct) Close() {
    func (m *databaseStruct) Execute(sql string) int64 {
    func (m *databaseStruct) RenameTable(oldTableName string, newTableNname string) {
    func (m *databaseStruct) tables() (res []string) {
    func (m *databaseStruct) createTable(tableName string, engineName ...string) *databaseOrmStruct {
    func (m *databaseStruct) Table(tbname string) *databaseOrmStruct {
        func (m *databaseOrmStruct) Fields(items ...string) *databaseOrmStruct {
        func (m *databaseOrmStruct) Where(key string, operator string, value interface{}) *databaseOrmStruct {
        func (m *databaseOrmStruct) WhereIn(key string, value []interface{}) *databaseOrmStruct {
        func (m *databaseOrmStruct) WhereNotIn(key string, value []interface{}) *databaseOrmStruct {
        func (m *databaseOrmStruct) WhereNull(columnName string) *databaseOrmStruct {
        func (m *databaseOrmStruct) WhereNotNull(columnName string) *databaseOrmStruct {
        func (m *databaseOrmStruct) OrWhere(key string, operator string, value interface{}) *databaseOrmStruct {
        func (m *databaseOrmStruct) OrWhereIn(key string, value []interface{}) *databaseOrmStruct {
        func (m *databaseOrmStruct) Orderby(item ...string) *databaseOrmStruct {
        func (m *databaseOrmStruct) Limit(number int) *databaseOrmStruct {
        func (m *databaseOrmStruct) Get() (res []gorose.Data) {
        func (m *databaseOrmStruct) Paginate(pagesize int, page int) []gorose.Data {
        func (m *databaseOrmStruct) First() (res gorose.Data) {
        func (m *databaseOrmStruct) find(id int) gorose.Data {
        func (m *databaseOrmStruct) Count() (res int64) {
        func (m *databaseOrmStruct) exists() (res bool) {
        func (m *databaseOrmStruct) chunk(limit int, callback func([]gorose.Data) error) {
        func (m *databaseOrmStruct) buildSQL() (string, []interface{}) {
        func (m *databaseOrmStruct) data(data interface{}) *databaseOrmStruct {
        func (m *databaseOrmStruct) offset(offset int) *databaseOrmStruct {
        func (m *databaseOrmStruct) insertGetID() (num int64) {
        func (m *databaseOrmStruct) Insert() (num int64) {
        func (m *databaseOrmStruct) Update(data ...interface{}) (num int64) {
        func (m *databaseOrmStruct) Delete() (num int64) {
        func (m *databaseOrmStruct) DropTable() int64 {
        func (m *databaseOrmStruct) TruncateTable() (status int64) {
        func (m *databaseOrmStruct) AddColumn(columnName string, columnType string, defaultValue ...string) *databaseOrmStruct {
        func (m *databaseOrmStruct) DropColumn(columnName string) *databaseOrmStruct {
        func (m *databaseOrmStruct) AddIndex(columnName ...string) *databaseOrmStruct {
        func (m *databaseOrmStruct) IndexExists(columnName ...string) (exists bool) {
        func (m *databaseOrmStruct) DropIndex(columnName ...string) *databaseOrmStruct {
        func (m *databaseOrmStruct) Columns() (res map[string]string) {
RabbitMQ:     func(string, string) *golanglibs.rabbitConnectionStruct {...},
    func (m *rabbitConnectionStruct) Send(data map[string]string) {
    func (m *rabbitConnectionStruct) Recv() chan map[string]string {
RateLimit:    func(int) *golanglibs.rateLimitStruct {...},
    func (m *rateLimitStruct) Take() {
Redis:        func(string, int, string, int, ...golanglibs.redisConfig) *golanglibs.RedisStruct {...},
    func (m *RedisStruct) Ping() string {
    func (m *RedisStruct) Del(key string) {
    func (m *RedisStruct) Set(key string, value string, ttl ...interface{}) {
    func (m *RedisStruct) Get(key string) *string {
    func (m *RedisStruct) GetLock(key string, timeoutsec int) *RedisLockStruct {
        func (m *RedisLockStruct) acquire() {
        func (m *RedisLockStruct) Release() {
Selenium:     func(string) *golanglibs.seleniumStruct {...},
    func (c *seleniumStruct) Close() {
    func (c *seleniumStruct) Cookie() (co string) {
    func (c *seleniumStruct) Url() string {
    func (c *seleniumStruct) ScrollRight(pixel int) {
    func (c *seleniumStruct) ScrollLeft(pixel int) {
    func (c *seleniumStruct) ScrollUp(pixel int) {
    func (c *seleniumStruct) ScrollDown(pixel int) {
    func (c *seleniumStruct) ResizeWindow(width int, height int) *seleniumStruct {
    func (c *seleniumStruct) Find(xpath string, nowait ...bool) *seleniumElementStruct {
        func (c *seleniumElementStruct) Clear() *seleniumElementStruct {
        func (c *seleniumElementStruct) Click() *seleniumElementStruct {
        func (c *seleniumElementStruct) Text() string {
        func (c *seleniumElementStruct) Input(s string) *seleniumElementStruct {
        func (c *seleniumElementStruct) Submit() *seleniumElementStruct {
        func (c *seleniumElementStruct) PressEnter() *seleniumElementStruct {
SSH:          func(string, string, string, int) *golanglibs.sshStruct {...},
    func (m *sshStruct) Close() {
    func (m *sshStruct) Exec(cmd string) (output string, status int) {
    func (m *sshStruct) PushFile(local string, remote string) {
    func (m *sshStruct) PullFile(remote string, local string) {
StatikOpen:   func(string) *golanglibs.statikFileStruct {...},
    func (m *statikFileStruct) Readlines() chan string {
    func (m *statikFileStruct) Readline() string {
    func (m *statikFileStruct) Close() {
    func (m *statikFileStruct) Read(num ...int) string {
    func (m *statikFileStruct) Seek(num int64) {
Table:        func(...string) *golanglibs.tableStruct {...},
    func (m *tableStruct) SetMaxCellWidth(width ...int) {
    func (m *tableStruct) AddRow(row ...interface{}) {
    func (m *tableStruct) Render() string {
TelegramBot:  func(string) *golanglibs.telegramBotStruct {...},
    func (m *telegramBotStruct) SetChatID(chatid int64) *telegramBotStruct {
    func (m *telegramBotStruct) SendFile(path string) tgbotapi.Message {
    func (m *telegramBotStruct) SendImage(path string) tgbotapi.Message {
    func (m *telegramBotStruct) Send(text string, cfg ...tgMsgConfig) tgbotapi.Message {
Telegraph:    func(string) *golanglibs.telegraphStruct {...},
    func (m *telegraphStruct) Post(title string, content string) *telegraphPageInfo
URL:          func(string) *golanglibs.urlStruct {...},
    func (u *urlStruct) Parse() *urlComponents {
    func (u *urlStruct) Encode() string {
    func (u *urlStruct) Decode() string {
TTLCache:     func(interface {}) *golanglibs.ttlCacheStruct {...},
    func (m *ttlCacheStruct) Set(key string, value string) {
    func (m *ttlCacheStruct) Get(key string) string {
    func (m *ttlCacheStruct) Exists(key string) bool {
    func (m *ttlCacheStruct) Count() int {
VNC:          func(string, ...golanglibs.VNCCfg) *golanglibs.vncStruct {...},
    func (m *vncStruct) Close() {
    func (m *vncStruct) Move(x, y int) *vncStruct {
    func (m *vncStruct) Click() *vncStruct {
    func (m *vncStruct) RightClick() *vncStruct {
    func (m *vncStruct) Input(s string) *vncStruct {
    func (m *vncStruct) Key() *vncKeyStruct {
        func (m *vncKeyStruct) Enter() *vncKeyStruct {
        func (m *vncKeyStruct) Ctrl_a() *vncKeyStruct {
        func (m *vncKeyStruct) Ctrl_c() *vncKeyStruct {
        func (m *vncKeyStruct) Ctrl_v() *vncKeyStruct {
        func (m *vncKeyStruct) Ctrl_z() *vncKeyStruct {
        func (m *vncKeyStruct) Ctrl_x() *vncKeyStruct {
        func (m *vncKeyStruct) Ctrl_f() *vncKeyStruct {
        func (m *vncKeyStruct) Ctrl_d() *vncKeyStruct {
        func (m *vncKeyStruct) Ctrl_s() *vncKeyStruct {
        func (m *vncKeyStruct) Ctrl_r() *vncKeyStruct {
        func (m *vncKeyStruct) Ctrl_e() *vncKeyStruct {
        func (m *vncKeyStruct) delete() *vncKeyStruct {
        func (m *vncKeyStruct) tab() *vncKeyStruct {
WebSocket:    func(string) *golanglibs.websocketStruct {...},
    func (c *websocketStruct) Send(text string) {
    func (c *websocketStruct) Recv(timeout ...int) string {
    func (c *websocketStruct) Close() {
Xlsx:         func(string) *golanglibs.xlsxStruct {...},
    func (c *xlsxStruct) Save() {
    func (c *xlsxStruct) Close() {
    func (c *xlsxStruct) GetSheet(name string) *xlsxSheetStruct {
        func (c *xlsxSheetStruct) Get(coordinate string) string {
        func (c *xlsxSheetStruct) Set(coordinate string, value string) *xlsxSheetStruct {
XPath:        func(string) *golanglibs.xpathStruct {...},
    func (m *xpathStruct) First(expr string) (res *xpathStruct) {
    func (m *xpathStruct) Find(expr string) (res []*xpathStruct) {
    func (m *xpathStruct) Text() string {
    func (m *xpathStruct) GetAttr(attr string) string {
    func (m *xpathStruct) Html() string {
JsonXPath:    func(string) *golanglibs.xpathJsonStruct {...},
    func (m *xpathJsonStruct) Exists(expr string) bool {
    func (m *xpathJsonStruct) First(expr string) (res *xpathJsonStruct) {
    func (m *xpathJsonStruct) Find(expr string) (res []*xpathJsonStruct) {
    func (m *xpathJsonStruct) Text() string {

Random:
	Int    func(min, max int64) int64
	Choice func(array interface{}) interface{}
	String func(length int, charset ...string) string


Re:
	FindAll func(pattern string, text string, multiline ...bool) [][]string
	Replace func(pattern string, newstring string, text string) string

Socket:

  KCP: golanglibs.kcpStruct{
    Listen:     func(string, int, string, string) *golanglibs.kcpServerSideListener {...},
    Connect:    func(string, int, string, string) *golanglibs.kcpClientSideConn {...},
    RawListen:  func(string, int, string, string) *golanglibs.kcpRawServerSideListener {...},
    RawConnect: func(string, int, string, string) *kcp.UDPSession {...},
  },
  Smux: golanglibs.smuxStruct{
    ServerWrapper: func(io.ReadWriteCloser, ...golanglibs.SmuxConfig) *golanglibs.smuxServerSideListener {...},
    ClientWrapper: func(io.ReadWriteCloser, ...golanglibs.SmuxConfig) *golanglibs.smuxClientSideSession {...},
  },
  SSL: golanglibs.sslStruct{
    Listen:        func(string, int, string, string) *golanglibs.tcpServerSideListener {...},
    ServerWrapper: func(net.Conn, string, string) *golanglibs.tcpServerSideConn {...},
    Connect:       func(string, int, ...golanglibs.sslCfg) *golanglibs.sslClientSideConn {...},
    ClientWrapper: func(net.Conn, ...golanglibs.sslCfg) *golanglibs.sslClientSideConn {...},
  },
  TCP: golanglibs.tcpStruct{
    Listen:  func(string, int) *golanglibs.tcpServerSideListener {...},
    Connect: func(string, int, ...int) *golanglibs.tcpClientSideConn {...},
  },
  UDP: golanglibs.udpStruct{
    Listen:  func(string, int) golanglibs.udpServerSideConn {...},
    Connect: func(string, int) golanglibs.udpClientSideConn {...},
  },

String:

Time:
  Now:            func() float64 {...},
  TimeDuration:   func(interface {}) time.Duration {...},
  FormatDuration: func(int64) string {...},
  Sleep:          func(interface {}) {...},
  Strptime:       func(string, string) int64 {...},
  Strftime:       func(string, interface {}) string {...},

Argparser(description string) *argparseIniStruct
    func (m *argparseIniStruct) Get(section, key, defaultValue, comment string) (res string) {
    func (m *argparseIniStruct) GetInt(section, key, defaultValue, comment string) int {
    func (m *argparseIniStruct) GetInt64(section, key, defaultValue, comment string) int64 {
    func (m *argparseIniStruct) GetFloat64(section, key, defaultValue, comment string) float64 {
    func (m *argparseIniStruct) GetBool(section, key, defaultValue, comment string) bool {
    func (m *argparseIniStruct) Save(fpath ...string) (exist bool) {
    func (m *argparseIniStruct) GetHelpString() (h string) {
    func (m *argparseIniStruct) ParseArgs() *argparseIniStruct {

Base64:
	Encode func(str string) string
	Decode func(str string) string



Binary:
	Map2bin func(m map[string]string) string
	Bin2map func(s string) (res map[string]string)

Cmd:
	GetOutput                func(command string, timeoutSecond ...interface{}) string
	GetStatusOutput          func(command string, timeoutSecond ...interface{}) (int, string)
	GetOutputWithShell       func(command string, timeoutSecond ...interface{}) string
	GetStatusOutputWithShell func(command string, timeoutSecond ...interface{}) (int, string)
	Tail                     func(command string) chan string
	Exists                   func(cmd string) bool


Crypto
	Xor func(data, key string) string
	Aes func(key string) *aesStruct

File(filePath string) *fileStruct
    func (f *fileStruct) Time() *fileTimeStruct {
    func (f *fileStruct) Stat() os.FileInfo {
    func (f *fileStruct) Size() int64 {
    func (f *fileStruct) Touch() {
    func (f *fileStruct) Chmod(mode os.FileMode) bool {
    func (f *fileStruct) Chown(uid, gid int) bool {
    func (f *fileStruct) Mtime() int64 {
    func (f *fileStruct) Unlink() {
    func (f *fileStruct) Copy(dest string) {
    func (f *fileStruct) Move(newPosition string) {

Open(args ...string) *fileIOStruct
    func (m *fileIOStruct) Readlines() chan string {
    func (m *fileIOStruct) Readline() string {
    func (m *fileIOStruct) Close() {
    func (m *fileIOStruct) Write(str interface{}) *fileIOStruct {
    func (m *fileIOStruct) Read(num ...int) string {
    func (m *fileIOStruct) Seek(num int64) {


Funcs
	Nslookup               func(name string, querytype string, dnsService ...string) [][]string
	FakeName               func() string
	FileType               func(fpath string) string
	Inotify                func(path string) chan *fsnotifyFileEventStruct
	IPLocation             func(ip string, dbpath ...string) *ipLocationInfo
	HightLightHTMLForCode  func(code string, codeType ...string) (html string)
	Markdown2html          func(md string) string
	CPUUsagePerProgress    func() (res map[int64]progressCPUUsageStruct)
	ResizeImg              func(srcPath string, dstPath string, width int, height ...int)
	GetRSS                 func(url string, config ...rssConfig) *gofeed.Feed
	GbkToUtf8              func(s string) string
	Utf8ToGbk              func(s string) string
	GetSnowflakeID         func(nodeNumber ...int) int64
	GetRemoteServerSSLCert func(host string, port ...int) []*x509.Certificate
	Tailf                  func(path string, startFromEndOfFile ...bool) chan *tail.Line
	BaiduTranslateAnyToZH  func(text string) string
	ParseUserAgent         func(UserAgent string) ua.UserAgent
	Wget                   func(url string, cfg ...WgetCfg) (filename string)
	Whois                  func(s string, servers ...string) string
	IpInNet                func(ip string, Net string, mask ...string) bool
	Int2ip                 func(ipnr int64) string
	Ip2int                 func(ipnr string) int64
    Zh2PinYin              func(zh string) (ress []string)


Hash
	Md5sum   func(str string) string
	Md5File  func(path string) string
	Sha1sum  func(str string) string
	Sha1File func(path string) string

Html
	Encode func(str string) string
	Decode func(str string) string

Http
	Head     func(uri string, args ...interface{}) httpResp
	PostFile func(uri string, filePath string, args ...interface{}) httpResp
	PostRaw  func(uri string, body string, args ...interface{}) httpResp
	PostJSON func(uri string, json interface{}, args ...interface{}) httpResp
	Post     func(uri string, args ...interface{}) httpResp
	Get      func(uri string, args ...interface{}) httpResp
	PutJSON  func(uri string, json interface{}, args ...interface{}) httpResp
	Put      func(uri string, args ...interface{}) httpResp
	PutRaw   func(uri string, body string, args ...interface{}) httpResp

Json
	Dumps     func(v interface{}, pretty ...bool) string
	Loads     func(str string) map[string]interface{}
	Valid     func(j string) bool
	Yaml2json func(y string) string
	Json2yaml func(j string) string
	Format    func(js string) string

Math
	Abs          func(number float64) float64
	Sum          func(array interface{}) (sumresult float64)
	Average      func(array interface{}) (avgresult float64)
	DecimalToAny func(num, n int64) string
	AnyToDecimal func(num string, n int64) int64

Os
	Mkdir           func(filename string)
	Getcwd          func() string
	Exit            func(status int)
	Touch           func(filePath string)
	Chmod           func(filePath string, mode os.FileMode) bool
	Chown           func(filePath string, uid, gid int) bool
	Copy            func(filePath, dest string)
	Rename          func(filePath, newPosition string)
	Move            func(filePath, newPosition string)
	Path            pathStruct
	System          func(command string, timeoutSecond ...interface{}) int
	SystemWithShell func(command string, timeoutSecond ...interface{}) int
	Hostname        func() string
	Envexists       func(varname string) bool
	Getenv          func(varname string) string
	Walk            func(path string) chan string
	Listdir         func(path string) (res []string)
	SelfDir         func() string
	TempFilePath    func() string
	TempDirPath     func() string
	Getuid          func() int
	ProgressAlive   func(pid int) bool
	GoroutineID     func() int64
	Unlink          func(filename string)

