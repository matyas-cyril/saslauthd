package ldap_test

import (
	"fmt"
	"log"
	myLdap "saslauthd/plugin/ldap/ldap"
	"testing"
)

// go test -timeout 15s -run ^TestLdap$
func TestLdap(t *testing.T) {

	args := make(map[string]any)

	args["uri"] = "192.168.56.198"
	args["admin"] = "cn=admin,dc=gendarmerie,dc=defense,dc=gouv,dc=fr"
	args["pwd"] = "test"
	args["baseDN"] = "dmdName=Personnes,dc=gendarmerie,dc=defense,dc=gouv,dc=fr"
	args["filter"] = "(uid=%s)"
	args["port"] = 389
	args["timeout"] = 10
	args["tls"] = false

	l, err := myLdap.New(&args)
	if err != nil {
		log.Fatal(err)
	}

	err = l.Connect()

	defer l.Close()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(l.Auth("adam.barry", "test"))
}
