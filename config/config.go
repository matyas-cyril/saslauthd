package config

import (
	"crypto/sha512"
	"fmt"
	"math/rand"
	"net"
	"os"
	"slices"
	"time"

	"github.com/matyas-cyril/logme"
	toml "github.com/pelletier/go-toml/v2"
)

func LoadConfig(tomFile, appName, appPath string) (*Config, error) {

	// Charger le fichier Toml
	f, err := os.Open(tomFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var v interface{}

	if err := toml.NewDecoder(f).Decode(&v); err != nil {
		return nil, err
	}

	// Traiter les données du Toml en Struture Config
	conf, err := initConfigFromToml(v, appPath)
	if err != nil {
		return nil, err
	}

	// Contrôle des données après conversion
	if err := conf.postProcessConfig(appName); err != nil {
		return nil, err
	}

	return conf, nil
}

func initConfigFromToml(toml any, appPath string) (*Config, error) {

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
	c.Server.Graceful = 5

	c.Debug.File = "/tmp/saslauthd.debug"

	c.Cache.Category = "LOCAL"
	c.Cache.OK = 60
	c.Cache.KO = 60
	c.Cache.Check = 1

	c.Cache.Local.Path = "/tmp"
	c.Cache.Local.Sweep = 60

	c.Cache.MemCache.Host = "127.0.0.1"
	c.Cache.MemCache.Port = 6379
	c.Cache.MemCache.DB = 0
	c.Cache.MemCache.Timeout = 3

	c.Auth.MechList = []string{"NO"}
	c.Auth.Plugin = make(map[string]*DefinePlugin)

	// Si des plugins sont déclarés
	plugins := make(map[string]any)

	for name, k := range toml.(map[string]any) {

		switch name {

		case "SERVER":
			if err := c.decodeTomlServer(k); err != nil {
				return nil, err
			}

		case "DEBUG":
			if err := c.decodeTomlDebug(k); err != nil {
				return nil, err
			}

		case "CACHE":
			if err := c.decodeTomlCache(k); err != nil {
				return nil, err
			}

		case "AUTH":
			if err := c.decodeTomlAuth(k); err != nil {
				return nil, err
			}

		case "PLUGIN":
			var err error
			plugins, err = castAnyToStringAny(k)
			if err != nil {
				return nil, fmt.Errorf("value of section [%s] is not a valid hash option", name)
			}

		default:
			return nil, fmt.Errorf("section [%s] not exist", name)
		}

	}

	// Traitement des plugins
	if err := c.decodeTomlPlugin(plugins); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) postProcessConfig(appName string) error {

	// Vérifier que les plugins utilisés sont existants
	authPlugins := []string{}
	for p := range c.Auth.Plugin {
		authPlugins = append(authPlugins, p)
	}
	for _, p := range c.Auth.MechList {
		switch p {

		case "YES", "NO":
			continue

		default:
			if !slices.Contains(authPlugins, p) {
				return fmt.Errorf("auth mechanism '%s' not available - missing section [PLUGIN.%s]", p, p)
			}

		}
	}

	// Contrôle qu'un serveur de cache est présent
	if c.Cache.Enable {

		if c.Cache.Category == "MEMCACHE" {
			cnx, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.Cache.MemCache.Host, c.Cache.MemCache.Port), time.Duration(c.Cache.Check)*time.Second)
			if err != nil {
				return err
			}
			cnx.Close()
		}

	}

	// Générer une clef de chiffrement aléatoire
	if c.Cache.KeyRand {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		k := sha512.Sum512([]byte(fmt.Sprintf("%f%s%d", rnd.ExpFloat64(), time.Now(), rand.Uint64())))
		c.Cache.Key = k[:]
	}

	if c.Server.SocketSize < c.Server.BufferSize {
		return fmt.Errorf("option [SERVER.socket_size]: %d must greater or equal than [SERVER.buffer_size]: %d", c.Server.SocketSize, c.Server.BufferSize)
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
		return err
	}
	c.Log = myLog

	return nil
}
