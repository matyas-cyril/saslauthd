package ldap

import (
	"crypto/tls"
	"fmt"
	"reflect"
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
}

type Ldap struct {
	opt LdapOpt
	cnx *myLdap.Conn
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
		Port:               389,
		Timeout:            10,
		Tls:                false,
		InsecureSkipVerify: true,
	}

	for k, v := range args {

		switch k {
		case "uri", "admin", "pwd", "baseDN":

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
			}

		case "filter":

			kV, kErr := v.(string)
			if !kErr {
				return nil, fmt.Errorf("ldap param key '%s' failed to typecast", k)
			}
			l.Filter = strings.TrimSpace(kV)

		case "port", "timeout":

			typeTarget := reflect.TypeFor[int]()
			rv := reflect.ValueOf(v)
			if !rv.Type().AssignableTo(typeTarget) {
				return nil, fmt.Errorf("ldap param key '%s' failed to typecast", k)
			}

			nbr := rv.Convert(typeTarget).Int()
			if nbr < 0 || nbr > 65535 {
				return nil, fmt.Errorf("ldap param key '%s' integer range invalid", k)
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

		case "tls", "tlsSkipVerify":
			kV, kErr := v.(bool)
			if !kErr {
				return nil, fmt.Errorf("ldap param key '%s' failed to typecast", k)
			}

			switch k {
			case "tls":
				l.Tls = kV

			case "tlsSkipVerify":
				l.InsecureSkipVerify = kV
			}

		default:
			return nil, fmt.Errorf("ldap param key '%s' not exist", k)
		}

	}

	return &Ldap{opt: l}, nil

}

func (l *Ldap) Connect() (err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("panic error plugin ldap : %s", pErr)
		}
	}()

	// Définir le timeout de connexion
	myLdap.DefaultTimeout = time.Duration(l.opt.Timeout) * time.Second

	if !l.opt.Tls {
		l.cnx, err = myLdap.DialURL(fmt.Sprintf("ldap://%s:%d", l.opt.Uri, l.opt.Port))
	} else {
		l.cnx, err = myLdap.DialURL(fmt.Sprintf("ldaps://%s:%d", l.opt.Uri, l.opt.Port), myLdap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: l.opt.InsecureSkipVerify}))
	}

	if err != nil {
		l.cnx = nil
		return err
	}

	l.cnx.SetTimeout(time.Duration(l.opt.Timeout) * time.Second)

	// On Bind avec le compte appli
	err = l.cnx.Bind(l.opt.Admin, l.opt.Passwd)
	if err != nil {
		l.cnx = nil
		return err
	}

	return nil
}

// searchUser : Permet de vérifier qu'un utilisateur existe
func (l *Ldap) searchUser(userName string) error {

	filter := fmt.Sprintf(l.opt.Filter, userName)
	searchRequest := myLdap.NewSearchRequest(l.opt.BaseDn, myLdap.ScopeWholeSubtree, myLdap.DerefAlways,
		0, 0, false,
		filter,
		[]string{"dn"},
		nil)

	sr, err := l.cnx.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) != 1 {
		return fmt.Errorf("user '%s' does not exist or too many entries returned", userName)
	}
	return nil
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

	if len(userName) == 0 {
		return fmt.Errorf("ldap auth password empty")
	}

	// On vérifie que l'utilisateur est authorisé/existe via le filtre
	if err := l.searchUser(userName); err != nil {
		return err
	}

	// On génére l'identifiant utilisateur en type LDAP
	dn := fmt.Sprintf("uid=%s,%s", userName, l.opt.BaseDn)

	// On Bind avec le compte de l'utilisateur pour contrôler le mot de passe
	bindRequest := myLdap.NewSimpleBindRequest(dn, passwd, nil)
	_, err = l.cnx.SimpleBind(bindRequest)
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

	if l.cnx != nil && !l.cnx.IsClosing() {
		l.cnx.Close()
		l.cnx = nil
		return true, nil
	}
	return false, nil
}
