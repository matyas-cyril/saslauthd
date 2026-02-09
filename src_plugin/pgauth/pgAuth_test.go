package main_test

import (
	"fmt"
	"testing"

	pgauth "saslauthd/plugin/postauth"
)

// go test -timeout 5s -run ^TestPgAuth$
func TestPgAuth(t *testing.T) {

	opt := make(map[string]any)
	opt["user"] = "photocopieur_bdd"
	opt["passwd"] = "PasswordBdd"
	opt["host"] = "172.0.200.30"
	opt["bdd"] = "photocopieur_bdd"
	opt["sql"] = "SELECT username, password FROM v_active_users WHERE username LIKE $1 LIMIT 1"
	opt["timeout"] = 3
	opt["realm"] = false

	args, err := pgauth.Check(opt)
	if err != nil {
		t.Fatal(err)
	}

	data := make(map[string][]byte)
	data["usr"] = []byte("cyril")
	data["dom"] = []byte("test")

	data["login"] = data["usr"]
	if len(data["dom"]) != 0 {
		data["login"] = fmt.Appendf(nil, "%s@%s", data["usr"], data["dom"])
	}

	auth, err := pgauth.Auth(data, args)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(auth)
}
