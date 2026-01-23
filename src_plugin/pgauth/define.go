package main

import (
	"fmt"
	"strings"
)

type PgAuth struct {
	Cnx     string
	Timeout uint16
	Sql     string // Requête permettant d'obtenir les informations d'authentification
	Realm   bool   /* Prise en compte des domaines pour composer le login
	Si false : login = username, si true login == username@real	*/
}

func interfaceToStruct(data map[string]any) (*PgAuth, error) {

	var (
		host    string = "127.0.0.1"
		port    uint16 = 5432
		user    string
		passwd  string
		bdd     string
		timeout uint16 = 5
		realm   bool
		sql     string
	)

	for k, v := range data {

		switch k {

		case "host", "user", "passwd", "bdd", "sql":
			kV, kErr := v.(string)
			if !kErr {
				return nil, fmt.Errorf("pgauth key '%s' failed to typecast", k)
			}
			kV = strings.TrimSpace(kV)
			if len(kV) == 0 {
				return nil, fmt.Errorf("pgauth key '%s' defined but empty", k)
			}
			switch k {
			case "host":
				host = kV
			case "user":
				user = kV
			case "passwd":
				passwd = kV
			case "bdd":
				bdd = kV
			case "sql":
				sql = kV
			}

		case "port", "timeout":
			kV, kErr := v.(uint16)
			if !kErr {
				return nil, fmt.Errorf("pgauth key '%s' failed to typecast", k)
			}
			switch k {
			case "port":
				port = kV
			case "timeout":
				timeout = kV
			}

		case "realm":
			kV, kErr := v.(bool)
			if !kErr {
				return nil, fmt.Errorf("pgauth key '%s' failed to typecast", k)
			}
			realm = kV

		default:
			return nil, fmt.Errorf("pgauth key '%s' not exist", k)
		}

	}

	// Vérifier la présences des paramètres obligatoires
	if len(user) == 0 {
		return nil, fmt.Errorf("pgauth key 'user' must defined")
	}
	if len(passwd) == 0 {
		return nil, fmt.Errorf("pgauth key 'passwd' must defined")
	}
	if len(bdd) == 0 {
		return nil, fmt.Errorf("pgauth key 'bdd' must defined")
	}
	if len(sql) == 0 {
		return nil, fmt.Errorf("pgauth key 'sql' must defined")
	}

	return &PgAuth{
		Cnx:     fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, passwd, host, port, bdd),
		Sql:     sql,
		Timeout: timeout,
		Realm:   realm,
	}, nil

}
