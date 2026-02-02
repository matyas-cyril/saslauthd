# PGAUTH

## PRÉSENTATION

## OPTIONS

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| host | string | 127.0.0.1 |
| port | int | 5432 |
| bdd | string | |
| user | string | |
| passwd | string | |
| timeout | int | 5 |
| realm | bool | false |
| sql | string | |

### **host :**

Adresse de connexion à la BDD

### **port :**

Port d'écoute de la BDD

### **bdd :**

Nom de la BDD

### **user :**

Utilisateur utilisé pour l'authentification à la BDD

### **passwd :**

Mot de passe d'authentification à la BDD

### **timeout :**

Définir le timeout (en seconde) de la requête

    0 :
        Pas de timeout
    
    5 :
        Défaut

    3600 :
        Max 

### **realm :**

Si **true** l'identifiant se présentant en uid@domain devra exister en BDD.  
Si **false**, même si l'utilisateur se présente en uid@domain, l'authentification se fera avec uid uniquement.

### **sql :**

Reqête SQL permettant d'obtenir l'username et le password correspond à l'identifiant.  
**$1** correspond à l'identifiant et est obligatoire.

```sql
SELECT username, password FROM users WHERE username LIKE $1 LIMIT 1
```