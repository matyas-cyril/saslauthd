package cache_generic

import (
	"fmt"
)

func (c *Cache) SetSucces(data map[string][]byte, hashKey []byte) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("panic error in SetSucces : %s", pErr)
		}
	}()
	return c.addInCache(data, hashKey, c.ok)
}

func (c *Cache) SetFailed(data map[string][]byte, hashKey []byte) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("panic error in SetFailed : %s", pErr)
		}
	}()
	return c.addInCache(data, hashKey, c.ko)
}

func (c *Cache) GetCache(hashKey []byte) (data map[string][]byte, err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			data = nil
			err = fmt.Errorf("panic error in GetCache : %s", pErr)
		}
	}()
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

func (c *Cache) Close() (err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("panic error in Close : %s", pErr)
		}
	}()

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
