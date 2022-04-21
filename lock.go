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

// 写锁定情况下，对读写锁进行读锁定或者写锁定，都将阻塞；而且读锁与写锁之间是互斥的；
// 读锁定情况下，对读写锁进行写锁定，将阻塞；加读锁时不会阻塞；
// 对未被写锁定的读写锁进行写解锁，会引发 Panic；
// 对未被读锁定的读写锁进行读解锁的时候也会引发 Panic；
// 写解锁在进行的同时会试图唤醒所有因进行读锁定而被阻塞的 goroutine；
// 读解锁在进行的时候则会试图唤醒一个因进行写锁定而被阻塞的 goroutine。

type RWLockStruct struct {
	lock *sync.RWMutex
}

func getRWLock() *RWLockStruct {
	var rwMutex sync.RWMutex
	return &RWLockStruct{lock: &rwMutex}
}

// 获取读锁
func (m *RWLockStruct) RAcquire() {
	m.lock.RLock()
}

// 释放读锁
func (m *RWLockStruct) RRelease() {
	m.lock.RUnlock()
}

// 获取写锁
func (m *RWLockStruct) WAcquire() {
	m.lock.Lock()
}

// 释放写锁
func (m *RWLockStruct) WRelease() {
	m.lock.Unlock()
}
