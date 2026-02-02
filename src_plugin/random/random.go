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

func Check(opt map[string]any) (buffer bytes.Buffer, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			buffer = bytes.Buffer{}
			err = fmt.Errorf("panic error plugin random : %s", pErr)
		}
	}()

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

func Auth(data map[string][]byte, args bytes.Buffer) (valid bool, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			valid = false
			err = fmt.Errorf("panic error plugin random : %s", pErr)
		}
	}()

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

func interfaceToStruct(data map[string]any) (*randStruct, error) {

	r := randStruct{}

	for k := range data {

		if k == "rand" {

			// Type souhaité
			typeTarget := reflect.TypeFor[int]()

			rv := reflect.ValueOf(data[k])
			if !rv.Type().AssignableTo(typeTarget) {
				return nil, fmt.Errorf("random param key %s must be an integer", k)
			}

			nbr := rv.Convert(typeTarget).Int()
			if nbr < 0 {
				nbr = 0
			} else if nbr > 120 {
				nbr = 120
			}

			r.Rand = uint8(nbr)

		} else {
			return nil, fmt.Errorf("random param key %s not exist", k)
		}

	}

	return &r, nil
}
