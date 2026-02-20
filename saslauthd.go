package saslauthd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	myLog "github.com/matyas-cyril/logme"
	myCache "github.com/matyas-cyril/saslauthd/cache"
	myConfig "github.com/matyas-cyril/saslauthd/config"
	mySdnotify "github.com/matyas-cyril/saslauthd/sdnotify"
)

func Check(confFile string) error {

	_, err := myConfig.LoadConfig(confFile, APP_NAME, APP_PATH)
	if err != nil {
		return err
	}
	return nil
}

func varEnv(env, defaut string) string {

	e := strings.TrimSpace(env)
	if len(e) > 0 && len(e) < 32 {
		return e
	}

	return defaut
}

func Start(confFile, appPath string) {

	APP_NAME = varEnv(APP_NAME, APP_NAME_DEF)
	VERSION = varEnv(VERSION, VERSION_DEF)
	BUILD_TIME = varEnv(BUILD_TIME, BUILD_TIME_DEF)
	APP_PATH = appPath

	conf, err := myConfig.LoadConfig(confFile, APP_NAME, APP_PATH)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	func() {
		defer ending.mu.Unlock()
		ending.mu.Lock()
		ending.timeout = conf.Server.Graceful
	}()

	conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("%s %s starting up...", APP_NAME, VERSION))

	if conf.Debug.Enable {
		d, err := newLogInFile(conf.Debug.File)
		if err != nil {
			txtErr := fmt.Sprintf("failed to init debug file : '%s'\n", err)
			fmt.Fprintln(os.Stderr, txtErr)
			os.Exit(2)
		}
		debug = d
		debug.addLogInFile("***** mode debug enabled *****")
		conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("***** Debug mode activated: '%s'", conf.Debug.File))

		appverB := fmt.Sprintf("# -> APP_NAME[%s] VERSION[%s] BUILD_TIME[%s]", APP_NAME, VERSION, BUILD_TIME)
		appPath := fmt.Sprintf("# -> APP_PATH[%s]", APP_PATH)
		debug.addLogInFile(appverB)
		debug.addLogInFile(appPath)
		conf.Log.Info(myLog.MSGID_EMPTY, appverB)
		conf.Log.Info(myLog.MSGID_EMPTY, appPath)
	}

	// Pas de mechanisme d'auth activé
	if len(conf.Auth.MechList) == 0 {
		if Debug() {
			debug.addLogInFile("# -> go -> Server -> No Auth Mechanism")
		}
		conf.Log.Info(myLog.MSGID_EMPTY, "no auth mechanism defined")
		os.Exit(1)
	}

	// Suppression de la socket au démarrage
	if conf.Server.SelfRuling {
		if err := deleteSocket(conf.Server.Socket); err != nil {

			if Debug() {
				debug.addLogInFile(fmt.Sprintf("# -> del socket: Err[%s]", err))
			}

			conf.Log.Info(myLog.MSGID_EMPTY, err.Error())
		}
	}

	if Debug() {
		tmpConf := *conf
		// Cacher la clef
		tmpConf.Cache.Key = []byte(strings.Repeat("x", len(tmpConf.Cache.Key)))
		debug.addLogInFile(fmt.Sprintf("# -> Config File('%s') : %+v\n", confFile, tmpConf))
	}

	// Init du cache ?
	if conf.Cache.Enable {

		opt := make(map[string]any)

		// Préparer les options en fonction du type de cache
		switch conf.Cache.Category {
		case "LOCAL":
			opt["path"] = conf.Cache.Local.Path

		case "MEMCACHE":
			opt["host"] = conf.Cache.ExternalCache.Host
			opt["port"] = conf.Cache.ExternalCache.Port
			opt["timeout"] = conf.Cache.ExternalCache.Timeout

		case "REDIS":
			opt["host"] = conf.Cache.ExternalCache.Host
			opt["port"] = conf.Cache.ExternalCache.Port
			opt["timeout"] = conf.Cache.ExternalCache.Timeout
			opt["db"] = conf.Cache.ExternalCache.DB

		}

		cache, err = myCache.New(conf.Cache.Category, conf.Cache.Key, conf.Cache.OK, conf.Cache.KO, opt)
		// Si echec on désactive, ce n'est pas bloquant
		if err != nil {
			conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("failed to enable cache: '%s'", err))
			conf.Log.Info(myLog.MSGID_EMPTY, "warning cache option is now disabled !!!")

			if Debug() {
				debug.addLogInFile(fmt.Sprintf("# -> failed to enable cache: '%s'", err))
				debug.addLogInFile("# -> warning cache option is now disabled !!!")
			}

			conf.Cache.Enable = false
		}
	}

	// Déclaration du compteur pour connaître le nbr de clients
	clients = NewSync()

	sigterm := make(chan os.Signal, 1)
	exitChan := make(chan int, 1)
	signal.Notify(sigterm, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {

		for {

			sig := <-sigterm
			switch sig {

			// Effectuer une purge
			// kill -SIGUSR1 XXXX [XXXX - PID]
			case syscall.SIGUSR1:
				err := Check(confFile)
				if err != nil {
					conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("check - config file '%s' verification failed : %s", confFile, err.Error()))
				} else {
					conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("check - config file '%s' verification successfully", confFile))
				}

			// Effectuer une nettoyage
			// kill -SIGUSR2 XXXX [XXXX - PID]
			case syscall.SIGUSR2:
				fmt.Println("syscall.SIGUSR2")
				// conf.Log.Info(logme.MSGID_EMPTY, "Reload")
				//exit_chan <- 3

			// kill -SIGINT XXXX [XXXX - PID] - CTRL D - STOP - QUIT
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:

				func() {
					ending.mu.Lock()
					defer ending.mu.Unlock()
					ending.flag = true
				}()

				conf.Log.Info(myLog.MSGID_EMPTY, "SIGTERM signal received")

				// Notification que l'on va arrêter le serveur
				if conf.Server.Notify {
					notif, err := mySdnotify.SdNotify_Stopping()
					if err != nil {
						conf.Log.Error(myLog.MSGID_EMPTY, fmt.Sprintf("failed to send notify STOP : %s", err))
						if Debug() {
							debug.addLogInFile(fmt.Sprintf("# -> go -> Send notify STOP - Failed : %s", err))
						}
					}
					if Debug() {
						if !notif {
							conf.Log.Warning(myLog.MSGID_EMPTY, "failed to send notify STOP")
							debug.addLogInFile("# -> go -> Send notify STOP - Failed")
						} else {
							conf.Log.Info(myLog.MSGID_EMPTY, "succes to send notify STOP")
							debug.addLogInFile("# -> go -> Send notify STOP - OK")
						}
					}
				}

				if clients.Get() > 0 {

					if Debug() {
						debug.addLogInFile("#[SIGTERM] -> signal received", "#[SIGTERM] -> Graceful shutdown")
					}

					conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("graceful shutdown for %d client(s)", clients.Get()))

					maxEnding := time.Now().Add(time.Duration(ending.timeout) * time.Second)
					for time.Now().Compare(maxEnding) < 0 {
						time.Sleep(50 * time.Millisecond)
						if clients.Get() < 1 {
							break
						}
					}
				}

				// Il reste des clients, c'est dommage !!!
				if clients.Get() > 0 {
					if Debug() {
						debug.addLogInFile("#[SIGTERM] -> signal received", "#[SIGTERM] -> Hard shutdown")
					}
					conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("hard shutdown for %d client(s)", clients.Get()))
				}

				if Debug() {
					debug.addLogInFile("#[SIGTERM] -> signal received", "#[SIGTERM] -> Server stopped", "#[SIGTERM] -> Exit: 0")
				}

				conf.Log.Info(myLog.MSGID_EMPTY, "Bye, Bye - Server stopped !!!")

				exitChan <- 0
				return

			default:
				conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("signal '%v' received", sig))
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("#[%v] -> signal received", sig))
				}
				exitChan <- 2
				return
			}
		}

	}()

	// Statistique clients
	go func() {

		// ByPass les statistiques
		if conf.Server.Stat == 0 {
			return
		}

		for {

			time.Sleep(time.Duration(conf.Server.Stat) * time.Second)

			conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("stat (last %d seconds): incoming request(s)[%d]", conf.Server.Stat, statClientReq.Get()))
			conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("stat (last %d seconds): auth succes[%d] - auth failure[%d] - cnx rejected[%d]", conf.Server.Stat, statClientOK.Get(), statClientKO.Get(), statClientReject.Get()))

			// Remise à zéro des compteurs
			statClientReq.Reset()
			statClientOK.Reset()
			statClientKO.Reset()
			statClientReject.Reset()
		}

	}()

	// Watchdog
	if conf.Server.Notify {
		go func() {

			time.Sleep(1 * time.Second)

			for {

				time.Sleep(1 * time.Second)

				notif, err := mySdnotify.SdNotify_Watchdog()
				if err != nil {
					conf.Log.Error(myLog.MSGID_EMPTY, fmt.Sprintf("failed to send Watchdog : %s", err))
					if Debug() {
						debug.addLogInFile(fmt.Sprintf("# -> go -> Send Watchdog - Failed : %s", err))
					}
				}
				if Debug() {
					if !notif {
						conf.Log.Warning(myLog.MSGID_EMPTY, "failed to send Watchdog")
						debug.addLogInFile("# -> go -> Send Watchdog - Failed")

					} else {
						conf.Log.Info(myLog.MSGID_EMPTY, "succes to send Watchdog")
						debug.addLogInFile("# -> go -> Send Watchdog - OK")
					}
				}
			}

		}()
	}

	go func() {

		if Debug() {
			debug.addLogInFile(fmt.Sprintf("# -> go -> Server: Network[%s] - Socket[%s]", conf.Server.Network, conf.Server.Socket))
		}

		// Créér les dossiers pour la socket
		if conf.Server.SelfRuling {
			if err := mkDirForSocket(conf.Server.Socket); err != nil {
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("# -> go -> Server -> Make Socket Folder: Err[%s]", err))
				}
				conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("failed to create folders for socket - %s", err.Error()))
				exitChan <- 3
				return
			}
		}

		// Création de la socket
		ln, err := net.Listen(conf.Server.Network, conf.Server.Socket)
		if err != nil {
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("# -> go -> Server -> Listen: Err[%s]", err))
			}
			conf.Log.Info(myLog.MSGID_EMPTY, err.Error())
			exitChan <- 3
			return
		}

		// Déterminer l'ID de l'utilisateur de la socket
		usr, err := user.Lookup(conf.Server.User)
		if err != nil || usr == nil {
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("# -> go -> Server -> lookup user id[%s]: Err[%s]", conf.Server.User, err))
			}
			conf.Log.Info(myLog.MSGID_EMPTY, err.Error())
			exitChan <- 3
			return
		}

		usrId, err := strconv.Atoi(usr.Uid)
		if err != nil {
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("# -> go -> Server -> lookup user id[%s]: Err[%s]", conf.Server.User, err))
			}
			conf.Log.Info(myLog.MSGID_EMPTY, err.Error())
			exitChan <- 3
			return
		}

		// Déterminer le groupe de l'utilisateur de la socket
		grp, err := user.LookupGroup(conf.Server.Group)
		if err != nil || grp == nil {
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("# -> go -> Server -> lookup group id[%s]: Err[%s]", conf.Server.Group, err))
			}
			conf.Log.Info(myLog.MSGID_EMPTY, err.Error())
			exitChan <- 3
			return
		}

		grpId, err := strconv.Atoi(grp.Gid)
		if err != nil {
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("# -> go -> Server -> lookup group id[%s]: Err[%s]", conf.Server.Group, err))
			}
			conf.Log.Info(myLog.MSGID_EMPTY, err.Error())
			exitChan <- 3
			return
		}

		// Changer les droits de la socket
		err = os.Chown(conf.Server.Socket, usrId, grpId)
		if err != nil {
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("# -> go -> Server -> chown[%s:%s]: Err[%s]", conf.Server.User, conf.Server.Group, err))
			}
			conf.Log.Info(myLog.MSGID_EMPTY, err.Error())
			exitChan <- 3
			return
		}

		// Changer les droits de la socket
		err = os.Chmod(conf.Server.Socket, conf.Server.UGO)
		if err != nil {
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("# -> go -> Server -> chmod[%v]: Err[%s]", conf.Server.UGO, err))
			}
			conf.Log.Info(myLog.MSGID_EMPTY, err.Error())
			exitChan <- 3
			return
		}

		defer func() {

			if err := recover(); err != nil {
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("# -> go -> Server -> Defer: Err[%s]", err))
				}
				conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("failed to start server. check socket '%s' access", conf.Server.Socket))
			}

			if ln != nil {
				if err := ln.Close(); err != nil {
					conf.Log.Info(myLog.MSGID_EMPTY, err.Error())
					if Debug() {
						debug.addLogInFile(fmt.Sprintf("# -> go -> Server -> Close: Err[%s]", err))
					}

				}
			}

		}()

		conf.Log.Info(myLog.MSGID_EMPTY, "* Server started - Let's Go !!!")
		conf.Log.EnableMessageID()

		if Debug() {
			debug.addLogInFile("# -> go -> Server started")
		}

		// Notifier à l'OS que le serveur est READY
		if conf.Server.Notify {
			notif, err := mySdnotify.SdNotify_Ready()
			if err != nil {
				conf.Log.Error(myLog.MSGID_EMPTY, fmt.Sprintf("failed to send notify : %s", err))
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("# -> go -> Send notify READY - Failed : %s", err))
				}
			}
			if Debug() {
				if !notif {
					conf.Log.Warning(myLog.MSGID_EMPTY, "failed to send notify READY")
					debug.addLogInFile("# -> go -> Send notify READY - Failed")
				} else {
					conf.Log.Info(myLog.MSGID_EMPTY, "succes to send notify READY")
					debug.addLogInFile("# -> go -> Send notify READY - OK")
				}
			}
		}

		for {

			// Génération d'un ID unique
			msgID := conf.Log.MessageID()

			cnx, err := ln.Accept()
			if err != nil {
				conf.Log.Info(msgID, err.Error())
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("#[%s] -> go -> Server -> listener failed: Err[%s]", msgID, err))
				}
				//break
				continue
			}

			// On n'accepte plus de cnx
			if ending.flag {
				cnx.Close()
				conf.Log.Info(msgID, "connection refused - shutdown in progress")
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("#[%s] -> go -> Server -> cnx refused -> shutdown in progress", msgID))
				}
				continue
			}

			if conf.Server.Stat > 0 {
				statClientReq.Inc()
			}

			// Nbr max des clients
			if uint32(clients.Get()) > conf.Server.ClientMax {
				conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("max clients[%d] reached", conf.Server.ClientMax))

				if conf.Server.Stat > 0 {
					statClientReject.Inc()
				}

				if Debug() {
					debug.addLogInFile(fmt.Sprintf("#[%s] -> go -> Server -> max clients[%d] reached", msgID, conf.Server.ClientMax))
					debug.addLogInFile(fmt.Sprintf("#[%s] -> go -> Server -> clients: nbr[%d]", msgID, clients.Get()))
					debug.addLogInFile(fmt.Sprintf("#[%s] -> go -> Server -> clients: max[%d]", msgID, conf.Server.ClientMax))
					debug.addLogInFile(fmt.Sprintf("#[%s] -> go -> Server -> cnx close", msgID))
				}

				cnx.Close()
				continue
			}

			clients.Add(1)

			if Debug() {
				debug.addLogInFile(fmt.Sprintf("#[%s] -> go -> Server -> client cnx", msgID))
			}

			go requestClient(cnx, conf, msgID)

		}

	}()

	// Utilisation du cache local - Activation du nettoyage
	if conf.Cache.Enable && conf.Cache.Category == "LOCAL" && conf.Cache.Local.Sweep > 0 {

		go func() {

			// Purge du cache au démarrage
			err := cache.Flush()
			if err != nil {
				conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("cache purge failure: %s", err))
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("# -> go -> cache -> purge failure: %s", err))
				}
			}

			for {

				time.Sleep(time.Duration(conf.Cache.Local.Sweep) * time.Second)

				c_ok, c_ko, c_err, err := cache.Clean()
				if err != nil {
					conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("cache cleaning failure: %s", err))
					if Debug() {
						debug.addLogInFile(fmt.Sprintf("# -> go -> cache -> cleaning failure: %s", err))
					}
					continue
				}
				conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("cache cleaning: OK[%d] - KO[%d]", c_ok, c_ko))
				if Debug() {
					debug.addLogInFile(fmt.Sprintf("# -> go -> cache -> cleaning: OK[%d] - KO[%d]", c_ok, c_ko))
				}
				for _, c := range c_err {
					conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("cache cleaning error: %s", c))
					if Debug() {
						debug.addLogInFile(fmt.Sprintf("# -> go -> cache -> cleaning eror: %s", c))
					}
				}

			}

		}()

	}

	// Afficher dans les logs les informations du serveur
	if conf.Server.RateInfo > 0 {

		go func() {

			for {

				time.Sleep(time.Duration(conf.Server.RateInfo) * time.Second)

				conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("current total client(s): %d / %d", clients.Get(), conf.Server.ClientMax))
				if Debug() {
					conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("nbr goroutine(s): %d", runtime.NumGoroutine()))
					conf.Log.Info(myLog.MSGID_EMPTY, fmt.Sprintf("nbr CPU: %d", runtime.NumCPU()))
				}

			}

		}()

	}

	exitCode := <-exitChan

	// Suppression de la socket
	if err := deleteSocket(conf.Server.Socket); err != nil {
		conf.Log.Info(myLog.MSGID_EMPTY, err.Error())
	}

	// supprimer les répertoires de la socket si vide
	rmDirForSocket(conf.Server.Socket)

	os.Exit(exitCode)
}
