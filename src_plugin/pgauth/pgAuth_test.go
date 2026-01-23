package main_test

import (
	"fmt"
	"testing"

	pgauth "saslauthd/plugin/postauth"
)

// go test -timeout 5s -run ^TestPgAuth$
func TestPgAuth(t *testing.T) {

	opt := make(map[string]any)
	args, err := pgauth.Check(opt)
	if err != nil {
		t.Fatal(err)
	}

	data := make(map[string][]byte)
	data["d0"] = []byte("cyril")
	data["d1"] = []byte("test")

	auth, err := pgauth.Auth(data, args)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(auth)
}
