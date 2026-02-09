package main

import (
	"fmt"
	"reflect"
	"strings"
)

type PgAuth struct {
	Cnx     string
	Timeout uint16
	Sql     string // Requête permettant d'obtenir les informations d'authentification
	Virtdom bool   /* Prise en compte des domaines pour composer le login
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
		virtdom bool
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
				if !strings.Contains(kV, "$1") {
					return nil, fmt.Errorf("pgauth key '%s' not valid", k)
				}
				sql = kV
			}

		case "port", "timeout":

			// Type souhaité
			typeTarget := reflect.TypeFor[int]()
			rv := reflect.ValueOf(data[k])
			if !rv.Type().AssignableTo(typeTarget) {
				return nil, fmt.Errorf("pgauth key '%s' failed to typecast", k)
			}

			nbr := rv.Convert(typeTarget).Int()
			if nbr < 0 || nbr > 65535 {
				return nil, fmt.Errorf("pgauth key '%s' integer range invalid", k)
			}

			switch k {
			case "port":
				port = uint16(nbr)
			case "timeout":
				timeout = uint16(nbr)
			}

		case "virdom":
			kV, kErr := v.(bool)
			if !kErr {
				return nil, fmt.Errorf("pgauth key '%s' failed to typecast", k)
			}
			virtdom = kV

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
		Virtdom: virtdom,
	}, nil

}
