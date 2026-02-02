package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5"
	argon2id "github.com/matyas-cyril/argon2id"
)

func Check(opt map[string]any) (buffer bytes.Buffer, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			buffer = bytes.Buffer{}
			err = fmt.Errorf("panic error plugin pgauth : %s", pErr)
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
			err = fmt.Errorf("panic error plugin pgauth : %s", pErr)
		}
	}()

	var arg PgAuth
	dec := gob.NewDecoder(&args)
	if err := dec.Decode(&arg); err != nil && err != io.EOF {
		return false, err
	}

	if len(data["d0"]) == 0 {
		return false, fmt.Errorf("login not defined")
	}

	if len(data["d1"]) == 0 {
		return false, fmt.Errorf("password not defined")
	}

	login := data["d0"]
	password := data["d1"]

	// On utilise les domaines pour composer l'authentification
	if arg.Realm {
		if len(data["d3"]) == 0 {
			return false, fmt.Errorf("login generation failure. Realm not defined")
		}
		// Login -> login@realm
		login = append(append(append([]byte{}, login...), byte('@')), data["d3"]...)
	}

	// Connexion à la BDD
	cnx, err := pgx.Connect(context.Background(), arg.Cnx)
	if err != nil {
		return false, fmt.Errorf("unable to connect to database: %v\n", err)
	}
	defer cnx.Close(context.Background())

	// Obtenir les informations en fonction du login utilisé
	var bddUser []byte
	var bddPass []byte
	err = cnx.QueryRow(context.Background(), arg.Sql, login).Scan(&bddUser, &bddPass)
	if err != nil {
		return false, fmt.Errorf("query failed: %v", err)
	}

	// On vérifie que l'username retourné par la BDD correspond bien à celui attendu...
	// Cela peut se produire si la requête peut retourner plusieurs lignes.
	// QueryRow retourne 0 ou 1 ligne uniquement. Toutes les autres sont ignorées
	if !bytes.Equal(login, bddUser) {
		return false, fmt.Errorf("different login from database return. Possible duplication or usurpation")
	}

	// Vérifier le mot de passe argon2id
	return argon2id.Check(password, bddPass)
}
