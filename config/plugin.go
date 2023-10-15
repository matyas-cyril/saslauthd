package config

import (
	"bytes"
	"fmt"
	"plugin"
	"strings"

	toml "github.com/pelletier/go-toml"
)

func processPlugin(tomlTree *toml.Tree, c *Config) (_Config *Config, _error error) {

	_plugin := ""

	defer func() {
		if err := recover(); err != nil {
			_Config = nil
			if len(_plugin) > 0 {
				_error = fmt.Errorf("plugin '%s' processing panic[%s]", _plugin, err)
			} else {
				_error = fmt.Errorf("plugin processing panic[%s]", err)
			}
		}
	}()

	// Si rien plugin par défault est "NO"
	if len(c.Auth.MechList) == 0 {
		c.Auth.MechList = []string{"NO"}
		return c, nil
	}

	// Chargement de la configuration des plugins présent dans auth_method
	for _, k := range c.Auth.MechList {

		if len(k) == 0 {
			return nil, fmt.Errorf("empty string in section [AUTH.mech]")
		}

		_plugin = k

		// ByPass des methodes internes
		if k == "YES" || k == "NO" {
			continue
		}

		// Traitement des plugins en sous-section [PLUGIN.*]
		opt := make(map[string]interface{})
		plugin := tomlTree.Get(fmt.Sprintf("PLUGIN.%s", k))
		if plugin != nil {
			opt = plugin.(*toml.Tree).ToMap()
		}

		// Chargement du plugin
		dataPlugin, err := loadPlugin(c.Server.PluginPath, k, opt)
		if err != nil {
			return nil, err
		}

		// Ajout du plugin
		c.Auth.Plugin[k] = dataPlugin

	}

	return c, nil
}

func loadPlugin(path, name string, opt map[string]interface{}) (_defPlug *DefinePlugin, _err error) {

	defer func() {
		if err := recover(); err != nil {
			_defPlug = nil
			_err = fmt.Errorf(fmt.Sprintf("%s", err))
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
	sFuncCheck := symbolCheck.(func(param map[string]interface{}) (bytes.Buffer, error))
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
