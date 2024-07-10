package config

import (
	"fmt"
	"os"

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

	decodeToml(v, appPath)

	return nil, nil
}

func decodeToml(toml any, appPath string) (*Config, error) {

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
		case "PLUGIN":

		default:
			return nil, fmt.Errorf(fmt.Sprintf("section [%s] not exist", name))
		}

	}

	return &c, nil

}
