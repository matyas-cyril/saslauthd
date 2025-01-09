package cache_generic

import (
	"github.com/bradfitz/gomemcache/memcache"
	myLocalCache "github.com/matyas-cyril/cache-file"
)

/*
	type CacheType interface {
		myLocalCache.CacheFile | memcache.Client
	}

	type CacheGeneric[Cachetype CacheType] struct {
		name    string // Nom du cache utilisé
		key     []byte // clef de chiffrement
		ok      uint32
		ko      uint32
		f_cache *Cachetype
	}
*/
type CacheGeneric struct {
	name       string // Nom du cache utilisé
	key        []byte // clef de chiffrement
	ok         uint32
	ko         uint32
	f_local    *myLocalCache.CacheFile
	f_memcache *memcache.Client
}
