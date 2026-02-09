---
title: go-saslauthd
section: 4
header: manuel utilisateur
footer: go-saslauthd
---

# NOM

go-saslauthd - Serveur d'authentification SASL

# SYNOPSIS

    go-saslauthd  [OPTIONS]

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

    /opt/go-saslauthd/saslauthd.toml

## **[SERVER]**

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| socket | string | /var/run/saslauthd/mux |
| user | string | mail |
| group | string | mail |
| rate_info | int | 30 |
| client_max | int | 100 |
| client_timeout | int | 30 |
| graceful | int | 5 |
| buffer_size | int | 256 |
| buffer_timeout | int | 50 |
| buffer_hash | string | sha256 |
| socket_size | int | 8192 |
| plugin_path | string | $APP_PATH/plugins |
| log | string | TERM |
| log_facility | string | AUTH |

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

### **graceful :**

Durée maximum en secondes pendant laquelle le serveur attend la fin des transactions avec le client. Pendant cette période il n'accepte plus de nouvelles connexions.

    0 :
        Pas de graceful shutdown, directement du hard shutdown

    5 :
        Valeur par défaut

    60 :
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
    Défaut: 8192

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

### **log_facility :**

Définir le sous-système applicatif pour l'enregistrement des logs.

```
Les possibilités sont :  
    AUTH | MAIL | SYSLOG | USER | LOCAL0 | LOCAL1 | LOCAL2 | LOCAL3 | LOCAL4 | LOCAL5 | LOCAL6 | LOCAL7  

Défaut :
    AUTH
``````

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
        Nécessite la présence de MemcacheD. Voir bloc [CACHE.MEMCACHE]

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

### **check :**

Timeout en secondes de la vérification de présence d'un serveur en écoute sur le port et l'host renseigné.  
Ce contrôle est effectué durant la phase de configuration.

    1 :
        Valeur mini
    
    3 :
        Valeur par défaut

    3600 :
        Valeur max

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

## **[CACHE.MEMCACHE]**

### **host :**

Adresse du serveur memcacheD

    défaut : 127.0.0.1

### **port :**

Port d'écoute du serveur memacacheD.

    défaut : 11211

### **timeout :**

Durée maximum en secondes d'une transaction vers le serveur de cache.

    0 : 
        Pas de timeout

    3 :
        Valeur par défaut

    60 :
        Valeur max autorisée


## **[CACHE.REDIS]**

### **host :**

Adresse du serveur Redis.

    défaut : 127.0.0.1

### **port :**

Port d'écoute du serveur Redis.

    défaut : 6379

### **timeout :**

Durée maximum en secondes d'une transaction vers le serveur de cache.

    0 : 
        Pas de timeout

    3 :
        Valeur par défaut

    60 :
        Valeur max autorisée

### **db :**

Numéro de la database utilisée pour l'isolation.

    défaut : 0

### **user :**

Utilisateur d'accès

    défaut : 

### **password :**

Mot de passe d'accès

    défaut : 

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
| virtdom | bool | true |

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

### **virtdom :**

Si virtdom est **true** alors l'authentifcation sera par défaut en user@dom si dom est définir, sinon user.  
Si virtdom est **false** l'authentification sera effectuée avec user.

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

Si virtdom est **true** alors l'authentifcation sera par défaut en user@dom si dom est définir, sinon user.  
Si virtdom est **false** l'authentification sera effectuée avec user.

## **[PLUGIN.LEMON]**

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

## **[PLUGIN.PGAUTH]**

| CLEF | TYPE | DÉFAUT |
|:----:|:----:|:-------:|
| host | string | 127.0.0.1 |
| port | int | 5432 |
| bdd | string | |
| user | string | |
| passwd | string | |
| timeout | int | 5 |
| virtdom | bool | false |
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

### **virtdom :**

Si **true** l'identifiant se présentant en uid@domain devra exister en BDD.  
Si **false**, même si l'utilisateur se présente en uid@domain, l'authentification se fera avec uid uniquement.

### **sql :**

Reqête SQL permettant d'obtenir l'username et le password correspond à l'identifiant.  
**$1** correspond à l'identifiant et est obligatoire.

```sql
SELECT username, password FROM users WHERE username LIKE $1 LIMIT 1
```
