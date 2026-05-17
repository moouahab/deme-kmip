package metrics

import (
	"context"
	"sync"
)

type MemoryCollector struct {
	mu sync.RWMutex

	httpRequestsTotal int
	httpErrorsTotal   int

	createKeyTotal  int
	getKeyTotal     int
	destroyKeyTotal int

	successTotal  int
	notFoundTotal int
}

func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{}
}

func (c *MemoryCollector) IncHTTPRequests(ctx context.Context) {
	_ = ctx
	c.mu.Lock()
	defer c.mu.Unlock()
	c.httpRequestsTotal++
}

func (c *MemoryCollector) IncHTTPErrors(ctx context.Context) {
	_ = ctx
	c.mu.Lock()
	defer c.mu.Unlock()
	c.httpErrorsTotal++
}

func (c *MemoryCollector) IncCreateKey(ctx context.Context) {
	_ = ctx
	c.mu.Lock()
	defer c.mu.Unlock()
	c.createKeyTotal++
}

func (c *MemoryCollector) IncGetKey(ctx context.Context) {
	_ = ctx
	c.mu.Lock()
	defer c.mu.Unlock()
	c.getKeyTotal++
}

func (c *MemoryCollector) IncDestroyKey(ctx context.Context) {
	_ = ctx
	c.mu.Lock()
	defer c.mu.Unlock()
	c.destroyKeyTotal++
}

func (c *MemoryCollector) IncSuccess(ctx context.Context) {
	_ = ctx
	c.mu.Lock()
	defer c.mu.Unlock()
	c.successTotal++
}

func (c *MemoryCollector) IncNotFound(ctx context.Context) {
	_ = ctx
	c.mu.Lock()
	defer c.mu.Unlock()
	c.notFoundTotal++
}

func (c *MemoryCollector) Snapshot(ctx context.Context) Snapshot {
	_ = ctx

	c.mu.RLock()
	defer c.mu.RUnlock()

	return Snapshot{
		HTTPRequestsTotal: c.httpRequestsTotal,
		HTTPErrorsTotal:   c.httpErrorsTotal,
		CreateKeyTotal:    c.createKeyTotal,
		GetKeyTotal:       c.getKeyTotal,
		DestroyKeyTotal:   c.destroyKeyTotal,
		SuccessTotal:      c.successTotal,
		NotFoundTotal:     c.notFoundTotal,
	}
}
