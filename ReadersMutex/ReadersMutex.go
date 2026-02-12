package ReadersMutex

import "sync"

type RMutex struct {
	readersCount int
	readerMutex  sync.Mutex
	writerMutex  sync.Mutex
}

func (rm *RMutex) ReadLock() {
	/**
		We should lock readersMutex to only update the count.
		Otherwise a deadlock can arise here.
	 	Can also use atomic-int instead.
	**/
	rm.readerMutex.Lock()
	rm.readersCount++
	if rm.readersCount == 1 {
		rm.writerMutex.Lock()
	}
	rm.readerMutex.Unlock()
}

func (rm *RMutex) WriteLock() {
	rm.writerMutex.Lock()
}

func (rm *RMutex) ReadUnlock() {
	rm.readerMutex.Lock()
	rm.readersCount--
	if rm.readersCount == 0 {
		rm.writerMutex.Unlock()
	}
	rm.readerMutex.Unlock()
}

func (rm *RMutex) WriteUnlock() {
	rm.writerMutex.Unlock()
}

func (rw *RMutex) TryWriteLock() bool {
	return rw.writerMutex.TryLock()
}

func New() *RMutex {
	return &RMutex{
		readersCount: 0,
		readerMutex:  sync.Mutex{},
		writerMutex:  sync.Mutex{},
	}
}
