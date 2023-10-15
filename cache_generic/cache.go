package cache_generic

import (
	"fmt"
	"reflect"
	"strings"

	myLocalCache "github.com/matyas-cyril/cache_file"
)

func New[C myLocalCache.CacheFile](category string, key []byte, ok, ko uint32, opt string) (*CacheGeneric[myLocalCache.CacheFile], error) {

	category = strings.ToUpper(strings.TrimSpace(category))

	switch category {

	case "LOCAL":
		lcl, err := myLocalCache.New(opt)
		if err != nil {
			return nil, err
		}

		// Activation du chiffrement
		if len(key) > 0 {
			lcl.SetKey(key)
			lcl.EnableCrypt()
		}

		c := CacheGeneric[myLocalCache.CacheFile]{
			f_cache: lcl,
		}
		c.key = key
		c.ok = ok
		c.ko = ko
		c.category = category

		return &c, nil

	case "MEMCACHE":
		return nil, fmt.Errorf("cache 'MEMCACHE' not yet implemented")

	case "REDIS":
		return nil, fmt.Errorf("cache 'REDIS' not yet implemented")

	}

	return nil, fmt.Errorf(fmt.Sprintf("cache '%s' not exist", category))
}

func (c *CacheGeneric[C]) SetSucces(data map[string][]byte, hashKey []byte) error {

	if c == nil {
		return fmt.Errorf("type var %s is nil", reflect.TypeOf(c).String())
	}

	return c.addInCache(data, hashKey, c.ok)

}

func (c *CacheGeneric[C]) SetFailed(data map[string][]byte, hashKey []byte) error {

	if c == nil {
		return fmt.Errorf("type var %s is nil", reflect.TypeOf(c).String())
	}

	return c.addInCache(data, hashKey, c.ko)
}

func (c *CacheGeneric[C]) GetCache(hashKey []byte) (map[string][]byte, error) {

	if c == nil {
		return nil, fmt.Errorf("type var %s is nil", reflect.TypeOf(c).String())
	}

	return c.getInCache(hashKey)
}

func (c *CacheGeneric[C]) addInCache(data map[string][]byte, hashKey []byte, length uint32) error {

	switch c.category {

	case "LOCAL":
		//var lcl localCache.CacheFile
		lcl := myLocalCache.CacheFile(*c.f_cache)
		err := lcl.Write(hashKey, data, uint(length))
		return err

	}

	return fmt.Errorf("errorooororororororo")
}

func (c *CacheGeneric[C]) getInCache(hashKey []byte) (map[string][]byte, error) {

	switch c.category {

	case "LOCAL":
		//var lcl localCache.CacheFile
		lcl := myLocalCache.CacheFile(*c.f_cache)

		// Obtenir le nom du fichier de cache
		fileName, err := lcl.GetFileName(hashKey)
		if err != nil {
			return nil, err
		}

		// Vérifier la validité du cache
		data, err := lcl.Read(fileName)
		if err != nil {
			return nil, err
		}
		return data, err

	}

	return nil, fmt.Errorf("errorooororororororo")
}

func (c *CacheGeneric[C]) Clean() (uint64, uint64, []error, error) {

	switch c.category {

	case "LOCAL":
		//var lcl localCache.CacheFile
		lcl := myLocalCache.CacheFile(*c.f_cache)
		return lcl.Sweep()

	}

	return 0, 0, nil, fmt.Errorf("errorooororororororo")

}

func (c *CacheGeneric[C]) Purge() (uint64, uint64, []error, error) {

	switch c.category {

	case "LOCAL":
		//var lcl localCache.CacheFile
		lcl := myLocalCache.CacheFile(*c.f_cache)
		return lcl.Purge()

	}

	return 0, 0, nil, fmt.Errorf("errorooororororororo")

}
