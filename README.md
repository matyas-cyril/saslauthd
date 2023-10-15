---
title: saslauthd
section: 4
header: manuel utilisateur
footer: saslauthd 1.0.0
---

# NOM

saslauthd - Serveur d'authentification SASL

# SYNOPSIS

    saslauthd [OPTIONS]

# DESCRIPTION

    Serveur d'authentification plaintext SASL. 
    Permet de traiter plusieurs méthodes d'authentification.

# OPTIONS

    --conf :
        Permet de prendre en compte un fichier de configuration autre que celui par défaut (/opt/saslauthd/saslauthd.toml).

    --check :
        Permet de vérifier un fichier de configuration sans prise en compte.

# CODES RETOUR

    0 : RAS  
    1 : Fichier de configuration absent ou non valide
    2 : Echec d'activation du mode debug
    3 : Echec socket

# FICHIER DE CONFIGURATION DEFAUT

    /opt/saslauthd/saslauthd.toml

## **[SERVER]**

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| socket | string | /var/run/saslauthd/mux |
| user | string | mail |
| group | string | mail |
| rate_info | int | 30 |
| client_max | int | 100 |
| client_timeout | int | 30 |
| buffer_size | int | 256 |
| buffer_timeout | int | 50 |
| buffer_hash | string | sha256 |
| socket_size | int | 1024 |
| plugin_path | string | ./plugins |
| log | string | TERM |

### **socket :**

    Socket de communication.

### **user :**

    Définir l'utilisateur du socket

### **group :**

    Définir le groupe du socket

### **rate_info :**

Fréquence en secondes de l'export des informations techniques sur serveur dans les logs.  
Voir l'option 'log' pour connaître le type d'affichage de l'export.

    0 :
        Désactiver

    30 :
        Valeur par défaut

    3600 :
        Valeur max



### **client_max :**

    Nombre de clients autorisés à se connecter.  

    0 :
        Pas de restriction. Open bar.

    100 :
        Valeur par défaut

    500000000 :
        Valeur max autorisée

### **client_timeout :**

Durée maximum en secondes d'une connexion client.

    0 :
        Pas de restriction. Si le plugin ne gère pas de timeout, alors la connexion peut rester ouverte !!!

    30 :
        Valeur par défaut

    240 :
        Valeur max autorisée

### **buffer_size :**

Taille du buffer (byte) utilisée lors de la lecture de la socket.

    Minimum: 1
    Maximum: 2048
    Défaut: 256

> on doit avoir **socket_size** >= **buffer_size**

### **buffer_timeout :**

Mise en place d'un timeout (milliseconde) pour forcer la sortie, lorsque le contenu de la socket est inférieur ou égal au buffer.  

    Minimum: 1
    Maximum: 10000
    Défaut: 50

### **buffer_hash :**

Hash calculé à partir de la trame reçue.  
Cette valeur va servir de référence pour le cache.  

    Possibilités: 
        md5, sha1, sha256, sha512
    
    Défaut:
        sha256

### **socket_size :**

Taille maximum (byte) de la taille de socket. Au dela de cette taille, la connection client est interrompue. Et c'est un échec d'authentification.

    Minimum: 8
    Maximum: 32768
    Défaut: 1024

> on doit avoir **socket_size** >= **buffer_size**

### **plugin_path :**

Répertoire contenant les plugins.

### **log :**

Définir le type de journalisation : NO | TERM | SYSLOG | BOTH

    NO :
        Pas de d'affichage

    TERM :
        Affichage des informations dans le terminal

    SYSLOG :
        Affichage dans Syslog (identique aux mails).  
        Nécessite rsyslog.

    BOTH:
        Affichage dans le terminal (TERM) et dans syslog (SYSLOG) 

*Si **log = "NO"**, rien ne sera affiché.*

---

## **[DEBUG]**

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| enable | bool | false |
| file | string | /tmp/saslauthd.debug |

### **enable :**

    false :
        Désactiver le mode debug

    true :
        Activer le mode debug

### **file :**

Fichier de sortie des informations lors de l'action du mode debug.  
Les données de mot de passe ou de clef de chiffrement sont remplacés par des caractères 'x' -> 120 

---

## **[CACHE]**

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| enable | bool | false |
| type | string | LOCAL |
| key | string |
| key_rand| bool | false |
| ok | int | 60 |
| ko | int | 60 |

### **enable :**

    false :
        Désactiver le cache

    true :
        Activer le cache

### **type :**

Définir le type de mise en cache utilisé : LOCAL | MEMCACHE | REDIS

    LOCAL :  
        Mise en cache local. Voir le bloc [CACHE.LOCAL].

    MEMCACHE :  
        Pas encore implémenté.

    REDIS :  
        Pas encore implémenté.

### **key :**
    Définir une clef symétrique. Si le champ est vide ou reste à la valeur par défaut, le chiffrement est désactivé.  

### **key_rand :**

    false :
        Ne pas générer une clef de chiffrement aléatoire.

    true :
        Générer une clef de chiffrement aléatoire à chaque démarrage. 
        Écrase la valeur de 'key'.

*L'activation de cette fonctionnalité est déconseillée lors de l'utilisation du cache distribué (MEMCACHE | REDIS).*

### **ok :**

Durée en secondes de la mise en cache du succès d'authentifiation.

    0 :
        Désactiver la mise en cache
    
    60 :
        Valeur par défaut

    31536000 :
        Valeur max (1 an)

### **ko :**

Durée en secondes de la mise en cache de léchec d'authentifiation.

    0 :
        Désactiver la mise en cache
    
    60 :
        Valeur par défaut

    31536000 :
        Valeur max (1 an)

---

## **[CACHE.LOCAL]**

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| path | string | /tmp |
| sweep | int | 60 |
| purge_on_start | bool | false |

### **path :**

Répertoire de stockage des fichiers de lors de la mise en cache.  
Le path du répertoire doit être absolu.

### **sweep :**

Fréquence en **secondes** de l'exécution de la suppression des fichiers de caches obsolètes.   

    0 :
        Désactiver
    
    60 :
        Valeur par défaut

    86400 :
        Valeur max (24 heures)

### **purge_on_start :**

Supprime l'ensemble des fichiers correspondant au pattern du cache.  
Même si le fichier cache est valide, il sera supprimé.

    false :
        Désactivé

    true :
        Activé

---

## **[AUTH]**

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| mech | []string | ["NO"] |
| auth_multi | bool | false |

### **mech :**

    Liste des mechanismes d'authentification.  
    Il existe 2 mechanismes internes [YES | NO].  
    Les autres mechanismes sont des plugins.  
    Ils sont définis dans des sous-sections de [PLUGIN].
    L'absence de mechanisme correspond à NO.

### **auth_multi :**

Activer le traitement des authentifications par lot de 3.  
NON PRIS EN COMPTE DANS LA VERSION ACTUELLE.

    false :
        Désactivé

    true :
        Activé

## **[PLUGIN.RANDOM]**

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| rand | int | 0 |

### **rand :**

Définir la durée en seconde de la suspension d'exécution.

    0 :
        Min
    
    120 :
        Max

    0 :
        Défaut

## **[PLUGIN.LDAP]**

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| uri | string |
| admin | string |
| pwd | string |
| baseDN | string |
| filter | string |
| port | int | 389 |
| timeout | int | 10 |
| tls | bool | false |
| tlsSkipVerify | bool | true |

## **[PLUGIN.JWT]**

Permet d'authentier via un token JWT passé en tant que mot de passe.  

Les algorithmes symétriques pris en comptes :
- HS256
- HS384
- HS512

Les algorithmes asymétriques pris en comptes :
- RS256
- RS384
- RS512

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| aud | []string |
| *auth_user* | string |

### **aud :**

Identifier le destinataire auxquel le JWT est destiné.  
Il n'y a pas de relation entre la clef de déchiffrement et l'aud. 

### ***auth_user***:

***auth_user*** correspond à la valeur de **iss** ddu payload.

En fonction du type de chiffrement (symétrique ou asymétrique), la syntaxe de la clef différe.  

#### **clef symétrique**

auth_user: clef

#### **clef asymétrique**

auth_user: **@include:**FICHIER_CLEF_SYMETRIQUE

La syntaxe **@include:** signifie que la clef publique est présente dans un fichier externe.  

## Exemples

    aud: ["webmail", "thunderbird", "tph"]
    user1: "mot_de_passe_symetrique"
    user2: "@include:sample.rsa"


## **[PLUGIN.LEMON]**

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