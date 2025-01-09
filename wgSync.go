package saslauthd

import "sync"

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

func (w *wgSync) Add(delta int) {
	w.mu.Lock()
	w.wg.Add(delta)
	w.cpt += uint64(delta)
	w.mu.Unlock()
}

func (w *wgSync) Done(delta int) {
	w.mu.Lock()
	if (w.cpt - uint64(delta)) >= 0 {
		for i := 0; i < delta; i++ {
			w.wg.Done()
			w.cpt -= 1
		}
	}
	w.mu.Unlock()
}

func (w *wgSync) Purge() {
	w.Done(int(w.Get()))
}

func (w *wgSync) Get() uint64 {
	return w.cpt
}

func (w *wgSync) Wait() {
	w.wg.Wait()
}
