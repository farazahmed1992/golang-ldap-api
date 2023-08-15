package ldapfunc

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
)

//"golang.org/x/text/encoding/unicode"

// GetDN Gets and object CN and returns Active Directory DN
func GetDN(dn string, object string) (data string) {
	//Search starts on a root of the directory.
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(dn))
	if object == "group" {
		baseDN := fmt.Sprintf("OU=live,OU=b2b,DC=%s,DC=%s", Config.DomainFirst, Config.DomainLast)
		searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{}, []ldap.Control{})
		result, err := ConnLdap.Search(searchReq)
		if err != nil {
			log.Println(err)
		}
		data = "not match"
		//fmt.Println((result.Entries)[0].GetAttributeValues("objectClass")[1])
		if len(result.Entries) >= 1 {
			for _, entry := range result.Entries {
				//entry.GetAttributeValue("distinguishedName")
				data = entry.GetAttributeValue("distinguishedName")
			}
		}

	}
	if object == "user" {
		baseDN := fmt.Sprintf("DC=%s,DC=%s", Config.DomainFirst, Config.DomainLast)
		searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"distinguishedName"}, []ldap.Control{})
		result, err := ConnLdap.Search(searchReq)
		if err != nil {
			log.Println(err)
		}
		data = "not match"
		//fmt.Println((result.Entries)[0].GetAttributeValues("objectClass")[1])
		if len(result.Entries) >= 1 {
			for _, entry := range result.Entries {
				//entry.GetAttributeValue("distinguishedName")
				data = entry.GetAttributeValue("distinguishedName")
			}
		}

	}

	return
}

func GetMember(user string, dn string) (userexist bool) {

	searchReq := ldap.NewSearchRequest(dn, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{}, nil)
	result, err := ConnLdap.Search(searchReq)
	vals := result.Entries[0].GetAttributeValues("member")
	members := make([]string, len(vals))
	for i, dn := range vals {
		members[i] = dn[strings.Index(dn, "=")+1 : strings.Index(dn, ",")]
	}
	//log.Println(members)
	//fmt.Println(vals[1])
	userexist = false
	for _, checkuser := range members {
		if user == checkuser {
			log.Printf("User:%s Exist in the Group", checkuser)
			userexist = true
			break
		}
	}

	if !userexist {
		log.Printf("User:%s doesnot exist in the Group", user)
	}
	//og.Println("Members: %s", vals)
	//log.Println(result.Entries)
	if err != nil {
		log.Println(err)
	}

	return
}

func AddNewUser(user string, password string, ou string) (string, string) {
	bind_string := LdapBind()
	if bind_string == "UTC-LDAP" || bind_string == "UTC-BIND" {
		data := "Unable to Connect to LDAP"
		status := "Error"
		return data, status
	}

	// Create a disable User
	userDN := fmt.Sprintf("CN=%s,OU=%s,dc=%s,dc=%s", user, ou, Config.DomainFirst, Config.DomainLast)
	fmt.Println(userDN)
	addReq := ldap.NewAddRequest(userDN, []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "organizationalPerson", "user", "person"})
	addReq.Attribute("name", []string{user})
	addReq.Attribute("sAMAccountName", []string{user})
	addReq.Attribute("userAccountControl", []string{fmt.Sprintf("%d", 0x0202)})
	addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	addReq.Attribute("userPrincipalName", []string{fmt.Sprintf("%s@%s.%s", user, Config.DomainFirst, Config.DomainLast)})
	addReq.Attribute("accountExpires", []string{fmt.Sprintf("%d", 0x00000000)})
	if err := ConnLdap.Add(addReq); err != nil {
		ConnLdap.Close()
		return err.Error(), "failed"
	}

	// Modify the Password
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	pwdEncoded, err := utf16.NewEncoder().String(fmt.Sprintf("%q", password))
	if err != nil {
		ConnLdap.Close()
		return err.Error(), "failed"

	}
	modReq := ldap.NewModifyRequest(userDN, []ldap.Control{})
	modReq.Replace("unicodePwd", []string{pwdEncoded})
	if err := ConnLdap.Modify(modReq); err != nil {
		ConnLdap.Close()
		return err.Error(), "failed"
	}
	// Enable the User
	modReq.Replace("userAccountControl", []string{fmt.Sprintf("%d", 0x0200)})

	if err := ConnLdap.Modify(modReq); err != nil {
		ConnLdap.Close()
		return err.Error(), "failed"
	}
	ConnLdap.Close()
	return "User Created Successful", "successful"
}

// AddUsertoGroup Gets a User CN name and add it to the provide group.
func AddUsertoGroup(user string, group string) (data string, status string) {
	log.Println("TESTLOG: ldapfunc calling AddUsertoGroup user:%s, group=%s", user, group)
	bind_string := LdapBind()
	if bind_string == "UTC-LDAP" || bind_string == "UTC-BIND" {
		data = "Unable to Connect to LDAP"
		status = "Error"
		return
	}
	userDN := GetDN(user, "user")
	groupDN := GetDN(group, "group")
	if groupDN != "not match" {
		if userDN != "not match" {

			modify := ldap.NewModifyRequest(groupDN, nil)
			modify.Add("member", []string{userDN})
			for i := 1; i <= 3; i++ {
				err := ConnLdap.Modify(modify)
				usercheck := GetMember(user, groupDN)
				if err != nil {
					log.Println(err)
					ConnLdap.Close()
					if strings.Contains(err.Error(), "Entry Already Exists") {
						data = "User Added Successful"
						status = "successful"
						return
					}
					data = err.Error()
					status = "Error"
					return
				}

				if usercheck {
					log.Printf("User:%s added to the Group: %s", user, group)
					data = "User Added Successful"
					status = "successful"
					break
				} else {
					log.Printf("User: %s Unable to be added to the Group: %s", user, group)
					data = "Unable to add user in the Group"
					status = "failed"
				}
			}
		} else {
			log.Printf("User: %s does not exist in the Directory", user)
			data = "User does not exist in the Directory"
			status = "failed"
		}
	} else {
		log.Printf("Group: %s does not exist in the Directory", group)
		data = "Group does not exist in the Directory"
		status = "failed"
	}
	ConnLdap.Close()
	return
}

func CheckUsertoGroup(user string, group string) (data string, status string) {
	//log.Println("TESTLOG: ldapfunc calling AddUsertoGroup user:%s, group=%s", user, group)
	bind_string := LdapBind()
	if bind_string == "UTC-LDAP" || bind_string == "UTC-BIND" {
		data = "Unable to Connect to LDAP"
		status = "Error"
		return
	}
	groupDN := GetDN(group, "group")
	if groupDN != "not match" {
		members := GetMember(user, groupDN)
		//log.Println(members)
		ConnLdap.Close()
		if members {
			data = "true"
			status = "ok"
			return
		} else {
			data = "false"
			status = "error"
			return
		}
	} else {
		log.Printf("Group: %s does not exist in the Directory", group)
		data = "Group does not exist in the Directory"
		status = "failed"
		return
	}
}

func AddNewGroup(group string, ou string) (string, string) {
	bind_string := LdapBind()
	if bind_string == "UTC-LDAP" || bind_string == "UTC-BIND" {
		data := "Unable to Connect to LDAP"
		status := "Error"
		return data, status
	}
	groupDN := fmt.Sprintf("CN=%s,OU=%s,OU=b2b,DC=%s,DC=%s", group, ou, Config.DomainFirst, Config.DomainLast)
	addReq := ldap.NewAddRequest(groupDN, []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "group"})
	addReq.Attribute("name", []string{group})
	addReq.Attribute("sAMAccountName", []string{group})
	addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	addReq.Attribute("groupType", []string{fmt.Sprintf("%d", 0x00000004|0x80000000)})

	if err := ConnLdap.Add(addReq); err != nil {
		ConnLdap.Close()
		return err.Error(), "failed"
	}
	ConnLdap.Close()
	return "Group Created Successful", "successful"
}
