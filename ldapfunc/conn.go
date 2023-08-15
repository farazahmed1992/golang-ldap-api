package ldapfunc

import (
	"crypto/tls"
	"fmt"
	"time"

	"go-ldap-api/util"

	"github.com/go-ldap/ldap/v3"
)

var Config, _ = util.LoadConfig(".")
var ldapURL = Config.LdapServer
var ldapURL2 = Config.LdapServer2
var ConnLdap, _ = ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: false}))
var Bind_return string

func LdapBind() string {

	bindUserDN := fmt.Sprintf("cn=%s, ou=IT-Admins, dc=%s,dc=%s", Config.LdapUser, Config.DomainFirst, Config.DomainLast)
	err := ConnLdap.Bind(bindUserDN, Config.LdapPassword)
	if err != nil {
		ConnLdap, err = ldap.DialURL(ldapURL2, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: false}))
		if err != nil {
			fmt.Println(err.Error())
			Bind_return = "UTC-LDAP"
			return Bind_return
		}

		err := ConnLdap.Bind(bindUserDN, Config.LdapPassword)
		if err != nil {
			fmt.Println(err.Error())
			Bind_return = "UTC-BIND"
			return Bind_return

		}
	}
	Bind_return = "LDAP-Connected"
	ConnLdap.SetTimeout(10 * time.Second)
	return Bind_return
}
