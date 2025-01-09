package saslauthd

import (
	"bytes"
	"fmt"
	"net"

	myLog "github.com/matyas-cyril/logme"
)

func handleConnection(cnx net.Conn, conf configFile, msgID myLog.MsgID) (_request request) {

	defer func() {

		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> end", msgID))
		}

		// Interception des erreurs Panic
		if err := recover(); err != nil {
			_request.auth = false
			_request.err = fmt.Errorf("auth processing failure[%s]", err)

			if Debug() {
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> panic: Err[%s]", msgID, err))
			}

		}

	}()

	if Debug() {
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> start", msgID))
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> readSocket", msgID))
		conf.Log.Info(msgID, "[DEBUG]: reading socket")
	}

	// Obtenir le contenu de la socket
	raw, err := readSocket(cnx, int(conf.Server.SocketSize), int(conf.Server.BufferSize), int(conf.Server.BufferTimeout), msgID)
	if err != nil {
		conf.Log.Error(msgID, err.Error())
		return request{false, err}
	}

	if Debug() {
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> readSocket -> end", msgID))
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> start", msgID))
		conf.Log.Info(msgID, "[DEBUG]: read socket successfully")
		conf.Log.Info(msgID, "[DEBUG]: extracting socket")
	}

	// Transformer la trame en map[string][][]byte
	data, err := extractData(raw, conf.Server.BufferHashType, msgID)
	if err != nil {
		conf.Log.Error(msgID, err.Error())

		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData: Err[%s]", msgID, err))
		}

		return request{false, err}
	}

	logUser := string(data["d0"])
	if len(data["d3"]) != 0 {
		logUser = fmt.Sprintf("%s@%s", data["d0"], data["d3"])
	}

	if Debug() {
		debug.addLogInFile(fmt.Sprintf("#[%s] -> ..-> handle -> extractData -> end", msgID))
		debug.addLogInFile(fmt.Sprintf("#[%s] -> ..-> handle -> user[%s]", msgID, logUser))
		debug.addLogInFile(fmt.Sprintf("#[%s] -> ..-> handle -> service[%s]", msgID, data["d2"]))
	}

	conf.Log.Info(msgID, fmt.Sprintf("extract socket successfully: login[%s] - service[%s]", logUser, data["d2"]))

	auth_ok := []byte{'1'}
	auth_ko := []byte{'0'}

	// Est-ce en cache
	if conf.Cache.Enable {

		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> cache -> login[%s]", msgID, logUser))
		}

		conf.Log.Info(msgID, fmt.Sprintf("looking for %s auth in cache", logUser))

		// Obtention du nom du fichier de cache
		dataCache, err := cache.GetCache(data["key"])
		if err != nil {
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> getCache[%s]: Err[%s]", msgID, conf.Cache.Category, err))
			}

		} else {

			// Vérification de la validité du cache
			aut, err := getAuthInCached(data, dataCache)
			if err != nil {
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> getAuthInCached[%s]: Err[%s]", msgID, conf.Cache.Category, err))
				}
			} else {

				if Debug() {
					debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> getstatus in cache:[%v] ", msgID, aut))
				}

				if bytes.Equal(aut, auth_ok) {
					conf.Log.Info(msgID, fmt.Sprintf("cached %s return auth success", logUser))
					if Debug() {
						debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> getAuthInCached[%s] -> login[%s] -> auth success", msgID, conf.Cache.Category, logUser))
					}

				} else {
					conf.Log.Info(msgID, fmt.Sprintf("cached %s return auth failed", logUser))
					if Debug() {
						debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> getAuthInCached[%s] -> login[%s] -> auth failed", msgID, conf.Cache.Category, logUser))
					}
				}

				// Retourne la valeur présente dans le cache
				return request{bytes.Equal(aut, auth_ok), nil}
			}

		}

		// Pas présent dans le cache
		conf.Log.Info(msgID, fmt.Sprintf("auth %s not cached", logUser))
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> login:[%v] not cached", msgID, logUser))
		}
	}

	conf.Log.Info(msgID, fmt.Sprintf("auth request for %s ...", logUser))
	if Debug() {
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> auth request -> login:[%v]", msgID, logUser))
	}

	/*
		Traitement de l'authentification
	*/
	auth := auth(data, conf, msgID)
	authRequest := request{}

	if auth { // Succès d'authentification
		conf.Log.Info(msgID, fmt.Sprintf("auth request for %s return auth success", logUser))
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> auth request -> login:[%v] -> auth success", msgID, logUser))
		}

		authRequest.auth = true

	} else { // Echec d'authentification
		conf.Log.Info(msgID, fmt.Sprintf("auth request for %s return auth failed", logUser))
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> auth request -> login:[%v] -> auth failed", msgID, logUser))
		}

		authRequest.auth = false
	}

	// Enregistre les données en cache, si fonctionnement normal (pas d'err Panic)
	if conf.Cache.Enable {

		hashKey := data["key"]

		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> cache -> login[%s] -> save", msgID, logUser))
		}

		// suppression de "key", car le hash est obligatoire dans le cache. Fait doublon
		delete(data, "key")

		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> cache -> delete data['key']", msgID))
		}

		// Insertion du l'état de l'auth
		if auth {
			data["aut"] = auth_ok
		} else {
			data["aut"] = auth_ko
		}

		// Mise en cache
		err := cache.SetSucces(data, hashKey)
		if err != nil {
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> cache -> SetSucces[%s]: Err[%s]", msgID, conf.Cache.Category, err))
			}
			conf.Log.Info(msgID, fmt.Sprintf("caching failure of '%s' : %s", logUser, err))

		} else {

			if auth {
				conf.Log.Info(msgID, fmt.Sprintf("caching success of '%s' authentication success", logUser))
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> cache -> login[%s] -> status:[success]", msgID, logUser))
				}
			} else {
				conf.Log.Info(msgID, fmt.Sprintf("caching success of '%s' authentication failure", logUser))
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> cache -> login[%s] -> status:[failed]", msgID, logUser))
				}
			}

		}

	}

	return authRequest
}

// Obtenir la valeur de l'authentification en cache
func getAuthInCached(data, dataCache map[string][]byte) ([]byte, error) {

	// Vérifie la présence du status d'authentification
	aut := dataCache["aut"]
	if aut == nil {
		return nil, fmt.Errorf("key 'aut' is missing in data cache")
	}

	// Les Hash sont différents....
	if !bytes.Equal(data["key"], dataCache["_key_"]) {
		return nil, fmt.Errorf("value of hash is different between data frame and data cache")
	}

	for k := range data {
		if k == "key" {
			continue
		}
		if !bytes.Equal(data[k], dataCache[k]) {
			return nil, fmt.Errorf("value of key '%s' is different between data frame and data cache", k)

		}
	}

	return aut, nil
}
