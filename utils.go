package saslauthd

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// deleteSocket: supprimer la socket si elle existe.
// On vérifie que le paramètre 'socket' et bien un
// fichier de type socket avant d'effectuer la suppression
// @return:
// - nil: pas de fichier ou suppression bien effectuée
// - error: problème pendant le traitement
func deleteSocket(socket string) error {

	// Voir si un fichier ou rep existe
	file, err := os.Stat(socket)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to access socket '%s' : ", err.Error())
	}

	// Le fichier est une socket existante
	if file.Mode().Type() != fs.ModeSocket {
		return fmt.Errorf("file '%s' is not a unix socket", socket)
	}

	// Suppression de la socket unix
	if err := os.Remove(socket); err != nil {
		return fmt.Errorf("failed to delete socket '%s' : ", err.Error())
	}

	return nil
}

func newLogInFile(logFile string) (*logInFile, error) {

	logFile = strings.TrimSpace(logFile)

	l := logInFile{}

	path := filepath.Dir(logFile)

	p, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	if !p.IsDir() {
		return nil, fmt.Errorf("debug path '%s' is not a directory : ", logFile)
	}

	f, err := os.Stat(logFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("debug file '%s' error : %s", logFile, err)
	}

	if errors.Is(err, os.ErrNotExist) {
		l.active = true
		l.file = logFile
		return &l, nil
	}

	if !f.IsDir() {
		l.active = true
		l.file = logFile
		return &l, nil
	}

	return nil, fmt.Errorf("debug file '%s' failed to init", logFile)
}

func (l logInFile) addLogInFile(dataLog ...string) error {

	if !l.active {
		return nil
	}

	file, err := os.OpenFile(l.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	logger := log.New(file, l.prefix, log.LstdFlags)
	logger.SetFlags(log.Ldate | log.Lmicroseconds)

	for _, l := range dataLog {
		logger.Println(l)
	}
	return nil
}

func Debug() bool {
	return debug != nil && debug.active
}
