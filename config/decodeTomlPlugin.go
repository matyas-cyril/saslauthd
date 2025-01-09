package config

import (
	"bytes"
	"fmt"
	"plugin"
	"strings"
)

func (c *Config) decodeTomlPlugin(d map[string]any) error {

	for name, opt := range d {

		switch NAME := strings.ToUpper(name); NAME {

		// ByPass des methodes internes
		case "YES", "NO":
			continue

		// Traitement des plugins en sous-section [PLUGIN.*]
		default:
			// Chargement du plugin
			dataPlugin, err := loadPlugin(c.Server.PluginPath, NAME, opt.(map[string]any))
			if err != nil {
				return err
			}
			c.Auth.Plugin[NAME] = dataPlugin
		}

	}

	return nil
}

func loadPlugin(path, name string, opt map[string]any) (_defPlug *DefinePlugin, _err error) {

	defer func() {
		if err := recover(); err != nil {
			_defPlug = nil
			_err = fmt.Errorf("%s", err)
		}
	}()

	// Nom complet du plugin repertoire/nom.sasl
	pluginName := fmt.Sprintf("%s.sasl", strings.ToLower(name))
	fullPathPlugin := fmt.Sprintf("%s/%s", path, pluginName)

	// Charger le fichier .sasl
	plugin, err := plugin.Open(fullPathPlugin)
	if err != nil {
		return nil, fmt.Errorf("plugin %s - %s", pluginName, err)
	}

	// Vérifier la présence de 2 fonctions Check et Auth
	symbolCheck, err := plugin.Lookup("Check")
	if err != nil {
		return nil, fmt.Errorf("plugin '%s' function 'Check' not defined", pluginName)
	}

	// On vérifie juste la présence de la fonction
	_, err = plugin.Lookup("Auth")
	if err != nil {
		return nil, fmt.Errorf("plugin '%s' function 'Auth' not defined", pluginName)
	}

	// L'analyse des paramètres se fait côté plugin
	sFuncCheck := symbolCheck.(func(param map[string]any) (bytes.Buffer, error))
	dataCheck, err := sFuncCheck(opt)
	if err != nil {
		return nil, err
	}

	return &DefinePlugin{
		Path: fullPathPlugin,
		Opt:  &dataCheck,
		File: plugin,
	}, nil
}
