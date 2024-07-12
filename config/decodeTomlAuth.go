package config

import "fmt"

func (c *Config) decodeTomlAuth(d any) error {

	for name, v := range d.(map[string]any) {

		switch name {

		case "mech":
			d, err := castAnyToStringTab(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			c.Auth.MechList = d

		case "auth_multi":
			d, err := castBool(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			c.Auth.AuthMulti = d

		default:
			return fmt.Errorf(fmt.Sprintf("value '%s' of key [%s.%s] is not a valid hash option", d, name, v))

		}

	}

	// Si pas de plugin par défault alors "NO"
	if len(c.Auth.MechList) == 0 {
		c.Auth.MechList = []string{"NO"}
	}

	return nil
}
