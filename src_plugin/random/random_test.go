package main_test

import (
	"fmt"
	"testing"

	random "saslauthd/plugin/random"
)

// go test -timeout 5s -run ^TestRandom$
func TestRandom(t *testing.T) {

	opt := make(map[string]any)
	opt["rand"] = 5

	data, err := random.Check(opt)
	if err != nil {
		t.Fatal(err)
	}

	auth, err := random.Auth(nil, data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(auth)
}
