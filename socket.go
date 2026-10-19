package saslauthd

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"net"
	"os"
	"path/filepath"
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
func extractData(frame []byte, bufferSize uint16, bufferHashType uint8, msgID myLog.MsgID) (data map[string][]byte, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			data = nil
			err = fmt.Errorf("panic error extract data : %s", pErr)
		}
	}()

	data = make(map[string][]byte)

	lenFrame := uint32(len(frame))

	if Debug() {
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> len(data_frame): %d", msgID, lenFrame))
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> BufferSize: %d", msgID, bufferSize))
	}

	var ptr uint32 = 0
	var opt uint8 = 0

	for {

		// On a défini usr, pwd, srv et dom dans data
		if opt > 3 {
			break
		}

		start := ptr + 2

		if start >= lenFrame {

			// Vérifie l'initialisation de usr, pwd et srv au minimum dans data
			if opt < 3 {
				return nil, fmt.Errorf("data frame invalid - user, password and service must define at least")
			}

			break
		}

		end := ptr + uint32(frame[ptr])*uint32(bufferSize) + uint32(frame[ptr+1]) + 2
		if end > lenFrame { // Contrôle pour éviter de déclencher une panic error
			return nil, fmt.Errorf("buffer_size value must be invalid - data frame analyze out of range [%d/%d]", end, lenFrame)
		}

		chunk := frame[start:end] // Contient la donnée

		ptr = end // Positionnement du pointeur pour le prochain tour

		switch opt {

		case 0: // User
			data["usr"] = bytes.TrimSpace(chunk)
			if len(data["usr"]) == 0 {
				return nil, fmt.Errorf("data frame invalid - user value not defined")
			}
			data["login"] = data["usr"]
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> %s: %v -> string[%s]", msgID, "usr", chunk, chunk))
			}

		case 1: // Password
			data["pwd"] = chunk
			if Debug() {
				valueFrameX := bytes.Repeat([]byte("x"), len(chunk))
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> %s: %v -> string[%s]", msgID, "pwd", valueFrameX, valueFrameX))
			}

		case 2: // Service
			data["srv"] = bytes.TrimSpace(chunk)
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> %s: %v -> string[%s]", msgID, "srv", chunk, chunk))
			}

		case 3: // Domain
			data["dom"] = bytes.TrimSpace(chunk)
			if len(data["dom"]) != 0 {
				data["login"] = fmt.Appendf(nil, "%s@%s", data["usr"], data["dom"])

			}
			if Debug() {
				debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> %s: %v -> string[%s]", msgID, "dom", chunk, chunk))
			}

		}
		opt++
	}

	if Debug() {
		debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> %s: %v -> string[%s]", msgID, "login", data["login"], data["login"]))
	}

	// Génération de la clef unique
	switch bufferHashType {

	case 0:
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> key -> hash: md5", msgID))
		}
		data["key"] = fmt.Appendf(nil, "%x", md5.Sum(frame))

	case 1:
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> key -> hash: sha1", msgID))
		}
		data["key"] = fmt.Appendf(nil, "%x", sha1.Sum(frame))

	case 2:
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> key -> hash: sha256", msgID))
		}
		data["key"] = fmt.Appendf(nil, "%x", sha256.Sum256(frame))

	case 3:
		if Debug() {
			debug.addLogInFile(fmt.Sprintf("#[%s] -> .. -> handle -> extractData -> key -> hash: sha512", msgID))
		}
		data["key"] = fmt.Appendf(nil, "%x", sha512.Sum512(frame))

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
