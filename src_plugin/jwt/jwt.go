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
	"strings"

	myJwt "github.com/cristalhq/jwt/v5"
	"golang.org/x/exp/slices"
)

// La fonction de vérification est présente, car obligatoire, mais ne fait rien
func Check(opt map[string]any) (buffer bytes.Buffer, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			buffer = bytes.Buffer{}
			err = fmt.Errorf("panic error plugin jwt : %s", pErr)
		}
	}()

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

	jwtToken, err := myJwt.ParseNoVerify(data["pwd"])
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
	if err := checkJwt(data["usr"], data["dom"], jwt); err != nil {
		return false, err
	}

	// Obtenir la clef
	key, err := getKey(jwt.Iss, args)
	if err != nil {
		return false, err
	}

	// Vérifier que l'aud est présent
	if !slices.Contains(key.Aud, jwt.Aud) {
		return false, fmt.Errorf("jwt aud '%s' not valid", jwt.Aud)
	}

	// Vérifier la signature du JWT
	if err := checkSignJwt(data["pwd"], key.Pwd, jwtToken.Header().Algorithm.String()); err != nil {
		return false, err
	}

	return true, nil
}

// Obtenir la clef de vérification en fonction de l'ISS
func getKey(iss string, args bytes.Buffer) (*jwtCredent, error) {

	var keyChain map[string]jwtCredent

	dec := gob.NewDecoder(&args)
	if err := dec.Decode(&keyChain); err != nil && err != io.EOF {
		return nil, err
	}

	key := keyChain[iss]
	if len(key.Aud) == 0 {
		return nil, fmt.Errorf("no key available to verify jwt signature")
	}

	return &key, nil
}

func interfaceToData(opt map[string]any) (map[string]jwtCredent, error) {

	data := make(map[string]jwtCredent)

	for k := range opt {

		jwtData := jwtCredent{}

		kBis := strings.TrimSpace(k)

		// On évite les doublons
		if len(data[kBis].Aud) != 0 {
			return nil, fmt.Errorf("duplicate entry for %s key", k)
		}

		var ref map[string]any
		if !reflect.ValueOf(ref).Type().ConvertibleTo(reflect.ValueOf(opt[k]).Type()) {
			return nil, fmt.Errorf("jwt param key %s must be a map[string]any", k)
		}

		for vK, vV := range opt[k].(map[string]any) {

			switch strings.TrimSpace(vK) {

			case "aud":

				// Convertir un []any en []string
				for _, data := range vV.([]string) {
					if !slices.Contains(jwtData.Aud, data) {
						jwtData.Aud = append(jwtData.Aud, data)
					}
				}

			case "pwd":
				if len(jwtData.Pwd) != 0 {
					return nil, fmt.Errorf("jwt param key %s value pwd already initialised", k)
				}
				jwtData.Pwd = []byte(vV.(string))

			case "inc":
				if len(jwtData.Pwd) != 0 {
					return nil, fmt.Errorf("jwt param key %s value inc already initialised", k)
				}

				// Charger le fichier texte
				fByte, err := loadFile(vV.(string))
				if err != nil {
					return nil, err
				}
				jwtData.Pwd = fByte

			default:
				return nil, fmt.Errorf("jwt param key %s not exist", k)
			}

		}
		data[kBis] = jwtData
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
