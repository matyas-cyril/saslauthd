# go-saslauthd

## [Manpage](man.md)

## Dépendances

```
$ sudo apt install pandoc make
```

## Makefile
```
make [option]

option:
    
    build:
        Générer go-saslauthd, le fichier man et les plugins dans le dossier 'BUILD'

    clean:
        Supprimer les dossiers 'BUILD', 'DEB' et 'plugins'

    deb:
        Créer le paquet 'deb'. Le fichier généré sera dans le dossier 'DEB'

    man:
        Générer le 'man' à partir du fichier 'man.md'

    plugins:
        Compiler les plugins (*.sasl) dans le répertoire 'plugins'

	help:
		Afficher la liste des commandes
```

### Plugin

La déclaration d'un nouveau plugin nécessite les 2 fonctions suivantes :  

```go
/* Vérifier la validité des arguments et générer des données pour l'exploitation du plugin lors de l'appel de la fonction Auth */
func Check(opt map[string]any) (buffer bytes.Buffer, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			buffer = bytes.Buffer{}
			err = fmt.Errorf("panic error plugin NAME_PLUGIN : %s", pErr)
		}
	}()
    /*
    [CODE PROPRE AU PLUGIN]
    */
}

/* Fonction appelée lors de l'authentification.
data contient la trame mise en forme
args les donées générées par la fonction Check */
func Auth(data map[string][]byte, args bytes.Buffer) (valid bool, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			valid = false
			err = fmt.Errorf("panic error plugin NAME_PLUGIN : %s", pErr)
		}
	}()
    /*
    [CODE PROPRE AU PLUGIN]
    */
}
```

#### Variables

##### opt

Variable de type : map[string]any

```go
opt["_pluginPath" ] : // Path complet du répertoire des plugings. Valeur de type string.
```

##### data

Variable de type : data map[string][]byte

```go
data["usr"] : username
data["pwd"] : password
data["srv"] : service
data["dom"] : realm
data["key"] : hash
data["login"] : Si realm existe alors login=usr@dom, sinon login=usr
```

