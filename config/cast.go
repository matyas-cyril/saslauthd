package config

import (
	"fmt"
	"strings"
)

func castString(v any) (str string, err error) {

	defer func() {
		if p := recover(); p != nil {
			str = ""
			err = fmt.Errorf("param type not valid to cast to string")
		}
	}()

	str = strings.TrimSpace(v.(string))
	if len(str) == 0 {
		return "", fmt.Errorf("string length equal zero")
	}

	return str, nil
}

func castUint8(v any) (nbr uint8, err error) {

	defer func() {
		if p := recover(); p != nil {
			nbr = 0
			err = fmt.Errorf("param type not valid to cast to uint8")
		}
	}()

	return uint8(v.(int64)), nil
}

func castUint16(v any) (nbr uint16, err error) {

	defer func() {
		if p := recover(); p != nil {
			nbr = 0
			err = fmt.Errorf("param type not valid to cast to uint16")
		}
	}()

	return uint16(v.(int64)), nil
}

func castUint32(v any) (nbr uint32, err error) {

	defer func() {
		if p := recover(); p != nil {
			nbr = 0
			err = fmt.Errorf("param type not valid to cast to uint32")
		}
	}()

	return uint32(v.(int64)), nil
}

func castBool(v any) (b bool, err error) {

	defer func() {
		if p := recover(); p != nil {
			b = false
			err = fmt.Errorf("param type not valid to cast to boolean")
		}
	}()

	b = v.(bool)

	return b, nil
}
