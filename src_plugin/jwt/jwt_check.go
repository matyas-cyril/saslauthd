package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	myJwt "github.com/cristalhq/jwt/v5"
)

// Contrôle du payload
func checkJwt(usr, dom []byte, jwt *jwtStruct) error {

	// Prise en compte de l'UID au lieu de USR et DOM
	if len(jwt.Uid) > 0 {

		uid := string(usr)
		// On compare usr@dom à jwt.uid
		if len(dom) != 0 {
			uid = fmt.Sprintf("%s@%s", usr, dom)
		}

		// On vérifie que l'uid correspond
		if !strings.EqualFold(uid, jwt.Uid) {
			return fmt.Errorf("jwt uid payload not match")
		}

	} else {

		// On vérifie que l'user correspond
		if !strings.EqualFold(string(usr), jwt.Usr) {
			return fmt.Errorf("jwt usr payload not match")
		}

		// On vérifie que le domaine correspond
		if !strings.EqualFold(string(dom), jwt.Dom) {
			return fmt.Errorf("jwt dom payload not match")
		}

	}

	// Iss obligatoire
	if len(strings.TrimSpace(jwt.Iss)) == 0 {
		return fmt.Errorf("jwt iss claim must be defined")
	}

	// Aud obligatoire
	if len(strings.TrimSpace(jwt.Aud)) == 0 {
		return fmt.Errorf("jwt aud claim must be defined")
	}

	// Exp doit être défini
	if jwt.Exp == 0 {
		return fmt.Errorf("jwt exp claim must be defined")
	}

	epoch := time.Now().Unix()

	// Token expiré !!!
	if jwt.Exp < uint64(epoch) {
		return fmt.Errorf("jwt token is expired - claim exp")
	}

	// utilisation du NotBefore
	if jwt.Nbf > 0 {

		// Not before >= à l'expiration !!!!
		if jwt.Nbf >= jwt.Exp {
			return fmt.Errorf("jwt token invalid - claim nbf must lower than claim exp")
		}

		// Token pas encore utilisable
		if jwt.Nbf > uint64(epoch) {
			return fmt.Errorf("jwt token not be accepted for processing - claim nbf")
		}
	}

	return nil
}

// Cotrôle de la signature
func checkSignJwt(rawToken, key []byte, algo string) error {

	algo = strings.TrimSpace(strings.ToUpper(algo))

	switch algo {

	case "HS256", "HS384", "HS512":
		return verifyJwtSignHS(rawToken, key, algo)

	case "RS256", "RS384", "RS512":
		return verifyJwtSignRS(rawToken, key, algo)

	case "ES256", "ES384", "ES512":
		return fmt.Errorf("jwt algo '%s' not yet implemented", algo)

	case "PS256", "PS384", "PS512":
		return fmt.Errorf("jwt algo '%s' not yet implemented", algo)

	}

	return fmt.Errorf("algo %s not valid", algo)
}

func verifyJwtSignHS(rawToken, key []byte, algo string) error {

	verif, err := myJwt.NewVerifierHS(myJwt.Algorithm(algo), key)
	if err != nil {
		return err
	}

	jwtToken, err := myJwt.Parse(rawToken, verif)
	if err != nil {
		return err
	}

	// On vérifie la validité de la signature
	if err := verif.Verify(jwtToken); err != nil {
		return err
	}

	return nil
}

func verifyJwtSignRS(rawToken, key []byte, algo string) (rtnErr error) {

	defer func() {

		if err := recover(); err != nil {
			rtnErr = fmt.Errorf("verify jwt rsa signature panic error : %s", err)
		}

	}()

	rsaPubKey, err := byteToRSAPub(key)
	if err != nil {
		return err
	}

	verif, err := myJwt.NewVerifierRS(myJwt.Algorithm(algo), rsaPubKey)
	if err != nil {
		return err
	}

	jwtToken, err := myJwt.Parse(rawToken, verif)
	if err != nil {
		return err
	}

	// On vérifie la validité de la signature
	if err := verif.Verify(jwtToken); err != nil {
		return err
	}

	return nil
}

func byteToRSAPub(key []byte) (rtnPubKey *rsa.PublicKey, rtnErr error) {

	defer func() {

		if err := recover(); err != nil {
			rtnPubKey = nil
			rtnErr = fmt.Errorf("decode rsa pub key panic : %s", err)
		}

	}()

	block, _ := pem.Decode(key)

	var pubKey *rsa.PublicKey
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pubKey = pubInterface.(*rsa.PublicKey)

	return pubKey, nil
}
