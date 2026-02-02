package main_test

import (
	"fmt"
	lemon "saslauthd/plugin/lemon"
	"testing"
)

// go test -timeout 5s -run ^TestLemon$
func TestLemon(t *testing.T) {

	opt := make(map[string]any)
	opt["url"] = "http://hello.world"
	opt["timeout"] = 3

	data, err := lemon.Check(opt)
	if err != nil {
		t.Fatal(err)
	}

	auth, err := lemon.Auth(nil, data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(auth)
}
