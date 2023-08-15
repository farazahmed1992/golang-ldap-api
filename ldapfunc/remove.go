package ldapfunc

import (
	"log"

	"github.com/go-ldap/ldap/v3"
)

// Gets and object CN and returns Active Directory DN

func DelGroup(group string) (data string, status string) {
	bind_string := LdapBind()
	if bind_string == "UTC-LDAP" || bind_string == "UTC-BIND" {
		data = "Unable to Connect to LDAP"
		status = "Error"
		return
	}
	groupDN := GetDN(group, "group")
	if groupDN != "not match" {
		delReq := ldap.NewDelRequest(groupDN, []ldap.Control{})
		if err := ConnLdap.Del(delReq); err != nil {
			log.Println(err)
			ConnLdap.Close()
			data = err.Error()
			status = "Error"
			return
		}
		ConnLdap.Close()
		data = "Team Group has been removed"
		status = "successful"
	} else {
		log.Printf("Group: %s does not exist in the Directory", group)
		data = "Group does not exist in the Directory"
		status = "failed"
	}
	return
}

func RemoveUserFromGroup(user string, group string) (data string, status string) {
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
			modify.Delete("member", []string{userDN})
			err := ConnLdap.Modify(modify)
			if err != nil {
				log.Println(err)
				ConnLdap.Close()
				data = err.Error()
				status = "Error"
				return
			}
			data = "User removed from the Group"
			status = "successful"

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
