package sdnotify

import (
	"fmt"
	"net"
	"os"
)

type NotifyState string

const (
	NotifyReady     NotifyState = "READY=1"
	NotifyStopping  NotifyState = "STOPPING=1"
	NotifyReloading NotifyState = "RELOADING=1"
	NotifyWatchdog  NotifyState = "WATCHDOG=1"
)

// NotifyStatus crée un signal de statut personnalisé
func NotifyStatus(msg string) NotifyState {
	return NotifyState(fmt.Sprintf("STATUS=%s", msg))
}

// SdNotify envoie un message au daemon d'initialisation.
// @param :
// Si unsetEnvVar est true alors la var d'env 'NOTIFY_SOCKET' est supprimée
// state correspond à la notification envoyée. Peut-être perso via NotifyStatus
// @return :
// (true, nil) -> notif envoyée avec succès
// (false, nil) -> NOTIFY_SOCKET n'est pas définie donc rien ne se passe. Si unsetEnvVAr est à true, c'est une "erreur" normale !!!
// (false, err) -> erreur lors de l'envoi de la notification
func SdNotify(unsetEnvVar bool, state NotifyState) (bool, error) {

	socketAddr := &net.UnixAddr{
		Name: os.Getenv("NOTIFY_SOCKET"),
		Net:  "unixgram",
	}

	if socketAddr.Name == "" {
		return false, nil
	}

	// Supprime NOTIFY_SOCKET pour éviter que les processus enfants ne tentent de notifier.
	if unsetEnvVar {
		if err := os.Unsetenv("NOTIFY_SOCKET"); err != nil {
			return false, err
		}
	}

	cnx, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
	if err != nil {
		return false, err
	}
	defer cnx.Close()

	if _, err = cnx.Write([]byte(state)); err != nil {
		return false, err
	}

	return true, nil
}

func SdNotify_Ready() (bool, error) {
	return SdNotify(false, NotifyReady)
}

func SdNotify_Stopping() (bool, error) {
	return SdNotify(true, NotifyStopping)
}

func SdNotify_Reloading() (bool, error) {
	return SdNotify(false, NotifyReloading)
}

func SdNotify_Watchdog() (bool, error) {
	return SdNotify(false, NotifyWatchdog)
}
