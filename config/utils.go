package config

import (
	"fmt"
	"os"
	"strings"
)

// dirExist vérifie que le répertoire en argument est bien existant.
// @return : bool true->exist false->absent
func dirExist(directory string) bool {

	f, err := os.Stat(directory)
	if err != nil {
		return false
	}

	return f.IsDir()
}

// anyToStringTab fonction
// data []any : donnée à traiter
// strCase int8 : 0 on ne fait rien, -1 lowercase, 1 uppercase
// uniq bool : Supprimer les doublons si true
// trim bool : Supprimer les espaces si true
func anyToStringTab(data []any, strCase int8, uniq bool, trim bool) ([]string, error) {

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
