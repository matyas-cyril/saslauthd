# LDAP

## PRÉSENTATION

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

Adresse permettant de vérifier le token.

### **admin :**

### **pwd :**

### **baseDN :**

### **filter :**

### **port :**

### **timeout :**

Définir le timeout (en seconde) de la requêtre ldap

    0 :
        Pas de timeout
    
    10 :
        Défaut

    3600 :
        Max 

### **att :**

### **attMatch :**

### **tls :**

### **tlsSkipVerify :**
