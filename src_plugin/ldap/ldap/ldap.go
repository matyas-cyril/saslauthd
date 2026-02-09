package ldap

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	myLdap "github.com/go-ldap/ldap/v3"
)

// Ldap : Structure
type LdapOpt struct {
	Uri                string
	Admin              string
	Passwd             string
	BaseDn             string
	Filter             string
	Port               uint16
	Timeout            uint16
	Tls                bool
	InsecureSkipVerify bool
	Attribute          string
	AttributeMatch     string
	VirtDom            bool
}

type Ldap struct {
	Opt LdapOpt
	Cnx *myLdap.Conn
}

func New(args map[string]any) (ldap *Ldap, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			ldap = nil
			err = fmt.Errorf("panic error plugin ldap : %s", pErr)
		}
	}()

	if args == nil {
		return nil, fmt.Errorf("no args to initialize ldap connection")
	}

	// Valeurs par défaut
	l := LdapOpt{
		Uri:                "127.0.0.1",
		Port:               389,
		Timeout:            10,
		Filter:             "(uid=%s)",
		Attribute:          "dn",
		AttributeMatch:     "uid",
		Tls:                false,
		InsecureSkipVerify: true,
		VirtDom:            true,
	}

	for k, v := range args {

		switch k {
		case "uri", "admin", "pwd", "baseDN", "att", "attMatch":

			kV, kErr := v.(string)
			if !kErr {
				return nil, fmt.Errorf("ldap param key '%s' failed to typecast", k)
			}
			kV = strings.TrimSpace(kV)
			if len(kV) == 0 {
				return nil, fmt.Errorf("ldap param key '%s' defined but empty", k)
			}

			switch k {

			case "uri":
				l.Uri = kV
			case "admin":
				l.Admin = kV
			case "pwd":
				l.Passwd = kV
			case "baseDN":
				l.BaseDn = kV
			case "att":
				l.Attribute = kV
			case "attMatch":
				l.AttributeMatch = kV
			}

		case "filter":

			kV, kErr := v.(string)
			if !kErr {
				return nil, fmt.Errorf("ldap param key '%s' failed to typecast", k)
			}
			kV = strings.TrimSpace(kV)
			if !strings.HasPrefix(kV, "(") || !strings.HasSuffix(kV, ")") {
				return nil, fmt.Errorf("ldap param key '%s' syntaxe invalid", k)
			}

			l.Filter = kV

		case "port", "timeout":

			nbr, cast := v.(int64)
			if !cast {
				return nil, fmt.Errorf("ldap param key '%s' failed to typecast", k)
			}

			switch k {
			case "port":
				if nbr < 1 || nbr > 65535 {
					return nil, fmt.Errorf("ldap param key '%s' must be an integer between 1 and 65535", k)
				}
				l.Port = uint16(nbr)

			case "timeout":
				if nbr < 0 || nbr > 3600 {
					return nil, fmt.Errorf("ldap param key '%s' must be an integer between 0 and 3600", k)
				}
				l.Timeout = uint16(nbr)
			}

		case "tls", "tlsSkipVerify", "virtdom":
			kV, kErr := v.(bool)
			if !kErr {
				return nil, fmt.Errorf("ldap param key '%s' failed to typecast", k)
			}

			switch k {
			case "tls":
				l.Tls = kV

			case "tlsSkipVerify":
				l.InsecureSkipVerify = kV

			case "virtdom":
				l.VirtDom = kV
			}

		default:
			return nil, fmt.Errorf("ldap param key '%s' not exist", k)
		}

	}

	return &Ldap{Opt: l}, nil

}

func (l *Ldap) Connect() (err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("panic error plugin ldap : %s", pErr)
		}
	}()

	// Définir le timeout de connexion
	myLdap.DefaultTimeout = time.Duration(l.Opt.Timeout) * time.Second

	if !l.Opt.Tls {
		l.Cnx, err = myLdap.DialURL(fmt.Sprintf("ldap://%s:%d", l.Opt.Uri, l.Opt.Port))
	} else {
		l.Cnx, err = myLdap.DialURL(fmt.Sprintf("ldaps://%s:%d", l.Opt.Uri, l.Opt.Port), myLdap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: l.Opt.InsecureSkipVerify}))
	}

	if err != nil {
		l.Cnx = nil
		return err
	}

	l.Cnx.SetTimeout(time.Duration(l.Opt.Timeout) * time.Second)

	// On Bind avec le compte appli
	err = l.Cnx.Bind(l.Opt.Admin, l.Opt.Passwd)
	if err != nil {
		l.Cnx = nil
		return err
	}

	return nil
}

// searchUser : Permet de vérifier qu'un utilisateur existe et le retourne
func (l *Ldap) searchUser(userName string) (string, error) {

	/*
	   baseDN (string) : point de départ de la recherche (ex : "dc=example,dc=com").
	   scope (int) : portée de la recherche (constantes) :
	       ldap.ScopeBaseObject — seulement l’entrée baseDN.
	       ldap.ScopeSingleLevel — enfants directs de baseDN.
	       ldap.ScopeWholeSubtree — baseDN et tout son sous-arbre.
	   derefAliases (int) : comportement de résolution des alias (constantes) :
	       ldap.NeverDerefAliases, ldap.DerefInSearching, ldap.DerefFindingBaseObj, ldap.DerefAlways
	   sizeLimit (int) : nombre max d’entrées renvoyées (0 = pas de limite côté client).
	   timeLimit (int) : temps max en secondes côté serveur pour exécuter la recherche (0 = pas de limite).
	   typesOnly (bool) : si true, le serveur ne renvoie que les noms d’attributs (pas les valeurs).
	   filter (string) : filtre LDAP en syntaxe standard (ex : "(objectClass=person)").
	   attributes ([]string) : liste des attributs à récupérer (nil ou [] pour tous).
	   controls ([]Control) : contrôles LDAP optionnels (ex : paging control).
	*/
	searchRequest := myLdap.NewSearchRequest(
		l.Opt.BaseDn,
		myLdap.ScopeWholeSubtree,
		myLdap.DerefInSearching,
		3,
		0,
		true,
		fmt.Sprintf(l.Opt.Filter, userName),
		[]string{l.Opt.Attribute},
		nil)

	sr, err := l.Cnx.Search(searchRequest)
	if err != nil {
		return "", err
	}

	if len(sr.Entries) != 1 {
		return "", fmt.Errorf("user '%s' does not exist or too many entries returned", userName)
	}

	dn := sr.Entries[0].DN
	parseDN, err := myLdap.ParseDN(dn)
	if err != nil {
		return "", fmt.Errorf("failed to parse DN '%s': %w", dn, err)
	}

	// Vérifier que l'on a bien l'user et non un user jocker (ex: c*)
	for _, rdn := range parseDN.RDNs {
		for _, at := range rdn.Attributes {
			if strings.ToLower(at.Type) == l.Opt.AttributeMatch {
				if strings.EqualFold(at.Value, userName) {
					return dn, nil
				}
			}
		}
	}

	return "", fmt.Errorf("value for attribute '%s' not found in DN '%s'", l.Opt.AttributeMatch, dn)
}

// Auth :
func (l *Ldap) Auth(userName, passwd string) (err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("panic error plugin ldap Auth : %s", pErr)
		}
	}()

	userName = strings.TrimSpace(userName)
	passwd = strings.TrimSpace(passwd)

	if len(userName) == 0 {
		return fmt.Errorf("ldap auth user name empty")
	}

	if len(passwd) == 0 {
		return fmt.Errorf("ldap auth password empty")
	}

	// On vérifie que l'utilisateur existe via le filtre
	dnUser, err := l.searchUser(userName)
	if err != nil {
		return err
	}

	// On Bind avec le compte de l'utilisateur pour contrôler le mot de passe
	bindRequest := myLdap.NewSimpleBindRequest(dnUser, passwd, nil)
	_, err = l.Cnx.SimpleBind(bindRequest)
	if err != nil {
		return err
	}

	return nil
}

// Close : Fermer les connexions
func (l *Ldap) Close() (status bool, err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("panic error plugin ldap Close : %s", pErr)
			status = false
		}
	}()

	if l.Cnx != nil && !l.Cnx.IsClosing() {
		l.Cnx.Close()
		l.Cnx = nil
		return true, nil
	}
	return false, nil
}
