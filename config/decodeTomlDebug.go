package config

import (
	"fmt"
	"path/filepath"
)

func (c *Config) decodeTomlDebug(d any) error {

	for name, v := range d.(map[string]any) {

		switch name {
		case "enable":
			d, err := castBool(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			c.Debug.Enable = d

		case "file":
			d, err := castString(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}

			f, err := filepath.Abs(d)
			if err != nil {
				return fmt.Errorf("key [%s.%s] : %s", name, v, err.Error())
			}
			c.Debug.File = f

		default:
			return fmt.Errorf("value '%s' of key [%s.%s] is not a valid hash option", d, name, v)

		}

	}

	return nil
}
