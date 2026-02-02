package main_test

import (
	"testing"

	ldap "saslauthd/plugin/ldap"
)

// go test -timeout 5s -run ^TestLdap$
func TestLdap(t *testing.T) {

	ldap.Check(nil)

}
