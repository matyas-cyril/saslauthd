package main_test

import (
	"fmt"
	"testing"

	ldap "saslauthd/plugin/ldap"
)

// go test -timeout 5s -run ^TestLdap$
func TestLdap(t *testing.T) {

	opt := make(map[string]any)
	opt["uri"] = "172.0.10.17"
	opt["admin"] = "cn=admin,dc=example,dc=com"
	opt["pwd"] = "AdminLDAP"
	opt["baseDN"] = "ou=people,dc=example,dc=com"

	data, err := ldap.Check(opt)
	if err != nil {
		t.Fatal(err)
	}

	login := make(map[string][]byte)
	login["usr"] = []byte("cyril")
	login["pwd"] = []byte("Cyril")

	auth, err := ldap.Auth(login, data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(auth)

}
