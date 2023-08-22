package server

import (
	"jkwksvr/jwks"
	"sync"
	"time"
)

type JWKSCache struct {
	Remote      string
	Filter      func(jwks *jwks.JWKS) bool
	RefreshTime time.Duration

	mut         sync.RWMutex
	lastRefresh time.Time
	cache       []byte
}

func (c *JWKSCache) Refresh() {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.refresh()
}

func (c *JWKSCache) refresh() {
	// Get the jwks from remote.
	// Unmarshal and filter
	// Marshal and cache.
}

func (c JWKSCache) Bytes() []byte {
	c.mut.RLock()
	if c.cache != nil {
		c.mut.RUnlock()
		return c.cache
	}
	// Need a refresh.
	c.mut.Lock()
	if c.cache != nil {
		// someone was faster
		c.mut.Unlock()
		return c.cache
	}
	c.refresh()
	c.mut.Unlock()
	return c.cache
}

