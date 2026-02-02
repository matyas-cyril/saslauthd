package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"reflect"
	"strings"

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

func interfaceToStruct(data map[string]any) (*ldap.LdapOpt, error) {

	// Déclarations et valeurs par défaut
	var (
		uri           string
		admin         string
		passwd        string
		baseDN        string
		port          uint16 = 389
		timeout       uint16 = 10
		tls           bool   = false
		tlsSkipVerify bool   = true
	)

	for k, v := range data {

		switch k {

		case "uri", "admin", "pwd", "baseDN":
			kV, kErr := v.(string)
			if !kErr {
				return nil, fmt.Errorf("ldap key '%s' failed to typecast", k)
			}
			kV = strings.TrimSpace(kV)
			if len(kV) == 0 {
				return nil, fmt.Errorf("ldap key '%s' defined but empty", k)
			}

			switch k {

			case "uri":
				uri = kV
			case "admin":
				admin = kV
			case "pwd":
				passwd = kV
			case "baseDN":
				baseDN = kV
			}

		case "port", "timeout":
			typeTarget := reflect.TypeFor[int]()
			rv := reflect.ValueOf(data[k])
			if !rv.Type().AssignableTo(typeTarget) {
				return nil, fmt.Errorf("ldap key '%s' failed to typecast", k)
			}

			nbr := rv.Convert(typeTarget).Int()
			if nbr < 0 || nbr > 65535 {
				return nil, fmt.Errorf("ldap key '%s' integer range invalid", k)
			}

			switch k {
			case "port":
				port = uint16(nbr)
			case "timeout":
				timeout = uint16(nbr)
			}

		case "tls", "tlsSkipVerify":
			kV, kErr := v.(bool)
			if !kErr {
				return nil, fmt.Errorf("ldap key '%s' failed to typecast", k)
			}

			switch k {
			case "tls":
				tls = kV
			case "tlsSkipVerify":
				tlsSkipVerify = kV
			}

		default:
			return nil, fmt.Errorf("ldap key '%s' not exist", k)

		}

	}

	return &ldap.LdapOpt{
		Uri:                uri,
		Port:               port,
		Admin:              admin,
		Passwd:             passwd,
		BaseDn:             baseDN,
		Timeout:            timeout,
		Tls:                tls,
		InsecureSkipVerify: tlsSkipVerify,
	}, nil
}
