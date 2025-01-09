package main

import (
	"bytes"
	"encoding/gob"
	"io"
)

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

	var arg ldapStruct

	dec := gob.NewDecoder(&args)
	if err := dec.Decode(&arg); err != nil && err != io.EOF {
		return false, err
	}

	return false, nil

}

type ldapStruct struct {
}

func interfaceToStruct(data map[string]interface{}) (*ldapStruct, error) {
	return nil, nil
}
