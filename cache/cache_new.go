package cache_generic

import (
	"context"
	"fmt"
	"strings"
	"time"

	memcache "github.com/bradfitz/gomemcache/memcache"
	myLocalCache "github.com/matyas-cyril/cache-file"
	redis "github.com/redis/go-redis/v9"
)

/*
Initialisation du type de cache utilisé.
name est le type de cache utilisé (LOCAL, MEMCACHE, REDIS, KEYDB)
key est la clef de chiffrement du cache
ok et ko sont les durées en cas de succès ou d'échec
opt correspond aux options. Les clefs sont différentes en fonction du type de cache
*/
func New(name string, key []byte, ok, ko uint32, opt map[string]any) (cache *Cache, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			cache = nil
			err = fmt.Errorf("panic error in cache instantiation: %s", pErr)
		}
	}()

	switch name = strings.ToUpper(strings.TrimSpace(name)); name {

	case "LOCAL":

		var path string

		for k, v := range opt {

			cr := false

			switch k {

			case "path":
				path, cr = v.(string)
				if !cr {
					return nil, fmt.Errorf("cast error - failed to initialised cache %s", name)
				}

			default:
				return nil, fmt.Errorf("key '%s' invalid - failed to initialised cache %s", k, name)
			}

		}

		lcl, err := myLocalCache.New(path)
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

		var host string
		var port, timeout uint16

		for k, v := range opt {

			cr := false

			switch k {

			case "host":
				host, cr = v.(string)
				if !cr {
					return nil, fmt.Errorf("cast %s error - failed to initialised %s", k, name)
				}

			case "port":
				vInt, cr := v.(int)
				if !cr {
					return nil, fmt.Errorf("cast %s error - failed to initialised %s", k, name)
				}
				port = uint16(vInt)

			case "timeout":
				timeoutInt, cr := v.(int)
				if !cr {
					return nil, fmt.Errorf("cast %s error - failed to initialised %s", k, name)
				}
				timeout = uint16(timeoutInt)

			default:
				return nil, fmt.Errorf("key '%s' invalid - failed to initialised cache %s", k, name)

			}

		}

		mc := memcache.New(fmt.Sprintf("%s:%d", host, port))
		mc.Timeout = time.Duration(timeout) * time.Second
		if err := mc.Ping(); err != nil {
			return nil, fmt.Errorf("ping connection failed to %s:%d server", host, port)
		}

		return &Cache{
			name:       name,
			key:        key,
			ok:         ok,
			ko:         ko,
			f_memcache: mc,
		}, nil

	case "REDIS", "KEYDB":

		var host, username, password string
		var port, timeout uint16
		var db uint8

		for k, v := range opt {

			cr := false

			switch k {

			case "host":
				host, cr = v.(string)
				if !cr {
					return nil, fmt.Errorf("cast %s error - failed to initialised %s", k, name)
				}

			case "port":
				vInt, cr := v.(int)
				if !cr {
					return nil, fmt.Errorf("cast %s error - failed to initialised %s", k, name)
				}
				port = uint16(vInt)

			case "timeout":
				timeoutInt, cr := v.(int)
				if !cr {
					return nil, fmt.Errorf("cast %s error - failed to initialised %s", k, name)
				}
				timeout = uint16(timeoutInt)

			case "username":
				host, cr = v.(string)
				if !cr {
					return nil, fmt.Errorf("cast %s error - failed to initialised %s", k, name)
				}

			case "password":
				host, cr = v.(string)
				if !cr {
					return nil, fmt.Errorf("cast %s error - failed to initialised %s", k, name)
				}

			case "db":
				dbInt, cr := v.(int)
				if !cr {
					return nil, fmt.Errorf("cast %s error - failed to initialised %s", k, name)
				}
				db = uint8(dbInt)

			default:
				return nil, fmt.Errorf("key '%s' invalid - failed to initialised cache %s", k, name)

			}

		}

		r := redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			Username:     username,
			Password:     password,
			DB:           int(db),
			ReadTimeout:  time.Duration(timeout) * time.Second,
			WriteTimeout: time.Duration(timeout) * time.Second,
		})

		_, err := r.Ping(context.Background()).Result()
		if err != nil {
			return nil, fmt.Errorf("ping connection failed to %s:%d server", host, port)
		}

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
