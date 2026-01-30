package saslauthd

import (
	"sync"
)

type wgSync struct {
	mu  sync.Mutex
	wg  *sync.WaitGroup
	cpt uint64
}

func NewSync() *wgSync {
	return &wgSync{
		wg: &sync.WaitGroup{},
	}
}

func (w *wgSync) Add(delta uint) {
	w.mu.Lock()
	w.wg.Add(int(delta))
	w.cpt += uint64(delta)
	w.mu.Unlock()
}

func (w *wgSync) Done(delta uint) {
	w.mu.Lock()
	if w.cpt >= uint64(delta) {
		for range delta {
			w.wg.Done()
			w.cpt -= 1
		}
	}

	w.mu.Unlock()
}

func (w *wgSync) Purge() {
	w.Done(uint(w.Get()))
}

func (w *wgSync) Get() uint64 {
	return w.cpt
}

func (w *wgSync) Wait() {
	w.wg.Wait()
}
