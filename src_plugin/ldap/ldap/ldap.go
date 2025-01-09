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

func New(args *map[string]any) (_rtnOpt *Ldap, _err error) {

	defer func() {
		if err := recover(); err != nil {
			_rtnOpt = nil
			_err = fmt.Errorf(fmt.Sprintf("fatal: failed to initialize ldap - %s", err))
		}
	}()

	if args == nil {
		return nil, fmt.Errorf("no args to initialize ldap connection")
	}

	l := LdapOpt{
		Port:               389,
		Timeout:            10,
		InsecureSkipVerify: true,
	}

	for k, v := range *args {

		switch k {
		case "uri", "admin", "pwd", "baseDN":
			if !reflect.ValueOf(string("")).Type().ConvertibleTo(reflect.ValueOf(v).Type()) {
				return nil, fmt.Errorf(fmt.Sprintf("ldap param key '%s' must be a string", k))
			}

			d := strings.TrimSpace(v.(string))
			if len(d) == 0 {
				return nil, fmt.Errorf(fmt.Sprintf("ldap param key '%s' must be a not empty string", k))
			}

			switch k {
			case "uri":
				l.Uri = d
			case "admin":
				l.Admin = d
			case "pwd":
				l.Passwd = d
			case "baseDN":
				l.BaseDn = d
			}

		case "filter":
			if !reflect.ValueOf(string("")).Type().ConvertibleTo(reflect.ValueOf(v).Type()) {
				return nil, fmt.Errorf(fmt.Sprintf("ldap param key '%s' must be a string", k))
			}

			l.Filter = strings.TrimSpace(v.(string))

		case "port", "timeout":
			if !reflect.ValueOf(int(0)).Type().ConvertibleTo(reflect.ValueOf(v).Type()) {
				return nil, fmt.Errorf(fmt.Sprintf("ldap param key '%s' must be an integer", k))
			}

			d := v.(int)

			switch k {
			case "port":
				if d < 1 || d > 65535 {
					return nil, fmt.Errorf(fmt.Sprintf("ldap param key '%s' must be an integer between 1 and 65535", k))
				}
				l.Port = uint16(d)

			case "timeout":
				if d < 0 || d > 3600 {
					return nil, fmt.Errorf(fmt.Sprintf("ldap param key '%s' must be an integer between 0 and 3600", k))
				}
				l.Timeout = uint16(d)
			}

		case "tls", "tlsSkipVerify":
			if !reflect.ValueOf(l.Tls).Type().ConvertibleTo(reflect.ValueOf(v).Type()) {
				return nil, fmt.Errorf(fmt.Sprintf("ldap param key '%s' must be a boolean", k))
			}
			switch k {
			case "tls":
				l.Tls = v.(bool)

			case "tlsSkipVerify":
				l.InsecureSkipVerify = v.(bool)
			}

		default:
			return nil, fmt.Errorf(fmt.Sprintf("arg '%s' not exist", k))
		}

	}

	return &Ldap{opt: l}, nil

}

func (l *Ldap) Connect() (_err error) {

	defer func() {
		if err := recover(); err != nil {
			_err = fmt.Errorf(fmt.Sprintf("fatal: failed to initialize ldap connection - %s", err))
		}
	}()

	// Définir le timeout de connexion
	myLdap.DefaultTimeout = time.Duration(l.opt.Timeout) * time.Second

	var err error

	if !l.opt.Tls {
		l.cnx, err = myLdap.Dial("tcp", fmt.Sprintf("%s:%d", l.opt.Uri, l.opt.Port))
	} else {
		l.cnx, err = myLdap.DialTLS("tcp", fmt.Sprintf("%s:%d", l.opt.Uri, l.opt.Port), &tls.Config{InsecureSkipVerify: l.opt.InsecureSkipVerify})
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
func (l *Ldap) Auth(userName, passwd string) (_err error) {

	defer func() {
		if err := recover(); err != nil {
			_err = fmt.Errorf(fmt.Sprintf("fatal: failed to initialize ldap Auth - %s", err))
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
	_, err := l.cnx.SimpleBind(bindRequest)
	if err != nil {
		return err
	}

	return nil
}

// Close : Fermer les connexions
func (l *Ldap) Close() (_bool bool, _err error) {

	defer func() {
		if err := recover(); err != nil {
			_err = fmt.Errorf(fmt.Sprintf("fatal: %s", err))
			_bool = false
		}
	}()

	if l.cnx != nil && !l.cnx.IsClosing() {
		l.cnx.Close()
		l.cnx = nil
		return true, nil
	}
	return false, nil
}
