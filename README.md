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
Check(opt map[string]interface{}) (bytes.Buffer, error)

/* Fonction appelée lors de l'authentification.
data contient la trame mise en forme
args les donées générées par la fonction Check */
Auth(data map[string][]byte, args bytes.Buffer) (bool, error)

```

#### Variables

##### opt

Variable de type : map[string]interface{}

```go
opt["_pluginPath" ] : // Path complet du répertoire des plugings. Valeur de type string.
```

##### data

Variable de type : data map[string][]byte

```go
data["d0"] : username
data["d1"] : password
data["d2"] : service
data["d3"] : realm
data["d4"] : hash 
```

