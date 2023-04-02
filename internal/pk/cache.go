package pk

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"sync"
)

const CacheWindow = 10000 // 10k is sufficient

type (
	Cache struct {
		pks sync.Map // implies usage by pointer, but allows init via new(Cache)
	}
	TableCache = *lru.Cache[string, struct{}]
)

func (c *Cache) For(tableName string) TableCache {
	v, ok := c.pks.Load(tableName)
	if ok {
		return v.(TableCache)
	}

	cache, _ := lru.New[string, struct{}](CacheWindow)
	v, _ = c.pks.LoadOrStore(tableName, cache)
	return v.(TableCache)
}

func (c *Cache) Done(tableName string) {
	c.pks.Delete(tableName)
}
