package golanglibs

import "sync"

type LockStruct struct {
	lock *sync.Mutex
}

func getLock() *LockStruct {
	var a sync.Mutex
	return &LockStruct{lock: &a}
}

func (m *LockStruct) Acquire() {
	m.lock.Lock()
}

func (m *LockStruct) Release() {
	m.lock.Unlock()
}
