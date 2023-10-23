package cache_generic

import (
	localCache "github.com/matyas-cyril/cache-file"
)

type CacheType interface {
	localCache.CacheFile
}

type CacheGeneric[C CacheType] struct {
	category string
	key      []byte
	ok       uint32
	ko       uint32
	f_cache  *C
}
