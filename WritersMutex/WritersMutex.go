package Writersmutex

import (
	"sync"
)

type WPMutex struct {
	totalPendingWriters int
	totalActiveReaders  int
	isAnyWriteActive    bool
	cond                sync.Cond
}

func (wm *WPMutex) ReadLock() {
	wm.cond.L.Lock()
	for wm.totalPendingWriters != 0 || wm.isAnyWriteActive {
		wm.cond.Wait()
	}
	wm.totalActiveReaders++
	wm.cond.L.Unlock()
}

func (wm *WPMutex) WriteLock() {
	wm.cond.L.Lock()
	wm.totalPendingWriters++
	for wm.totalActiveReaders != 0 || wm.isAnyWriteActive {
		wm.cond.Wait()
	}
	wm.totalPendingWriters--
	wm.isAnyWriteActive = true
	wm.cond.L.Unlock()
}

func (wm *WPMutex) ReadUnlock() {
	wm.cond.L.Lock()
	wm.totalActiveReaders--
	if wm.totalActiveReaders == 0 {
		wm.cond.Broadcast()
	}
	wm.cond.L.Unlock()
}

func (wm *WPMutex) WriteUnlock() {
	wm.cond.L.Lock()
	wm.isAnyWriteActive = false
	wm.cond.Broadcast()
	wm.cond.L.Unlock()
}

func New() *WPMutex {
	return &WPMutex{
		cond: *sync.NewCond(&sync.Mutex{}),
	}
}
