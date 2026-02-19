package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/matyas-cyril/logme"
)

func (c *Config) decodeTomlServer(d any) error {

	for name, v := range d.(map[string]any) {

		switch name {

		case "socket", "user", "group", "plugin_path", "buffer_hash", "ugo":
			d, err := castString(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}

			switch name {
			case "socket":
				c.Server.Socket = d

			case "user":
				c.Server.User = d

			case "group":
				c.Server.Group = d

			case "ugo":
				parsedMode, err := strconv.ParseUint(d, 8, 32)
				if err != nil {
					return fmt.Errorf("value '%s' of key [SERVER.%s] is not a valid", d, name)
				}
				c.Server.UGO = os.FileMode(parsedMode)

			case "plugin_path":
				if !dirExist(d) {
					return fmt.Errorf("value '%s' of key [SERVER.%s] is not a valid directory or not exist", d, name)
				}
				c.Server.PluginPath = d

			case "buffer_hash":
				switch strings.ToUpper(d) {
				case "MD5":
					c.Server.BufferHashType = 0

				case "SHA1":
					c.Server.BufferHashType = 1

				case "SHA256":
					c.Server.BufferHashType = 2

				case "SHA512":
					c.Server.BufferHashType = 3

				default:
					return fmt.Errorf("value '%s' of key [SERVER.%s] is not a valid hash option", d, name)

				}

			}

		case "log":
			d, err := castString(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}

			switch strings.ToUpper(d) {
			case "NO":
				c.Server.LogType = logme.LOGME_NO

			case "TERM":
				c.Server.LogType = logme.LOGME_TERM

			case "SYSLOG":
				c.Server.LogType = logme.LOGME_SYSLOG

			case "BOTH":
				c.Server.LogType = logme.LOGME_BOTH

			default:
				return fmt.Errorf("value '%s' of key [%s.%s] not valid must be [ NO | TERM | SYSLOG | BOTH ]", d, name, v)
			}

		case "log_facility":
			d, err := castString(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}

			switch strings.ToUpper(d) {
			case "AUTH":
				c.Server.LogFacility = logme.LOGME_F_AUTH
			case "MAIL":
				c.Server.LogFacility = logme.LOGME_F_MAIL
			case "SYSLOG":
				c.Server.LogFacility = logme.LOGME_F_SYSLOG
			case "USER":
				c.Server.LogFacility = logme.LOGME_F_USER
			case "LOCAL0":
				c.Server.LogFacility = logme.LOGME_F_LOCAL0
			case "LOCAL1":
				c.Server.LogFacility = logme.LOGME_F_LOCAL1
			case "LOCAL2":
				c.Server.LogFacility = logme.LOGME_F_LOCAL2
			case "LOCAL3":
				c.Server.LogFacility = logme.LOGME_F_LOCAL3
			case "LOCAL4":
				c.Server.LogFacility = logme.LOGME_F_LOCAL4
			case "LOCAL5":
				c.Server.LogFacility = logme.LOGME_F_LOCAL5
			case "LOCAL6":
				c.Server.LogFacility = logme.LOGME_F_LOCAL6
			case "LOCAL7":
				c.Server.LogFacility = logme.LOGME_F_LOCAL7
			default:
				return fmt.Errorf("value '%s' of key [%s.%s] not valid must be [ AUTH | MAIL | SYSLOG | USER |	LOCAL0 | LOCAL1 | LOCAL2 | LOCAL3 | LOCAL4 | LOCAL5 | LOCAL6 | LOCAL7 ]", d, name, v)
			}

		case "rate_info":
			d, err := castUint16(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			if d > 3600 {
				return fmt.Errorf("key [%s.%s] must be lower than or equal 3600", name, v)
			}
			c.Server.RateInfo = d

		case "client_max":
			d, err := castUint32(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			c.Server.ClientMax = d

		case "client_timeout":
			d, err := castUint8(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			c.Server.ClientTimeout = d

		case "buffer_size":
			d, err := castUint16(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			if d < 1 || d > 2048 {
				return fmt.Errorf("key [%s.%s] must be upper than 1 and lower than or equal 2048", name, v)
			}
			c.Server.BufferSize = d

		case "buffer_timeout":
			d, err := castUint16(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			if d < 1 || d > 10000 {
				return fmt.Errorf("key [%s.%s] must be upper than 1 and lower than or equal 10000", name, v)
			}
			c.Server.BufferTimeout = d

		case "socket_size":
			d, err := castUint16(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			if d < 8 || d > 32768 {
				return fmt.Errorf("key [%s.%s] must be upper than 8 and lower than or equal 32768", name, v)
			}
			c.Server.SocketSize = d

		case "graceful":
			d, err := castUint8(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			if d > 60 {
				return fmt.Errorf("key [%s.%s] must be upper than 0 and lower than or equal 60", name, v)
			}
			c.Server.Graceful = d

		case "notify":
			d, err := castBool(v)
			if err != nil {
				return fmt.Errorf("SERVER.%s - %s", name, err)
			}
			c.Server.Notify = d

		default:
			return fmt.Errorf("key SERVER.%s not exist", name)

		}

	}

	return nil
}
