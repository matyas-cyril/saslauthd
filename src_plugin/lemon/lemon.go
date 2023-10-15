package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// La fonction de vérification est présente, car obligatoire, mais ne fait rien
func Check(opt map[string]interface{}) (bytes.Buffer, error) {

	var buffer bytes.Buffer

	// convertir l'interface en structure compréhensible par le plugin
	data, err := interfaceToStruct(opt)
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

	var arg Lemon

	dec := gob.NewDecoder(&args)
	if err := dec.Decode(&arg); err != nil && err != io.EOF {
		return false, err
	}

	token := strings.TrimSpace(string(data["d1"]))
	if len(token) == 0 {
		return false, fmt.Errorf("authentication Lemon failed - token length not valid")
	}

	if arg.Timeout > 0 {
		http.DefaultClient.Timeout = time.Duration(arg.Timeout) * time.Second
	}

	//fmt.Println(fmt.Sprintf("%s%s", arg.Url, token))

	resp, err := http.Get(fmt.Sprintf("%s%s", arg.Url, token))
	if err != nil {
		//return false, fmt.Errorf(fmt.Sprintf("authentication Lemon failed - %s", err.Error()))
		return false, fmt.Errorf(fmt.Sprintf("authentication Lemon failed - %s URL: %s", err.Error(), arg.Url))
	}

	body, err := io.ReadAll(resp.Body)

	//fmt.Println(string(body))

	if err != nil {
		return false, fmt.Errorf(fmt.Sprintf("authentication Lemon failed - %s", err.Error()))
	}

	type reqLemon struct {
		Mail   string `json:"mail"`
		Uid    string `json:"uid"`
		Status string `json:"accountStatus"`
	}

	var l reqLemon
	err = json.Unmarshal(body, &l)
	if err != nil {
		return false, fmt.Errorf(fmt.Sprintf("authentication Lemon failed - %s", err.Error()))
	}

	l.Mail = strings.TrimSpace(l.Mail)
	l.Uid = strings.TrimSpace(l.Uid)
	l.Status = strings.TrimSpace(l.Status)

	// Le compte est-il non-actif ?
	if len(l.Status) > 0 && !strings.EqualFold(l.Status, "active") {
		return false, fmt.Errorf(fmt.Sprintf("authentication Lemon failed - uid %s not active", data["0"]))
	}

	// Mode virtualdomain
	if len(data["d3"]) != 0 {

		login := fmt.Sprintf("%s@%s", data["d0"], data["d3"])
		if !strings.EqualFold(l.Mail, login) {
			return false, fmt.Errorf(fmt.Sprintf("authentication Lemon failed - login %s invalid", login))
		}

	} else {

		if !strings.EqualFold(l.Uid, string(data["d0"])) {
			return false, fmt.Errorf(fmt.Sprintf("authentication Lemon failed - login %s invalid", data["d0"]))
		}

	}

	return true, nil
}

type Lemon struct {
	Url     string // URL pour vérifier le token
	Timeout uint16 // Timeout de la requête
	Active  string // Valeur du champ d'un signifiant qu'un compte est actif
}

func interfaceToStruct(data map[string]interface{}) (*Lemon, error) {

	lemon := Lemon{
		Timeout: 5,
		Active:  "active",
	}

	for k := range data {

		if k == "url" {

			if !reflect.ValueOf(string("")).Type().ConvertibleTo(reflect.ValueOf(data[k]).Type()) {
				return nil, fmt.Errorf(fmt.Sprintf("lemon param key %s must be a string", k))
			}

			v := strings.TrimSpace(data[k].(string))

			if len(v) < 8 {
				return nil, fmt.Errorf(fmt.Sprintf("lemon param key %s must be a string with length > 7", k))
			}

			// Ajout du / si absent en fin d'URL
			if v[len(v)-1] != '/' {
				v = fmt.Sprintf("%s/", v)
			}
			lemon.Url = v

		} else if k == "timeout" {

			if !reflect.ValueOf(int(0)).Type().ConvertibleTo(reflect.ValueOf(data[k]).Type()) {
				return nil, fmt.Errorf(fmt.Sprintf("lemon param key %s must be an integer", k))
			}

			v := data[k].(int)

			if v < 0 || v > 3600 {
				return nil, fmt.Errorf(fmt.Sprintf("lemon param key %s must be an integer greater than 0 and lower than 3600", k))
			}

			lemon.Timeout = uint16(v)

		} else if k == "active" {

			if !reflect.ValueOf(string("")).Type().ConvertibleTo(reflect.ValueOf(data[k]).Type()) {
				return nil, fmt.Errorf(fmt.Sprintf("lemon param key %s must be a string", k))
			}

			lemon.Active = strings.TrimSpace(data[k].(string))

		} else {
			return nil, fmt.Errorf(fmt.Sprintf("lemon param key %s not exist", k))
		}

	}
	return &lemon, nil
}
