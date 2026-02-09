# LEMON

## PRÉSENTATION

## OPTIONS

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| url | string |
| timeout | int | 5 |
| active | string | active |
| authkey | string | mail |
| virtdom | bool | true |

### **url :**

Adresse permettant de vérifier le token.

### **timeout :**

Définir le timeout (en seconde) de la requêtre http

    0 :
        Pas de timeout
    
    5 :
        Défaut

    3600 :
        Max 

### **active :**

Définir la variable permettant de définir le status d'un compte.  
La valeur par défaut est "active".  
Une chaine vide désactive le contrôle.  

### ** authkey :**

Définir la clef retournée par la SSO permettant d'effectuer l'authentification.  
Le choix possible est : MAIL ou UID.

### **virtdom :**

Si virtdom est **true** alors l'authentifcation sera par défaut en user@dom si dom est définir, sinon user.  
Si virtdom est **false** l'authentification sera effectuée avec user.
