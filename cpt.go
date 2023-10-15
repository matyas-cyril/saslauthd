package saslauthd

import "sync"

type cptSync struct {
	mu  sync.Mutex
	cpt uint64
}

func NewCpt() *cptSync {
	return &cptSync{}
}

func (c *cptSync) Inc() {

	defer func() {
		c.mu.Unlock()
	}()

	c.mu.Lock()
	c.cpt++

}

func (c *cptSync) Get() uint64 {
	return c.cpt
}

func (c *cptSync) Reset() {

	defer func() {
		c.mu.Unlock()
	}()

	c.mu.Lock()
	c.cpt = 0

}
