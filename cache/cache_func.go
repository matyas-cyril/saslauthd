package cache_generic

import "fmt"

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
