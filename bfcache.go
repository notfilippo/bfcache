package bfcache

import (
	"sync"

	"github.com/cespare/xxhash/v2"
)

type Cache struct {
	Map sync.Map
}

func New() *Cache {
	return &Cache{Map: sync.Map{}}
}

func (c *Cache) Get(k []byte) []byte {
	sum := xxhash.Sum64(k)
	value, ok := c.Map.Load(sum)
	if !ok {
		return nil
	}
	return value.([]byte)
}

func (c *Cache) Set(k []byte, v []byte) {
	sum := xxhash.Sum64(k)
	c.Map.Store(sum, v)
}
