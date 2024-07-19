package cache_generic

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	myLocalCache "github.com/matyas-cyril/cache-file"
)

func New(name string, key []byte, ok, ko uint32, opt []any) (*CacheGeneric, error) {

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

		c := CacheGeneric{
			name:    name,
			key:     key,
			ok:      ok,
			ko:      ko,
			f_local: lcl,
		}

		return &c, nil

	case "MEMCACHE":

		mc := memcache.New(fmt.Sprintf("%s:%d", opt[0], opt[1]))

		mc.Timeout = time.Duration(opt[2].(int)) * time.Second

		c := CacheGeneric{
			name:       name,
			key:        key,
			ok:         ok,
			ko:         ko,
			f_memcache: mc,
		}

		return &c, nil

	case "REDIS":
		return nil, fmt.Errorf("cache 'REDIS' not yet implemented")

	}

	return nil, fmt.Errorf(fmt.Sprintf("cache '%s' not exist", name))
}

func (c *CacheGeneric) SetSucces(data map[string][]byte, hashKey []byte) error {

	if c == nil {
		return fmt.Errorf("type var %s is nil", reflect.TypeOf(c).String())
	}

	return c.addInCache(data, hashKey, c.ok)

}

func (c *CacheGeneric) SetFailed(data map[string][]byte, hashKey []byte) error {

	if c == nil {
		return fmt.Errorf("type var %s is nil", reflect.TypeOf(c).String())
	}

	return c.addInCache(data, hashKey, c.ko)
}

func (c *CacheGeneric) GetCache(hashKey []byte) (map[string][]byte, error) {

	if c == nil {
		return nil, fmt.Errorf("type var %s is nil", reflect.TypeOf(c).String())
	}

	return c.getInCache(hashKey)
}

func (c *CacheGeneric) addInCache(data map[string][]byte, hashKey []byte, length uint32) error {

	switch c.name {

	case "LOCAL":

		if err := c.f_local.Write(hashKey, data, uint(length)); err != nil {
			return err
		}

	}

	return nil
}

func (c *CacheGeneric) getInCache(hashKey []byte) (map[string][]byte, error) {

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

	}

	return data, nil
}

func (c *CacheGeneric) Clean() (uint64, uint64, []error, error) {

	switch c.name {

	case "LOCAL":
		return c.f_local.Sweep()

	}

	return 0, 0, nil, fmt.Errorf("errorooororororororo")

}

func (c *CacheGeneric) Purge() (uint64, uint64, []error, error) {

	switch c.name {

	case "LOCAL":
		return c.f_local.Purge()

	}

	return 0, 0, nil, fmt.Errorf("errorooororororororo")

}
