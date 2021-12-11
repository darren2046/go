package golanglibs

import "sync"

type toolsStruct struct {
	Lock          func() *lockStruct
	AliDNS        func(accessKeyID string, accessKeySecret string) *alidnsStruct
	Chart         chartStruct
	CloudflareDNS func(key string, email string) *cloudflareStruct
	Compress      compressStruct
	Crontab       func() *crontabStruct
	GodaddyDNS    func(key string, secret string) *godaddyStruct
	Ini           func(fpath ...string) *iniStruct
	JavascriptVM  func() *javascriptVMStruct
}

var Tools toolsStruct

func init() {
	Tools = toolsStruct{
		Lock:          getLock,
		AliDNS:        getAlidns,
		Chart:         chartstruct,
		CloudflareDNS: getCloudflare,
		Compress:      compressstruct,
		Crontab:       getCrontab,
		GodaddyDNS:    getGodaddy,
		Ini:           getIni,
		JavascriptVM:  getJavascriptVM,
	}
}

type lockStruct struct {
	lock *sync.Mutex
}

func getLock() *lockStruct {
	var a sync.Mutex
	return &lockStruct{lock: &a}
}

func (m *lockStruct) acquire() {
	m.lock.Lock()
}

func (m *lockStruct) release() {
	m.lock.Unlock()
}
