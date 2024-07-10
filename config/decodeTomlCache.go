package config

import (
	"fmt"
	"strings"
)

func (c *Config) decodeTomlCache(d any) error {

	for name, v := range d.(map[string]any) {

		switch name {

		case "enable", "keyRand":
			d, err := castBool(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}

			switch v {
			case "enable":
				c.Cache.Enable = d

			case "keyRand":
				c.Cache.KeyRand = d

			}

		case "key":
			d, err := castString(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			c.Cache.Key = []byte(strings.TrimSpace(d))

		case "type":
			d, err := castString(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}

			switch strings.ToUpper(d) {
			case "LOCAL":
				c.Cache.Category = "LOCAL"
				if err := decodeTomlCacheLocal(v); err != nil {

				}

			case "MEMCACHE":
				c.Cache.Category = "MEMCACHE"
			case "REDIS":
				c.Cache.Category = "REDIS"
			default:
				return fmt.Errorf(fmt.Sprintf("value '%s' of key [%s.%s] not valid must be [ MEM | FILE | MEMCACHED | REDIS ]", d, name, v))
			}

		case "ok", "ko":
			d, err := castUint32(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}

			if d > 31536000 {
				return fmt.Errorf(fmt.Sprintf("key [%s.%s] must be lower than or equal 31536000", name, v))
			}

			switch name {
			case "ok":
				c.Cache.OK = d

			case "ko":
				c.Cache.KO = d
			}

		default:
			return fmt.Errorf(fmt.Sprintf("value '%s' of key [%s.%s] is not a valid hash option", d, name, v))

		}

	}

	return nil
}

func decodeTomlCacheLocal(d any) error {

	return nil
}
