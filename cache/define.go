package cache_generic

import (
	memcache "github.com/bradfitz/gomemcache/memcache"
	myLocalCache "github.com/matyas-cyril/cache-file"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	name       string // Nom du cache utilisé
	key        []byte // clef de chiffrement
	ok         uint32
	ko         uint32
	f_local    *myLocalCache.CacheFile
	f_memcache *memcache.Client // memcached
	f_redis    *redis.Client    // Redis ou KeyDB
}
