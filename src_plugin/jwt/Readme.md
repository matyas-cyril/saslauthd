# JWT

## PRÉSENTATION

Permet d'authentier via un token JWT passé en tant que mot de passe.  

Les algorithmes symétriques pris en comptes :
- HS256
- HS384
- HS512

Les algorithmes asymétriques pris en comptes :
- RS256
- RS384
- RS512

## UTILISATION

La déclaration est de la forme **clef = { option, option,... }**

| OPTION | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| aud | []string | |
| pwd | string | |
| inc | string | |
| virtdom | bool | true |

### **clef :**

Correspond à l'émetteur du token (**ISS**: Issuer).

### **aud :**

Identifier le destinataire auxquel le JWT est destiné.  
On peut spécifier plusieurs aud pour un même token afin d'être multi-services.  

### **pwd :**

Mot de passe au format texte.  
Afin de vérifier la signature du token reçu.

### **inc :**

Inclure un fichier texte externe contenant le mot de passe utile à la vérification de la signature.

### **virtdom :**


## Exemple

``` toml
   [PLUGING.JWT]  
   admin1 = { aud = ["webmail"], pwd = "password" }  
   admin2 = { aud = ["webmail", "printer"], inc = "sample.rsa" }
```
