package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// La fonction de vérification est présente, car obligatoire, mais ne fait rien
func Check(opt map[string]any) (buffer bytes.Buffer, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			buffer = bytes.Buffer{}
			err = fmt.Errorf("panic error plugin lemon : %s", pErr)
		}
	}()

	// convertir l'interface en structure compréhensible par le plugin
	data, err := interfaceToStruct(opt)
	if err != nil {
		return bytes.Buffer{}, err
	}

	// Vérifier que l'URL est définie.
	if len(data.Url) == 0 {
		return bytes.Buffer{}, fmt.Errorf("failed to init plugin lemon - url not defined")
	}

	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(data); err != nil {
		return bytes.Buffer{}, err
	}

	return buffer, nil
}

func Auth(data map[string][]byte, args bytes.Buffer) (valid bool, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			valid = false
			err = fmt.Errorf("panic error plugin lemon : %s", pErr)
		}
	}()

	var arg Lemon

	dec := gob.NewDecoder(&args)
	if err := dec.Decode(&arg); err != nil && err != io.EOF {
		return false, err
	}

	token := strings.TrimSpace(string(data["pwd"]))
	if len(token) == 0 {
		return false, fmt.Errorf("authentication Lemon failed - token length not valid")
	}

	if arg.Timeout > 0 {
		http.DefaultClient.Timeout = time.Duration(arg.Timeout) * time.Second
	}

	resp, err := http.Get(fmt.Sprintf("%s%s", arg.Url, token))
	if err != nil {
		return false, fmt.Errorf("authentication Lemon failed - %s URL: %s", err.Error(), arg.Url)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return false, fmt.Errorf("authentication Lemon failed - %s", err.Error())
	}

	type reqLemon struct {
		Mail   string `json:"mail"`
		Uid    string `json:"uid"`
		Status string `json:"accountStatus"`
	}

	var l reqLemon
	err = json.Unmarshal(body, &l)
	if err != nil {
		return false, fmt.Errorf("authentication Lemon failed - %s", err.Error())
	}

	l.Mail = strings.TrimSpace(l.Mail)
	l.Uid = strings.TrimSpace(l.Uid)
	l.Status = strings.TrimSpace(l.Status)

	// Le compte est-il non-actif ?
	if len(l.Status) > 0 && !strings.EqualFold(l.Status, "active") {
		return false, fmt.Errorf("authentication Lemon failed - uid %s not active", data["usr"])
	}

	var login string
	if arg.VirtDom {
		login = string(data["login"])
	} else {
		login = string(data["usr"])
	}

	var authKey string
	if strings.EqualFold(arg.AuthKey, "MAIL") {
		authKey = l.Mail
	} else {
		authKey = l.Uid
	}

	if !strings.EqualFold(authKey, login) {
		return false, fmt.Errorf("authentication Lemon failed - login %s invalid", login)
	}

	return true, nil
}

type Lemon struct {
	Url     string // URL pour vérifier le token
	Timeout uint16 // Timeout de la requête
	Active  string // Valeur du champ d'un signifiant qu'un compte est actif
	AuthKey string // Clef servant pour l'authentification [uid | mail]
	VirtDom bool   // Utilisation des virtdom
}

func interfaceToStruct(data map[string]any) (*Lemon, error) {

	lemon := Lemon{
		Timeout: 5,
		Active:  "active",
		AuthKey: "MAIL",
		VirtDom: true,
	}

	for k := range data {

		switch k {

		case "url":

			v, cast := data[k].(string)
			if !cast {
				return nil, fmt.Errorf("lemon param key %s must be a string", k)
			}

			v = strings.TrimSpace(v)

			if len(v) < 8 {
				return nil, fmt.Errorf("lemon param key %s must be a string with length > 7", k)
			}

			// Ajout du / si absent en fin d'URL
			if v[len(v)-1] != '/' {
				v = fmt.Sprintf("%s/", v)
			}
			lemon.Url = v

		case "timeout":

			v, cast := data[k].(int32)
			if !cast {
				return nil, fmt.Errorf("lemon param key %s must be an integer", k)
			}

			if v < 0 || v > 3600 {
				return nil, fmt.Errorf("lemon param key %s must be an integer greater than 0 and lower than 3600", k)
			}

			lemon.Timeout = uint16(v)

		case "active":

			v, cast := data[k].(string)
			if !cast {
				return nil, fmt.Errorf("lemon param key %s must be a string", k)
			}

			lemon.Active = strings.TrimSpace(v)

		case "authkey":

			v, cast := data[k].(string)
			if !cast {
				return nil, fmt.Errorf("lemon param key %s must be a string", k)
			}
			v = strings.TrimSpace(v)

			switch kV := strings.ToUpper(v); kV {

			case "MAIL", "UID":
				lemon.AuthKey = kV

			default:
				return nil, fmt.Errorf("lemon param key %s must be [mail|uid]", k)
			}

			lemon.AuthKey = v

		case "virtdom":
			v, cast := data[k].(bool)
			if !cast {
				return nil, fmt.Errorf("lemon param key %s must be a boolean", k)
			}
			lemon.VirtDom = v

		default:
			return nil, fmt.Errorf("lemon param key %s not exist", k)
		}

	}
	return &lemon, nil
}
