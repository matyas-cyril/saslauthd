package saslauthd

import (
	"sync"

	myCache "github.com/matyas-cyril/saslauthd/cache_generic"
	myConfig "github.com/matyas-cyril/saslauthd/config"
)

type saslResponse string

var (
	APP_NAME   string = ""
	APP_PATH   string = ""
	VERSION    string = ""
	BUILD_TIME string = ""
)

const (
	APP_NAME_DEF   string = "go-saslauthd"
	VERSION_DEF    string = "undef"
	BUILD_TIME_DEF string = "undef"
)

const (
	// SASL SUCCESS RESPONSE
	SASL_SUCCESS saslResponse = "\x00\x02\x4F\x4B\x00"

	// SASL FAIL RESPONSE
	SASL_FAIL saslResponse = "\x00\x02\x4E\x4F\x00"
)

type request struct {
	auth bool
	err  error
}

type logInFile struct {
	file   string
	prefix string
	active bool
}

type endingPrgm struct {
	mu      sync.Mutex
	flag    bool  // True -> Demande d'arrêt du serveur
	timeout uint8 // Durée en sec avant de faire un hard stop
}

var debug *logInFile

// Nombre de clients connectés
var clients *wgSync

var cache *myCache.Cache

// Pour faire du graceful shutdown
var ending = endingPrgm{}

// Importer les structures
type configFile = *myConfig.Config
type DefinePlugin = myConfig.DefinePlugin

// Déclaration des compteurs
var statClientReq = NewCpt()
var statClientOK = NewCpt()
var statClientKO = NewCpt()
var statClientReject = NewCpt()
