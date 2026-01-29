package cache_generic

import (
	"fmt"
)

func (c *Cache) SetSucces(data map[string][]byte, hashKey []byte) error {
	return c.addInCache(data, hashKey, c.ok)
}

func (c *Cache) SetFailed(data map[string][]byte, hashKey []byte) error {
	return c.addInCache(data, hashKey, c.ko)
}

func (c *Cache) GetCache(hashKey []byte) (map[string][]byte, error) {
	/*
		if c == nil {
			return nil, fmt.Errorf("type var %s is nil", reflect.TypeOf(c).String())
		}
	*/
	return c.getInCache(hashKey)
}

func (c *Cache) Clean() (uint64, uint64, []error, error) {

	switch c.name {

	case "LOCAL":
		return c.f_local.Sweep()

	default:
		return 0, 0, nil, fmt.Errorf("failed to Clean - cache type '%s' not exist", c.name)

	}

}

func (c *Cache) Purge() (uint64, uint64, []error, error) {

	switch c.name {

	case "LOCAL":
		return c.f_local.Purge()

	default:
		return 0, 0, nil, fmt.Errorf("failed to Purge - cache type '%s' not exist", c.name)

	}

}

func (c *Cache) Close() error {

	switch c.name {
	case "LOCAL":
		return nil

	case "MEMCACHE":
		return c.f_memcache.Close()

	case "REDIS", "KEYDB":
		return c.f_redis.Close()
	}

	return nil
}
