package golanglibs

import "sync"

type toolsStruct struct {
	Lock func() *lockStruct
}

var Tools toolsStruct

func init() {
	Tools = toolsStruct{
		Lock: getLock,
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
