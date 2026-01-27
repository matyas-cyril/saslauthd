package cache_generic

import (
	"fmt"
	"strings"
	"time"

	memcache "github.com/bradfitz/gomemcache/memcache"
	myLocalCache "github.com/matyas-cyril/cache-file"
	redis "github.com/redis/go-redis/v9"
)

func New(name string, key []byte, ok, ko uint32, opt []any) (cache *Cache, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			cache = nil
			err = fmt.Errorf("panic error : %s", pErr)
		}
	}()

	switch name = strings.ToUpper(strings.TrimSpace(name)); name {

	case "LOCAL":
		lcl, err := myLocalCache.New(opt[0].(string))
		if err != nil {
			return nil, err
		}

		// Activation du chiffrement
		if len(key) > 0 {
			lcl.SetKey(key)
			lcl.EnableCrypt()
		}

		return &Cache{
			name:    name,
			key:     key,
			ok:      ok,
			ko:      ko,
			f_local: lcl,
		}, nil

	case "MEMCACHE":

		mc := memcache.New(fmt.Sprintf("%s:%d", opt[0], opt[1]))
		defer mc.Close()

		mc.Timeout = time.Duration(opt[2].(int)) * time.Second

		return &Cache{
			name:       name,
			key:        key,
			ok:         ok,
			ko:         ko,
			f_memcache: mc,
		}, nil

	case "REDIS|KEYDB":
		r := redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6749),
			Password:     "",
			DB:           0,
			ReadTimeout:  time.Duration(opt[2].(int)) * time.Second,
			WriteTimeout: time.Duration(opt[2].(int)) * time.Second,
		})
		defer r.Close()

		return &Cache{
			name:    name,
			key:     key,
			ok:      ok,
			ko:      ko,
			f_redis: r,
		}, nil

	}

	return nil, fmt.Errorf("cache '%s' not exist", name)
}

func (c *Cache) SetSucces(data map[string][]byte, hashKey []byte) error {
	/*
		if c == nil {
			return fmt.Errorf("type var %s is nil", reflect.TypeFor[*Cache]().String())
		}
	*/
	return c.addInCache(data, hashKey, c.ok)

}

func (c *Cache) SetFailed(data map[string][]byte, hashKey []byte) error {
	/*
		if c == nil {
			return fmt.Errorf("type var %s is nil", reflect.TypeOf(c).String())
		}
	*/
	return c.addInCache(data, hashKey, c.ko)
}

func (c *Cache) addInCache(data map[string][]byte, hashKey []byte, length uint32) error {

	switch c.name {

	case "LOCAL":

		if err := c.f_local.Write(hashKey, data, uint(length)); err != nil {
			return err
		}

	default:
		return fmt.Errorf("failed to Add value in cache - cache type '%s' not exist", c.name)

	}

	return nil
}

func (c *Cache) GetCache(hashKey []byte) (map[string][]byte, error) {
	/*
		if c == nil {
			return nil, fmt.Errorf("type var %s is nil", reflect.TypeOf(c).String())
		}
	*/
	return c.getInCache(hashKey)
}

func (c *Cache) getInCache(hashKey []byte) (map[string][]byte, error) {

	data := map[string][]byte{}

	switch c.name {

	case "LOCAL":

		// Obtenir le nom du fichier de cache
		fileName, err := c.f_local.GetFileName(hashKey)
		if err != nil {
			return nil, err
		}

		// Vérifier la validité du cache
		data, err = c.f_local.Read(fileName)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("failed to Get value in cache - cache type '%s' not exist", c.name)

	}

	return data, nil
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
