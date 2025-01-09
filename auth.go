package saslauthd

import (
	"bytes"
	"fmt"

	myLog "github.com/matyas-cyril/logme"
	"github.com/matyas-cyril/saslauthd/config"
)

// Effectue le ou les requêtes d'auth. data contient les données décodées de la trame, authMethods est la liste
//
// des plugins que l'on va utiliser et plugins contient les conf des plugins
// @return: bool -> true (succès d'auth) / false (échec d'auth). C'est le résultat d'exécution du dernier plugin déclaré dans server
//
//	string -> Nom du dernier plugin exécuté
//	[]error -> contient la liste des erreurs retournées par le / les plugins
//func auth(data map[string][]byte, authMethods []string, plugins map[string]*DefinePlugin, authMulti bool) (int8, string, []error) {

func auth(data map[string][]byte, conf configFile, msgID myLog.MsgID) bool {

	defer func() {
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> end", msgID))
		}
	}()

	if conf.Auth.AuthMulti {
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> multi", msgID))
		}
		return auth_multi(data, conf, msgID)
	}

	if Debug() {
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq", msgID))
	}
	return auth_seq(data, conf, msgID)
}

// @return: bool -> true (succès d'auth) / false (échec d'auth). C'est le résultat d'exécution du dernier plugin
// @retunn: string -> nom du dernier plugin d'auth exécuté
// @return: error -> erreur ne résultant pas d'un échec d'authentification
func auth_seq(data map[string][]byte, conf configFile, msgID myLog.MsgID) bool {

	if len(conf.Auth.MechList) == 0 {
		conf.Log.Info(msgID, "auth method - no plugin available")
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> no plugin available", msgID))
		}

		return false
	}

	// Parcourir tant que le retour n'est pas true
	for _, authMethod := range conf.Auth.MechList {

		if authMethod == "YES" {
			if Debug() {
				conf.Log.Info(msgID, fmt.Sprintf("[DEBUG]: plugin '%s' used", authMethod))
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> plugin[%s]", msgID, authMethod))
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> auth: success", msgID))
			}

			conf.Log.Info(msgID, fmt.Sprintf("plugin '%s' auth success", authMethod))
			return true

		} else if authMethod == "NO" {
			if Debug() {
				conf.Log.Info(msgID, fmt.Sprintf("[DEBUG]:plugin '%s' used", authMethod))
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> plugin[%s]", msgID, authMethod))
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> next", msgID))
			}
			continue

		} else {

			if Debug() {
				conf.Log.Info(msgID, fmt.Sprintf("[DEBUG]:plugin '%s' used", authMethod))
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> plugin[%s]", msgID, authMethod))
			}

			// Obtention des données plugins
			p := conf.Auth.Plugin[authMethod]
			if p.File == nil || p.Opt == nil {
				conf.Log.Info(msgID, fmt.Sprintf("plugin '%s' not exist", authMethod))
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> plugin[%s]: Err[Not Exist]", msgID, authMethod))
					debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> next", msgID))
				}
				continue
			}

			if Debug() {
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> plugin[%s] -> ptr[%p]", msgID, authMethod, p))
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> plugin[%s] -> request", msgID, authMethod))
			}

			// Authentification avec plugin
			rtnAuth, err := authPlugin(data, authMethod, p)
			if err != nil {
				if rtnAuth < 0 {
					// Log de l'erreur Panic

					if Debug() {
						debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> plugin[%s] -> panic: Err[%s]", msgID, authMethod, err))
						debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> plugin[%s] -> next", msgID, authMethod))
					}

					conf.Log.Info(msgID, err.Error())
					continue
				}
				// Ici les errreurs d'auth hors les panics
				conf.Log.Info(msgID, fmt.Sprintf("plugin '%s' auth failure: %s", authMethod, err))

				if Debug() {
					debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> plugin[%s] -> failure: Err[%s]", msgID, authMethod, err))
					debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> plugin[%s] -> next", msgID, authMethod))
				}

			}

			// Succés d'auth
			if rtnAuth > 0 {
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> auth: success", msgID))
				}
				conf.Log.Info(msgID, fmt.Sprintf("plugin '%s' auth success", authMethod))
				return true
			}

		}

	}

	if Debug() {
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> .. -> auth -> seq -> auth: no allowed plugin passed authentication with success", msgID))
	}
	conf.Log.Info(msgID, "no allowed plugin passed authentication with success")
	return false
}

func auth_multi(data map[string][]byte, conf configFile, msgID myLog.MsgID) bool {
	return auth_seq(data, conf, msgID)
}

// authPlugin fait la liaison entre le serveur et les plugins externes
// @return int8: -1 capturer panic erreur, 0 echec d'authentification, 1 succès d'authentification
// @return error: contenu de l'erreur
func authPlugin(data map[string][]byte, name string, plugin *config.DefinePlugin) (_return int8, _err error) {

	defer func() {

		if err := recover(); err != nil {
			_return = -1
			_err = fmt.Errorf("panic error plugin '%s' processing [%s]", name, err)
		}

	}()

	// On récupère la fonction Auth
	symbolAuth, err := plugin.File.Lookup("Auth")
	if err != nil {
		return 0, fmt.Errorf("plugin %s : %s", name, err.Error())
	}

	sFuncInfo := symbolAuth.(func(data map[string][]byte, args bytes.Buffer) (bool, error))
	rslt, err := sFuncInfo(data, *plugin.Opt)
	if err != nil {
		return 0, fmt.Errorf("plugin %s : %s", name, err.Error())
	}

	if rslt {
		return 1, nil
	}

	return 0, nil
}
