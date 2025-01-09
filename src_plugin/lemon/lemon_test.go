package main_test

import (
	"fmt"
	lemon "saslauthd/plugin/lemon"
	"testing"
)

func TestLemon(t *testing.T) {

	opt := make(map[string]interface{})

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
