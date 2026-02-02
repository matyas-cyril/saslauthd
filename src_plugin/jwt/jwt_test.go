package main_test

import (
	"fmt"
	pluginJwt "saslauthd/plugin/jwt"
	"testing"
)

// go test -timeout 5s -run ^TestCheck$
func TestCheck(t *testing.T) {

	// map[admin1:map[aud:[webmail] pwd:password] admin2:map[aud:[webmail] inc:sample.rsa]]
	opt := make(map[string]any)
	opt["admin1"] = map[string]any{
		"aud": []string{"webmail"},
		"pwd": "password",
	}
	opt["admin2"] = map[string]any{
		"aud": []string{"webmail"},
		"inc": "sample.rsa",
	}

	dataOpt, err := pluginJwt.Check(opt)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(dataOpt)

}

// go test -timeout 5s -run ^TestAuthRSA$
func TestAuthRSA(t *testing.T) {

	opt := make(map[string]any)
	opt["admin"] = map[string]any{
		"aud": []string{"webmail"},
		"inc": "sample.rsa",
	}

	dataOpt, err := pluginJwt.Check(opt)
	if err != nil {
		t.Fatal(err)
	}

	// Token expiré -
	// rawToken := []byte("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhZG1pbiIsInVzciI6InJlbmUubGF0YXVwZSIsImRvbSI6InRlc3QuZnIiLCJhdWQiOiJ3ZWJtYWlsIiwiZXhwIjoxNzM1Njg5NjAwfQ.MM78fLdl9nYIm6mPCIvRR5yK5DN0tqsg7S3sP7JE1dl-HYk1d_n_7zhdng20R--GnX5hYGYJz6wjUsJN6cK4B4h2tinKeEIkpMnvWTYsUeR6t42wDtZlfGYFghujZZUhruRHHaEhTw1gb5xy5Jj9DguLDvkANFtoCR8upDWw9FRh469ehgdReUkHeqnP-9eI9ZvIMPJLRoXCLL63c2ml9atsRSBGGAizIa_y6AJfiABbwsoN8ysVFaJ0U9GQGXt79r13774vv_UvZG3bblCo7dzA03EB9GOgq_NSoWJ0uJDrPdj5mWfg3CCOxjeOu1Bxk0AImFJiOJaru9QUU2FJGQ")

	// Token valide jusqu'au 01 janvier 2035 - 12h00
	rawToken := []byte("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhZG1pbiIsInVzciI6InJlbmUubGF0YXVwZSIsImRvbSI6InRlc3QuZnIiLCJhdWQiOiJ3ZWJtYWlsIiwiZXhwIjoyMDUxMjY1NjAwfQ.E0fgS34jq01tO0Zpqbt3sR08VpxW5rshDbhXaB84FYUOGFaoj7qSL4zsJo3f6oKB3DlWP-8IW3En5AmJQuto2xRA2fJ363aRtQOILMzmx2dlS7Xjk_A1vTXOS1x9I1ra_IEp0Tw3ZCJCvRg9RCKGsJH__KeOOe7BXbCpfrw0soJqwCswHnuIs-0NpZJOS6Xtw--qKOrz5VjgSdA9R0SKQRTxZH9QOiRYOXxzQGJeNiBA-GC24vmy9Ag8cINeubE45GMWHAmyF-C0wsZK7dTIrXxyebB1pIkdcpc2iAl6RYebaspbH6ll1G2bjqkfaPRxrmOBvbLKpJKIPLZW_Guevg")

	data := make(map[string][]byte)
	data["d0"] = []byte("rene.lataupe")
	data["d1"] = rawToken
	data["d2"] = []byte("imap")
	data["d3"] = []byte("test.fr")

	rtn, err := pluginJwt.Auth(data, dataOpt)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(rtn)

}

// go test -timeout 5s -run ^TestAuth$
func TestAuth(t *testing.T) {

	opt := make(map[string]any)
	opt["admin"] = map[string]any{
		"aud": []string{"webmail"},
		"pwd": "password",
	}

	dataOpt, err := pluginJwt.Check(opt)
	if err != nil {
		t.Fatal(err)
	}

	rawToken := []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhZG1pbiIsInVzciI6InJlbmUubGF0YXVwZSIsImRvbSI6InRlc3QuZnIiLCJhdWQiOiJ3ZWJtYWlsIiwiZXhwIjoxNzM1Njg5NjAwfQ.QywpRimRgmWKI7IKQdreyioJXHKoTe6X2q1Ey21d8e4")

	data := make(map[string][]byte)
	data["d0"] = []byte("rene.lataupe")
	data["d1"] = rawToken
	data["d2"] = []byte("imap")
	data["d3"] = []byte("test.fr")

	rtn, err := pluginJwt.Auth(data, dataOpt)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(rtn)
}

// go test -timeout 5s -run ^TestAuthUid$
func TestAuthUid(t *testing.T) {

	opt := make(map[string]any)
	opt["admin"] = map[string]any{
		"aud": []string{"webmail"},
		"pwd": "password",
	}

	dataOpt, err := pluginJwt.Check(opt)
	if err != nil {
		t.Fatal(err)
	}

	rawToken := []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhZG1pbiIsImF1ZCI6IndlYm1haWwiLCJ1aWQiOiJyZW5lLmxhdGF1cGVAdGVzdC5mciIsImlhdCI6MTUxNjIzOTAyMiwibmJmIjoxNTE2MjM5MDIyLCJleHAiOjI1MTYyMzkwMjJ9.bGJO7u7WdkKY62gNNjS34ngcCBiwKiS-WhraEVWzrDY")
	data := make(map[string][]byte)
	data["d0"] = []byte("rene.lataupe")
	data["d1"] = rawToken
	data["d2"] = []byte("imap")
	data["d3"] = []byte("test.fR")

	rtn, err := pluginJwt.Auth(data, dataOpt)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(rtn)
}
