package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	ldap "saslauthd/plugin/ldap/ldap"
)

func Check(opt map[string]any) (buffer bytes.Buffer, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			buffer = bytes.Buffer{}
			err = fmt.Errorf("panic error plugin ldap : %s", pErr)
		}
	}()

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

func Auth(data map[string][]byte, args bytes.Buffer) (valid bool, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			valid = false
			err = fmt.Errorf("panic error plugin ldap : %s", pErr)
		}
	}()

	var arg ldap.LdapOpt

	dec := gob.NewDecoder(&args)
	if err := dec.Decode(&arg); err != nil && err != io.EOF {
		return false, err
	}

	return false, nil

}

func interfaceToStruct(data map[string]any) (*ldap.Ldap, error) {
	return ldap.New(data)
}
