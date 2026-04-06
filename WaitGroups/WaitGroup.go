package WaitGroup

import "sync"

type WGroup struct {
	cond      *sync.Cond
	groupSize uint // unsigned, hence not allowing to be negative
}

func (wg *WGroup) Add(delta uint) {
	if delta < 1 {
		panic("Delta must greater than/equal to 1")
	}
	wg.cond.L.Lock()
	wg.groupSize += delta
	wg.cond.L.Unlock()
}

func (wg *WGroup) Done() {
	if wg.groupSize == 0 {
		panic("Done() called more times than Add()")
	}
	wg.cond.L.Lock()
	wg.groupSize--
	if wg.groupSize == 0 {
		// Broadcast since more than routines might be waiting
		wg.cond.Broadcast()
	}
	wg.cond.L.Unlock()
}

func (wg *WGroup) Wait() {
	wg.cond.L.Lock()
	for wg.groupSize > 0 {
		wg.cond.Wait()
	}
	wg.cond.L.Unlock()
}

func NewWGroup() *WGroup {
	return &WGroup{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}
