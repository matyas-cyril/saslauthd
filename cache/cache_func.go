package cache_generic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	memcache "github.com/bradfitz/gomemcache/memcache"
)

func (c *Cache) addInCache(data map[string][]byte, hashKey []byte, exp uint32) error {

	switch c.name {

	case "LOCAL":

		if err := c.f_local.Write(hashKey, data, uint(exp)); err != nil {
			return err
		}

	case "MEMCACHE":

		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to cast data for %s value - %v", c.name, err)
		}

		if err := c.f_memcache.Add(&memcache.Item{
			Key:        string(hashKey),
			Value:      jsonData,
			Expiration: int32(time.Now().Add(time.Duration(exp) * time.Second).Unix()),
		}); err != nil {
			return fmt.Errorf("failed to add data to %s : %v", c.name, err)
		}

	case "REDIS":

		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to cast data for %s value - %v", c.name, err)
		}

		if err := c.f_redis.Set(context.Background(), string(hashKey), jsonData, time.Duration(exp)*time.Second).Err(); err != nil {
			return fmt.Errorf("failed to add data to %s : %v", c.name, err)
		}

	default:
		return fmt.Errorf("failed to Add value in cache type '%s' not available", c.name)

	}

	return nil
}

func (c *Cache) getInCache(hashKey []byte) (map[string][]byte, error) {

	var data map[string][]byte

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

	case "MEMCACHE":

		item, err := c.f_memcache.Get(string(hashKey))
		if err != nil {
			return nil, fmt.Errorf("failed to get data from %s - %v", c.name, err)
		}

		if err := json.Unmarshal(item.Value, &data); err != nil {
			return nil, fmt.Errorf("failed to cast data from %s value - %v", c.name, err)
		}

	case "REDIS":

		item, err := c.f_redis.Get(context.Background(), string(hashKey)).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get data from %s - %v", c.name, err)
		}

		if err := json.Unmarshal([]byte(item), &data); err != nil {
			return nil, fmt.Errorf("failed to cast data from %s value - %v", c.name, err)
		}

	default:
		return nil, fmt.Errorf("failed to Get value in cache type '%s' not exist", c.name)

	}

	return data, nil
}
