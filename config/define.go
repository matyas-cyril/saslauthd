package config

import (
	"bytes"
	"plugin"

	"github.com/matyas-cyril/logme"
)

type Config struct {
	Server struct {
		Network        string            // type de socket : 'unix' unquement pour l'instant
		Socket         string            // socket d'écoute
		User           string            // Utilisateur d'appartenance de la socket
		Group          string            // Groupe d'appartenance de la socket
		RateInfo       uint16            // Fréquence d'affichage des infos serveurs en seconde
		ClientMax      uint32            // Nombre MAX de clients
		ClientTimeout  uint8             // Durée MAX en seconde d'une transaction client
		BufferSize     uint16            // Taille du buffer de lecture de la socket
		BufferTimeout  uint16            // Définir un timeout à la derniere itération (en ms)
		BufferHashType uint8             // 0->md5, 1->sha1, 2->sha256, 3->sha512
		SocketSize     uint16            // Taille max autorisée pour la socket
		PluginPath     string            // Répertoire par, défaut des plugins
		LogType        logme.LogPrint    // logme.LOGME_TERM | logme.LOGME_SYSLOG | logme.LOGME_BOTH | logme.LOGME_NO
		LogFacility    logme.LogFacility // LOGME_F_AUTH | LOGME_F_MAIL | LOGME_F_SYSLOG | LOGME_F_USER | LOGME_F_LOCAL[0-7]
		Graceful       uint8             // Définir une temporisation en sec lors de l'arrêt du serveur
		Stat           uint              // Définir la fréquence d'affichage des statistiques
	}

	Debug struct {
		Enable bool   // Activer le mode debug
		File   string // Fichier de sortie du mode debug (Full Path)
	}

	Cache struct {
		Enable   bool   // Activer ou désactiver l'utilisation du cache
		Category string // LOCAL, KEYDB, REDIS
		Key      []byte // Clef de chiffrement des données en cache
		KeyRand  bool   // Générer une clef aléatoire
		OK       uint32 // Durée en seconde d'un succés d'auth
		KO       uint32 // Durée en seconde d'un échec d'auth
		Check    uint8  // Temps max de contrôle de présence d'un serveur de cache

		Local struct {
			Path  string // Patch du cache local
			Sweep uint32 // Fréquence en seconde de l'exécution de la purge du cache
			Purge bool   // Purger au démarrage
		}

		MemCache struct {
			Host    string // Host de KEYDB ou REDIS
			Port    uint16 // Port de KEYDB ou REDIS
			DB      uint8  // Numéro de la DB
			Timeout uint16 // Timeout de transaction
		}
	}

	Auth struct {
		MechList  []string // Liste des authentifications utilisées
		AuthMulti bool     // Activation du traitement par 3 des authentifications
		Plugin    map[string]*DefinePlugin
	}

	Log *logme.LogMe
}

type DefinePlugin struct {
	Path string
	Opt  *bytes.Buffer
	File *plugin.Plugin
}
