# LDAP

## PRÉSENTATION

Permet d'effectuer une authentificaiton utilisateur via LDAP.

## OPTIONS

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| uri | string | 127.0.0.1 |
| admin | string |
| pwd | string |
| baseDN | string |
| filter | string | (uid=%s) |
| port | int | 389 |
| timeout | int | 10 |
| att | string | dn |
| attMatch | string | uid |
| tls | bool | false |
| tlsSkipVerify | bool | true |

### **uri :**

Adresse du serveur LDAP.

### **port :**

Port sur lequel le serveur LDAP écoute. Le port 389 est standard pour LDAP.

### **admin :**

Utilisateur avec des droits élargis.

### **pwd :**

Mot de passe associé à l'utilisateur admin pour authentification.

### **baseDN :**

Le DN de la racine de la base de données LDAP à partir de laquelle les recherches sont effectuées.

### **filter :**

Critère de filtrage pour les requêtes LDAP. Le %s sera remplacé par la valeur de recherche.

### **timeout :**

Durée, en secondes, avant qu'une requête ne soit abandonnée.

    0 :
        Pas de timeout
    
    10 :
        Défaut

    3600 :
        Max 

### **att :**

Attribut utilisé pour déterminer l'utilisateur à authentifier.

### **attMatch :**

Attribut utiliser pour matcher avec l'identiant de connexion.

### **tls :**

Indique si la connexion doit être sécurisée avec TLS (LDAPS).  
Si le port n'est pas défini par l'utilisateur, alors port=636.

### **tlsSkipVerify :**

Détermine si la validation du certificat TLS doit être ignorée.  
Utile lors de la connexion à des serveurs avec des certificats non valides.