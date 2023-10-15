package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"time"
)

func Check(opt map[string]interface{}) (bytes.Buffer, error) {

	var buffer bytes.Buffer

	// convertir l'interface en structure compréhensible par le plugin
	data, err := interfaceToStruct(opt)
	if err != nil {
		return bytes.Buffer{}, err
	}

	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(data); err != nil {
		return bytes.Buffer{}, err
	}

	return buffer, nil
}

func Auth(data map[string][]byte, args bytes.Buffer) (bool, error) {

	var arg randStruct

	dec := gob.NewDecoder(&args)
	if err := dec.Decode(&arg); err != nil && err != io.EOF {
		return false, err
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	if arg.Rand != 0 {
		nbrRand := rnd.Intn(int(arg.Rand) + 1)
		time.Sleep(time.Duration(nbrRand) * time.Second)
	}

	nbr := rand.New(rand.NewSource(time.Now().UnixNano()))

	if nbr.Uint32()%2 == 0 {
		return true, nil
	}

	return false, nil
}

type randStruct struct {
	Rand uint8
}

func interfaceToStruct(data map[string]interface{}) (*randStruct, error) {

	r := randStruct{}

	for k := range data {

		if k == "rand" {

			if !reflect.ValueOf(r.Rand).Type().ConvertibleTo(reflect.ValueOf(data[k]).Type()) {
				return nil, fmt.Errorf(fmt.Sprintf("random param key %s must be an integer", k))
			}

			nbr := data[k].(int64)
			if nbr < 0 {
				nbr = 0
			} else if nbr > 120 {
				nbr = 120
			}
			r.Rand = uint8(nbr)

		} else {
			return nil, fmt.Errorf(fmt.Sprintf("random param key %s not exist", k))
		}

	}

	return &r, nil
}
