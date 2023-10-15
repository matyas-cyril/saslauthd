# LEMON

## PRÉSENTATION

## OPTIONS

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| url | string |
| timeout | int | 5 |
| active | string | active |

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