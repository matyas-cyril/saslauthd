package config

import (
	"crypto/sha512"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/matyas-cyril/logme"
	toml "github.com/pelletier/go-toml"
)

func LoadConfig(tomFile, appName, appPath string) (*Config, error) {

	// Charger le contenu du fichier toml
	tomlTree, err := toml.LoadFile(tomFile)
	if err != nil {
		return nil, err
	}

	// Décoder le toml, si succès alors on a une structure
	conf, err := decodeToml(tomlTree, appPath)
	if err != nil {
		return nil, err
	}

	// Traitement des données après 1er contrôle
	conf, err = processConfig(conf, appName)
	if err != nil {
		return nil, err
	}

	// Traitement des plugins
	conf, err = processPlugin(tomlTree, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func decodeToml(tomlTree *toml.Tree, appPath string) (*Config, error) {

	// Déclaration des paramètres par défaut
	c := Config{}

	// Détermination du Path d'install
	c.Server.Network = "unix"
	c.Server.Socket = "/var/run/saslauthd/mux"
	c.Server.User = "mail"
	c.Server.Group = "mail"
	c.Server.RateInfo = 30
	c.Server.ClientMax = 100
	c.Server.ClientTimeout = 30
	c.Server.BufferSize = 256
	c.Server.BufferTimeout = 50
	c.Server.BufferHashType = 2
	c.Server.SocketSize = 1024
	c.Server.LogType = logme.LOGME_TERM
	c.Server.LogFacility = logme.LOGME_F_AUTH
	c.Server.Stat = 60 // Valeur non accessible via le fichier externe Toml pour l'instant
	c.Server.PluginPath = fmt.Sprintf("%s/plugins", appPath)

	c.Debug.File = "/tmp/saslauthd.debug"

	c.Cache.Category = "LOCAL"
	c.Cache.OK = 60
	c.Cache.KO = 60

	c.Cache.Local.Path = "/tmp"
	c.Cache.Local.Sweep = 60

	c.Auth.MechList = []string{"NO"}
	c.Auth.Plugin = make(map[string]*DefinePlugin)

	for _, k := range tomlTree.Keys() {

		if k == "SERVER" {

			// Parser en fonction des clefs définies dans le fichier Toml
			for _, v := range tomlTree.Get(k).(*toml.Tree).Keys() {

				opt := fmt.Sprintf("%s.%s", k, v)

				if v == "socket" || v == "user" || v == "group" || v == "plugin_path" || v == "buffer_hash" {

					d, err := getValue(tomlTree, opt, "string")
					if err != nil {
						return nil, err
					}

					if d != nil && len(strings.TrimSpace(d.(string))) > 0 {

						switch v {
						case "socket":
							c.Server.Socket = strings.TrimSpace(d.(string))

						case "user":
							c.Server.User = strings.TrimSpace(d.(string))

						case "group":
							c.Server.Group = strings.TrimSpace(d.(string))

						case "plugin_path":
							if !dirExist(strings.TrimSpace(d.(string))) {
								return nil, fmt.Errorf(fmt.Sprintf("value '%s' of key [%s.%s] is not a valid directory or not exist", strings.TrimSpace(d.(string)), k, v))
							}
							c.Server.PluginPath = strings.TrimSpace(d.(string))

						case "buffer_hash":
							switch strings.ToUpper(strings.TrimSpace(d.(string))) {
							case "MD5":
								c.Server.BufferHashType = 0

							case "SHA1":
								c.Server.BufferHashType = 1

							case "SHA256":
								c.Server.BufferHashType = 2

							case "SHA512":
								c.Server.BufferHashType = 3

							default:
								return nil, fmt.Errorf(fmt.Sprintf("value '%s' of key [%s.%s] is not a valid hash option", strings.TrimSpace(d.(string)), k, v))

							}

						}

					} else {
						return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be not null", k, v))
					}

				} else if v == "rate_info" || v == "client_max" || v == "client_timeout" || v == "buffer_size" || v == "buffer_timeout" || v == "socket_size" {

					d, err := getValue(tomlTree, opt, "int64")
					if err != nil {
						return nil, err
					}

					if d != nil && d.(int64) >= 0 {

						switch v {

						case "rate_info":
							if d.(int64) > 3600 {
								return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be lower than or equal 3600", k, v))
							}
							c.Server.RateInfo = uint16(d.(int64))

						case "client_max":
							if d.(int64) > 500000000 {
								return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be lower than or equal 500000000", k, v))
							}
							c.Server.ClientMax = uint32(d.(int64))

						case "client_timeout":
							if d.(int64) > 240 {
								return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be lower than or equal 240", k, v))
							}
							c.Server.ClientTimeout = uint8(d.(int64))

						case "buffer_size":
							if d.(int64) < 1 || d.(int64) > 2048 {
								return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be upper than 1 and lower than or equal 2048", k, v))
							}
							c.Server.BufferSize = uint16(d.(int64))

						case "buffer_timeout":
							if d.(int64) < 1 || d.(int64) > 10000 {
								return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be upper than 1 and lower than or equal 10000", k, v))
							}
							c.Server.BufferTimeout = uint16(d.(int64))

						case "socket_size":
							if d.(int64) < 8 || d.(int64) > 32768 {
								return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be upper than 8 and lower than or equal 32768", k, v))
							}
							c.Server.SocketSize = uint16(d.(int64))

						}

					} else {
						return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be greater than 0", k, v))
					}

				} else if v == "log" {

					d, err := getValue(tomlTree, opt, "string")
					if err != nil {
						return nil, err
					}

					if d != nil && len(strings.TrimSpace(d.(string))) > 0 {

						switch strings.ToUpper(strings.TrimSpace(d.(string))) {
						case "NO":
							c.Server.LogType = logme.LOGME_NO

						case "TERM":
							c.Server.LogType = logme.LOGME_TERM

						case "SYSLOG":
							c.Server.LogType = logme.LOGME_SYSLOG

						case "BOTH":
							c.Server.LogType = logme.LOGME_BOTH

						default:
							return nil, fmt.Errorf(fmt.Sprintf("value '%s' of key [%s.%s] not valid must be [ NO | TERM | SYSLOG | BOTH ]", d, k, v))
						}

					} else {
						return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be not null", k, v))
					}

				} else if v == "log_facility" {

					d, err := getValue(tomlTree, opt, "string")
					if err != nil {
						return nil, err
					}

					if d != nil && len(strings.TrimSpace(d.(string))) > 0 {

						switch strings.ToUpper(strings.TrimSpace(d.(string))) {
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
							return nil, fmt.Errorf(fmt.Sprintf("value '%s' of key [%s.%s] not valid must be [ AUTH | MAIL | SYSLOG | USER |	LOCAL0 | LOCAL1 | LOCAL2 | LOCAL3 | LOCAL4 | LOCAL5 | LOCAL6 | LOCAL7 ]", d, k, v))
						}

					} else {
						return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be not null", k, v))
					}

				} else {
					return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] not exist", k, v))

				}

			}

		} else if k == "DEBUG" {

			// Parser en fonction des clefs définies dans le fichier Toml
			for _, v := range tomlTree.Get(k).(*toml.Tree).Keys() {

				opt := fmt.Sprintf("%s.%s", k, v)

				if v == "enable" {

					d, err := getValue(tomlTree, opt, "bool")
					if err != nil {
						return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] : %s", k, v, err.Error()))
					}
					if d != nil {
						c.Debug.Enable = d.(bool)
					}

				} else if v == "file" {

					d, err := getValue(tomlTree, opt, "string")
					if err != nil {
						return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] : %s", k, v, err.Error()))
					}

					f, err := filepath.Abs(d.(string))
					if err != nil {
						return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] : %s", k, v, err.Error()))
					}

					c.Debug.File = f

				} else {
					return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] not exist", k, v))
				}

			}

		} else if k == "CACHE" {

			// Parser en fonction des clefs définies dans le fichier Toml
			for _, v := range tomlTree.Get(k).(*toml.Tree).Keys() {

				opt := fmt.Sprintf("%s.%s", k, v)

				if v == "enable" || v == "keyRand" {

					d, err := getValue(tomlTree, opt, "bool")
					if err != nil {
						return nil, err
					}
					if d != nil {

						switch v {
						case "enable":
							c.Cache.Enable = d.(bool)

						case "keyRand":
							c.Cache.KeyRand = d.(bool)

						}
					}

				} else if v == "key" {

					d, err := getValue(tomlTree, opt, "string")
					if err != nil {
						return nil, err
					}
					if d != nil {
						c.Cache.Key = []byte(strings.TrimSpace(d.(string)))
					}

				} else if v == "type" {

					d, err := getValue(tomlTree, opt, "string")
					if err != nil {
						return nil, err
					}

					if d != nil && len(strings.TrimSpace(d.(string))) > 0 {

						switch strings.ToUpper(strings.TrimSpace(d.(string))) {
						case "LOCAL":
							c.Cache.Category = "LOCAL"
						case "MEMCACHE":
							c.Cache.Category = "MEMCACHE"
						case "REDIS":
							c.Cache.Category = "REDIS"
						default:
							return nil, fmt.Errorf(fmt.Sprintf("value '%s' of key [%s.%s] not valid must be [ MEM | FILE | MEMCACHED | REDIS ]", d, k, v))
						}

					} else {
						return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be 'MEM' currently", k, v))
					}

				} else if v == "ok" || v == "ko" {

					d, err := getValue(tomlTree, opt, "int64")
					if err != nil {
						return nil, err
					}

					if d != nil && d.(int64) >= 0 {

						switch v {
						case "ok":
							if d.(int64) > 31536000 {
								return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be lower than or equal 31536000", k, v))
							}
							c.Cache.OK = uint32(d.(int64))

						case "ko":
							if d.(int64) > 31536000 {
								return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be lower than or equal 31536000", k, v))
							}
							c.Cache.KO = uint32(d.(int64))

						}

					} else {
						return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be greater than 0", k, v))
					}

				} else if v == "LOCAL" { // Sous-Section : CACHE.LOCAL

					for _, ss_v := range tomlTree.Get(opt).(*toml.Tree).Keys() {

						ss_opt := fmt.Sprintf("%s.%s", opt, ss_v)

						if ss_v == "path" {

							d, err := getValue(tomlTree, ss_opt, "string")
							if err != nil {
								return nil, err
							}

							if d != nil && len(strings.TrimSpace(d.(string))) > 0 {
								c.Cache.Local.Path = strings.TrimSpace(d.(string))

							} else {
								return nil, fmt.Errorf(fmt.Sprintf("key [%s] must be not null", ss_opt))
							}

						} else if ss_v == "purge_on_start" {

							d, err := getValue(tomlTree, ss_opt, "bool")
							if err != nil {
								return nil, fmt.Errorf(fmt.Sprintf("key [%s] : %s", ss_opt, err.Error()))
							}
							if d != nil {
								c.Cache.Local.Purge = d.(bool)
							}

						} else if ss_v == "sweep" {

							d, err := getValue(tomlTree, ss_opt, "int64")
							if err != nil {
								return nil, err
							}

							if d != nil && d.(int64) >= 0 {

								if d.(int64) > 86400 {
									return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] must be lower than or equal 86400", k, v))
								}
								c.Cache.Local.Sweep = uint32(d.(int64))

							} else {
								return nil, fmt.Errorf(fmt.Sprintf("key [%s] must be greater than 0", ss_opt))
							}

						} else {
							return nil, fmt.Errorf(fmt.Sprintf("key [%s] not exist", ss_opt))
						}
					}

				} else {

					return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] not exist", k, v))

				}

			}

		} else if k == "AUTH" {

			// Parser en fonction des clefs définies dans le fichier Toml
			for _, v := range tomlTree.Get(k).(*toml.Tree).Keys() {

				opt := fmt.Sprintf("%s.%s", k, v)

				if v == "mech" {

					d, err := getValue(tomlTree, opt, "[]interface {}")
					if err != nil {
						return nil, err
					}

					if d != nil {
						tab, err := interfaceToStringTab(d.([]interface{}), 1, true, true)
						if err != nil {
							return nil, err
						}
						c.Auth.MechList = tab
					}

				} else if v == "auth_multi" {

					d, err := getValue(tomlTree, opt, "bool")
					if err != nil {
						return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] : %s", k, v, err.Error()))
					}
					if d != nil {
						c.Auth.AuthMulti = d.(bool)
					}

				} else {
					return nil, fmt.Errorf(fmt.Sprintf("key [%s.%s] not exist", k, v))
				}

			}

		} else if k == "PLUGIN" {
			continue

		} else {
			return nil, fmt.Errorf(fmt.Sprintf("section [%s] not exist", k))
		}

	}

	return &c, nil
}

func processConfig(c *Config, appName string) (*Config, error) {

	if c.Cache.KeyRand {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		k := sha512.Sum512([]byte(fmt.Sprintf("%f%s%d", rnd.ExpFloat64(), time.Now(), rand.Uint64())))
		c.Cache.Key = k[:]
	}

	if c.Server.SocketSize < c.Server.BufferSize {
		return nil, fmt.Errorf(fmt.Sprintf("option [SERVER.socket_size]: %d must greater or equal than [SERVER.buffer_size]: %d", c.Server.SocketSize, c.Server.BufferSize))
	}

	// Init de Syslog
	args := map[string]interface{}{
		"length":   10,
		"tag":      appName,
		"logger":   c.Server.LogType,
		"facility": c.Server.LogFacility,
	}

	myLog, err := logme.New(args)
	if err != nil {
		return nil, err
	}
	c.Log = myLog

	return c, nil
}
