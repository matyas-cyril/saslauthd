package saslauthd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	myLog "github.com/matyas-cyril/logme"
)

func readSocket(cnx net.Conn, sizeFrame int, sizeBuffer int, timeout int, msgID myLog.MsgID) (data []byte, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			data = nil
			err = fmt.Errorf("panic error read socket : %s", pErr)
		}
	}()

	buffer := make([]byte, sizeBuffer)

	if Debug() {
		debug.addLogInFile(fmt.Sprintf("#[%s] -> ..-> handle -> readSocket -> Size Buffer: [%d]", msgID, sizeBuffer))
	}

	for {

		n, err := cnx.Read(buffer)
		if err != nil {
			if err.(*net.OpError).Timeout() {
				// Erreur "maitrisée" suite au SetReadDeadline
				break
			}
			// Problème autre que le timeout prévu
			return nil, err
		}

		lenBuffer := len(data) + len(buffer[:n])
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> readSocket -> Read Buffer: [%d]", msgID, n))
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> readSocket -> Total: [%d]", msgID, lenBuffer))
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> readSocket -> Data_Frame_Max: [%d]", msgID, sizeFrame))
		}

		// On vérifie qu'on ne dépasse pas la longueur max autorisée pour la trame
		if lenBuffer > sizeFrame {
			return nil, fmt.Errorf("data frame exceeds the size of %d allowed", sizeFrame)
		}

		data = append(data, buffer[:n]...)
		if n <= sizeBuffer {
			// Mise en place d'un timer si la dernière itération est <=
			// Le cnx.Read suivant est en attente / bloquant
			// A chaque passage dans cette boucle le timer est reinitialisé
			cnx.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> readSocket -> Enable Dead Timer -> Read Buffer[%d] <= Size Buffer[%d] : Timer[%d]", msgID, n, sizeBuffer, timeout))
			}
		}
	}
	return data, nil
}

// extractData convertie un []byte en map[string][]byte.
// Les données retournées de la map sont usr, pwd, srv, dom, key
func extractData(frame []byte, bufferHashType uint8, msgID myLog.MsgID) (data map[string][]byte, err error) {

	data = map[string][]byte{}
	defer func() {
		if pErr := recover(); pErr != nil {
			data = nil
			err = fmt.Errorf("data frame not valid %s", pErr)
		}
	}()

	lenFrame := uint(len(frame))
	if Debug() {
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> len(data_frame): %d", msgID, lenFrame))
	}

	var ptr uint = 0
	var cpt uint = 0
	for {

		// condition de sortie
		if (ptr + 2) >= lenFrame {
			break
		}

		lenData := uint(frame[ptr])*256 + uint(frame[ptr+1])

		data[fmt.Sprintf("d%d", cpt)] = frame[ptr+2 : ptr+lenData+2]
		if Debug() {
			if cpt == 1 {
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> d%d: %v -> string[%s]", msgID, cpt, []byte(strings.Repeat("x", len(data[fmt.Sprintf("d%d", cpt)]))), strings.Repeat("x", len(data[fmt.Sprintf("d%d", cpt)]))))
			} else {
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> d%d: %v -> string[%s]", msgID, cpt, data[fmt.Sprintf("d%d", cpt)], data[fmt.Sprintf("d%d", cpt)]))
			}
		}

		ptr = ptr + lenData + 2
		cpt++
	}

	// Génération de la clef unique
	switch bufferHashType {
	case 0:
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> key -> hash: md5", msgID))
		}
		data["key"] = []byte(fmt.Sprintf("%x", md5.Sum(frame)))
	case 1:
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> key -> hash: sha1", msgID))
		}
		data["key"] = []byte(fmt.Sprintf("%x", sha1.Sum(frame)))
	case 2:
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> key -> hash: sha256", msgID))
		}
		data["key"] = []byte(fmt.Sprintf("%x", sha256.Sum256(frame)))
	case 3:
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> key -> hash: sha512", msgID))
		}
		data["key"] = []byte(fmt.Sprintf("%x", sha512.Sum512(frame)))
	}

	if Debug() {
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> key: %s", msgID, data["key"]))
	}
	return data, err
}

// mkDirForSocket permet de créer les dossiers
// et les sous-dossiers de lasocket
func mkDirForSocket(socket string) error {

	fullPath, err := filepath.Abs(filepath.Dir(socket))
	if err != nil {
		return err
	}
	return os.MkdirAll(fullPath, os.ModePerm)
}

// rmDirForSocket Supprimer les dossiers et sous-dossiers
// si ils sont vides
func rmDirForSocket(socket string) error {

	isEmptyDir := func(path string) (bool, error) {
		entries, err := os.ReadDir(path)
		if err != nil {
			return false, err
		}
		return len(entries) == 0, nil
	}

	fullPath, err := filepath.Abs(filepath.Dir(socket))
	if err != nil {
		return err
	}

	for fullPath != "/" {

		f, err := os.Stat(fullPath)
		if err != nil {
			return err
		}

		if !f.IsDir() {
			return fmt.Errorf("file %s is not a directory", fullPath)
		}

		rst, err := isEmptyDir(fullPath)
		if err != nil {
			return err
		}

		if !rst {
			return fmt.Errorf("directory %s not empty", fullPath)
		}

		if err = os.Remove(fullPath); err != nil {
			return err
		}
		fullPath = filepath.Dir(fullPath)
	}

	return nil
}
