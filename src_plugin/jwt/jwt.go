package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	myJwt "github.com/cristalhq/jwt/v5"
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

func interfaceToData(opt map[string]any) (data map[string]jwtCredent, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			data = nil
			err = fmt.Errorf("panic error plugin jwt : %s", pErr)
		}
	}()

	data = make(map[string]jwtCredent)

	for k := range opt {

		if len(k) != len(strings.TrimSpace(k)) {
			return nil, fmt.Errorf("left or right space for %s key", k)
		}

		// On évite les doublons
		if len(data[k].Aud) != 0 {
			return nil, fmt.Errorf("duplicate entry for %s key", k)
		}

		// Initialisation et valeur par défaut
		jwtData := jwtCredent{
			VirtDom: true,
		}

		for vK, vV := range opt[k].(map[string]any) {

			switch vK {

			case "virtdom":
				value, ok := vV.(bool)
				if !ok {
					return nil, fmt.Errorf("jwt param key %s.%s must be a boolean", k, vK)
				}
				jwtData.VirtDom = value

			case "aud":

				tab, ok := vV.([]string)
				if !ok {
					return nil, fmt.Errorf("jwt param key %s.%s must be a []string", k, vK)
				}

				// Convertir un []any en []string
				for _, data := range tab {
					data = strings.TrimSpace(data)
					if !slices.Contains(jwtData.Aud, data) {
						jwtData.Aud = append(jwtData.Aud, data)
					} else {
						// L'audience a déjà été déclarée - doublon
						return nil, fmt.Errorf("jwt param key %s.%s value '%s' for password already initialised", k, vK, data)
					}
				}

			case "pwd":
				if len(jwtData.Pwd) != 0 { // Vérifier que l'on utilise un mot de passe OU un fichier externe
					return nil, fmt.Errorf("jwt param key %s.%s value for password already initialised", k, vK)
				}
				value, ok := vV.(string)
				if !ok {
					return nil, fmt.Errorf("jwt param key %s.%s must be a string", k, vK)
				}
				jwtData.Pwd = []byte(value)

			case "inc":
				if len(jwtData.Pwd) != 0 { // Vérifier que l'on utilise un mot de passe OU un fichier externe
					return nil, fmt.Errorf("jwt param key %s.%s value for password already initialised", k, vK)
				}
				value, ok := vV.(string)
				if !ok {
					return nil, fmt.Errorf("jwt param key %s.%s must be a string", k, vK)
				}

				// Charger le fichier texte
				fByte, err := loadFile(value)
				if err != nil {
					return nil, err
				}
				jwtData.Pwd = fByte

			default:
				return nil, fmt.Errorf("jwt param key %s.%s not exist", k, vK)
			}

		}
		data[k] = jwtData
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
