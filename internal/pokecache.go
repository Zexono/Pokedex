package internal

import (
	"sync"
	"time"
)


type Cache struct{
	mu sync.Mutex
	hold 	  map[string]cacheEntry
	interval  time.Duration
}


type cacheEntry struct{
	createdAt time.Time
	val 	  []byte
}


func (c *Cache) NewCache(interval time.Duration) *Cache{

	cache := Cache{
		hold: map[string]cacheEntry{},
		interval: interval,
	}
	go c.reapLoop()
	return &cache
}

func (c *Cache) Add(key string,val []byte)  {
	c.mu.Lock()
	defer c.mu.Unlock()
	ce := cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	c.hold[key] = ce
	
}

func (c *Cache) Get(key string) (val []byte,found bool)  {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.hold[key]
	if ok {
		return entry.val,true
	}

	return nil,false
	
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mu.Lock()
		for key, v := range c.hold {
			if time.Since(v.createdAt) > c.interval {
				delete(c.hold, key)
			}
		}
		c.mu.Unlock()
	}
}

