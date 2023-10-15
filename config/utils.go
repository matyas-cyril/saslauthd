package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	toml "github.com/pelletier/go-toml"
)

// getValue permet d'analyser la valeur d'une clef extraite d'un Toml
// tomlTree *toml.Tree : structure du Toml
// key string : clef à analyser
// typeValue string : Type attendu
// @return :
//
//	interface{} : type retourné (bool, string......)
//	error : différent de nil si erreur
func getValue(tomlTree *toml.Tree, key string, typeValue string) (interface{}, error) {

	getType := func(object interface{}) string {

		if object == nil {
			return "nil"
		}

		return reflect.ValueOf(object).Type().String()

	}

	typeValue = strings.TrimSpace(typeValue)
	value := tomlTree.Get(key)

	// On n'a pas de valeur définie dans le fichier toml - Ce n'est pas forcement une erreur
	if value == nil {
		// return nil, fmt.Errorf("No value available for key '%s'", key)
		return nil, nil
	}

	if typeValue != getType(value) {
		return nil, fmt.Errorf("value must be a type '%s' for key '%s' not '%s'", typeValue, key, getType(value))
	}

	return value, nil
}

// dirExist vérifie que le répertoire en argument est bien existant.
// @return : bool true->exist false->absent
func dirExist(directory string) bool {

	f, err := os.Stat(directory)
	if err != nil {
		return false
	}

	return f.IsDir()
}

// interfaceToStringTab fonction
// data []interface{} : donnée à traiter
// strCase int8 : 0 on ne fait rien, -1 lowercase, 1 uppercase
// uniq bool : Supprimer les doublons si true
// trim bool : Supprimer les espaces si true
func interfaceToStringTab(data []interface{}, strCase int8, uniq bool, trim bool) ([]string, error) {

	s := make([]string, len(data))
	for i, j := range data {

		str := fmt.Sprint(j)

		// Supprimer les espaces gauche/droite
		if trim {
			str = strings.TrimSpace(str)
		}

		if strCase < 0 {
			s[i] = strings.ToLower(str)
		} else if strCase > 0 {
			s[i] = strings.ToUpper(str)
		} else {
			s[i] = str
		}
	}

	// Supprimer les doublons
	if uniq {
		s = deleteDuplicatesFromSlice(s)
	}

	return s, nil
}

func deleteDuplicatesFromSlice(s []string) []string {

	m := make(map[string]bool)

	for _, item := range s {

		if !m[item] {
			m[item] = true
		}
	}

	var result []string
	for _, item := range s {
		if m[item] {
			result = append(result, item)
			delete(m, item)
		}

	}

	return result
}
