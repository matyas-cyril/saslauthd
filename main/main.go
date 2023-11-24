package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sasl "github.com/matyas-cyril/saslauthd"
)

var APP_CONF string = ""

var APP_PATH = func() (_str string) {

	defer func() {
		if err := recover(); err != nil {
			_str = "."
		}
	}()

	ex, _ := os.Executable()
	return filepath.Dir(ex)
}

// fileExist vérifie que le fichier en argument est bien un fichier existant.
// Il ne vérifie pas la validité, juste l'existance.
// @return : nil -> valide    error -> cause de l'échec
func fileExist(file string) error {

	f, err := os.Stat(file)
	if err != nil {
		return err
	}

	if f.IsDir() {
		return fmt.Errorf(fmt.Sprintf("stat %s: is a directory", file))
	}

	fRead, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fRead.Close()

	return nil
}

func varEnv(env, defaut string) string {

	e := strings.TrimSpace(env)
	if len(e) > 0 && len(e) < 32 {
		return e
	}

	return defaut
}

func main() {

	defaultConfFile := varEnv(APP_CONF, fmt.Sprintf("%s/saslauthd.toml", APP_PATH()))
	var confFile string
	var checkConfFile bool
	flag.StringVar(&confFile, "conf", defaultConfFile, "fichier de configuration de saslauthd")
	flag.BoolVar(&checkConfFile, "check", false, "vérifier la validité du fichier de configuration")
	flag.Parse()

	if err := fileExist(confFile); err != nil {
		fmt.Printf("error loading configuration file:[%s]\n", err.Error())
		os.Exit(1)
	}

	if checkConfFile {

		if err := sasl.Check(confFile); err != nil {
			fmt.Printf("check - config file '%s' checked failed : %s\n", confFile, err.Error())
			os.Exit(1)
		}

		fmt.Printf("check - config file '%s' checked successfully\n", confFile)
		os.Exit(0)
	}

	sasl.Start(confFile, APP_PATH())
}
