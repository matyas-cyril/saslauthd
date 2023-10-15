package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"

	myJwt "github.com/cristalhq/jwt/v4"
)

// La fonction de vérification est présente, car obligatoire, mais ne fait rien
func Check(opt map[string]interface{}) (bytes.Buffer, error) {

	var buffer bytes.Buffer

	// convertir l'interface en structure compréhensible par le plugin
	data, err := interfaceToData(opt)
	if err != nil {
		return bytes.Buffer{}, err
	}

	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(data); err != nil {
		return bytes.Buffer{}, err
	}

	return buffer, nil
}

func Auth(data map[string][]byte, args bytes.Buffer) (bool, error) {

	jwtToken, err := myJwt.ParseNoVerify(data["d1"])
	if err != nil {
		return false, err
	}

	// "typ": "JWT"
	if strings.ToUpper(jwtToken.Header().Type) != "JWT" {
		return false, fmt.Errorf("jwt header type must be JWT")
	}

	// Json -> Structure
	jwt := &jwtStruct{}
	if err := json.Unmarshal(jwtToken.Claims(), jwt); err != nil {
		return false, err
	}

	// Contrôle validité syntaxique et horaire du jwt
	if err := checkJwt(data["d0"], data["d3"], jwt); err != nil {
		return false, err
	}

	// Obtenir la clef
	key, err := getKey(jwt.Iss, args)
	if err != nil {
		return false, err
	}

	// Vérifier la signature du JWT
	if err := checkSignJwt(data["d1"], key, jwtToken.Header().Algorithm.String()); err != nil {
		return false, err
	}

	return true, nil
}

// Obtenir la clef de vérification en fonction de l'ISS
func getKey(iss string, args bytes.Buffer) ([]byte, error) {

	var keyChain map[string][]byte

	dec := gob.NewDecoder(&args)
	if err := dec.Decode(&keyChain); err != nil && err != io.EOF {
		return nil, err
	}

	key := keyChain[iss]
	if len(key) == 0 {
		return nil, fmt.Errorf("no key available to verify jwt signature")
	}

	return key, nil
}

func interfaceToData(opt map[string]interface{}) (map[string][]byte, error) {

	data := make(map[string][]byte)

	for k := range opt {

		kBis := strings.TrimSpace(strings.ToLower(k))

		// On évite les doublons
		if len(data[kBis]) != 0 {
			return nil, fmt.Errorf("duplicate entry for %s key", k)
		}

		if !reflect.ValueOf(string("")).Type().ConvertibleTo(reflect.ValueOf(opt[k]).Type()) {
			return nil, fmt.Errorf(fmt.Sprintf("jwt param key %s must be a string", k))
		}

		value := strings.TrimSpace(opt[k].(string))

		// Si value commence par @include: il s'agit d'un fichier externe
		regExp := regexp.MustCompile(`^(?i)\@include:(.+)$`)
		r := regExp.FindStringSubmatch(value)
		if len(r) != 0 {

			// Erreur pendant la décomposition...
			if len(r) != 2 {
				return nil, fmt.Errorf(fmt.Sprintf("jwt param key %s syntax error", k))
			}

			file := strings.TrimSpace(r[1])
			if len(file) == 0 {
				return nil, fmt.Errorf(fmt.Sprintf("jwt param key %s syntax error", k))
			}

			// Charger le fichier texte
			fByte, err := loadFile(file)
			if err != nil {
				return nil, err
			}

			data[kBis] = fByte
			continue

		}

		// Donc c'est mot de passe et non un fichier externe
		data[kBis] = []byte(value)

	}

	return data, nil
}

func loadFile(file string) ([]byte, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	data := ""

	fd := bufio.NewScanner(f)
	for fd.Scan() {
		txt := fd.Text()
		txt = strings.TrimSpace(txt)
		if len(txt) == 0 {
			continue
		}
		data = fmt.Sprintf("%s\n%s", data, txt)
	}

	if err := fd.Err(); err != nil {
		return nil, err
	}

	return []byte(data), nil
}
