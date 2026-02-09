package cache_generic

import (
	"context"
	"fmt"
)

// Ajouter les données avec la durée spécique de cache pour succès d'auth (ok)
func (c *Cache) SetSucces(data map[string][]byte, hashKey []byte) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("panic error in SetSucces : %s", pErr)
		}
	}()
	return c.addInCache(data, hashKey, c.ok)
}

// Ajouter les données avec la durée spécique de cache pour échec d'auth (ko)
func (c *Cache) SetFailed(data map[string][]byte, hashKey []byte) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("panic error in SetFailed : %s", pErr)
		}
	}()
	return c.addInCache(data, hashKey, c.ko)
}

// Obtenir des données si elles sont présentes en cache
func (c *Cache) GetCache(hashKey []byte) (data map[string][]byte, err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			data = nil
			err = fmt.Errorf("panic error in GetCache : %s", pErr)
		}
	}()
	return c.getInCache(hashKey)
}

// Supprime toutes les informations dans le cache.
func (c *Cache) Flush() (err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("panic error in Flush : %s", pErr)
		}
	}()

	switch c.name {

	case "LOCAL":
		_, _, _, localErr := c.f_local.Purge()
		return localErr

	case "MEMCACHE": // Vide tout memcache
		return c.f_memcache.DeleteAll()

	case "REDIS": // Vide toutes les infomations correspondant à la DB
		return c.f_redis.FlushDB(context.Background()).Err()

	default:
		return fmt.Errorf("failed to Flush - cache type '%s' not exist", c.name)

	}

}

// Supprimer les fichiers invalides ou dépassés
func (c *Cache) Clean() (uint64, uint64, []error, error) {

	switch c.name {

	case "LOCAL":
		return c.f_local.Sweep()

	case "MEMCACHE", "REDIS": // C'est auto--géré
		return 0, 0, nil, nil

	default:
		return 0, 0, nil, fmt.Errorf("failed to Clean - cache type '%s' not exist", c.name)

	}

}

// Fermer la connexion à un serveur de cache
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

	case "REDIS":
		return c.f_redis.Close()
	}

	return nil
}
