package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"go-ldap-api/ldapfunc"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
)

// swagger:model UserCreate
type UserCreate struct {
	// The request body
	// in: body
	// username for account creation.
	// required: true
	Username string `json:"username" validate:"required,excludesall=!_#? ,max=20"`
	// password for account.
	// required: true
	Password string `json:"password" validate:"required,excludesall= ,min=10,max=45"`
	// OU name for account to be placed.
	// required: true
	OU string `json:"ou" validate:"required,max=64"`
}

// swagger:model ModifyGroup
type ModifyGroup struct {
	// the GroupDn is the CN name for the Group.
	// required: true
	GroupDn string `json:"groupDN" validate:"required,excludesall=!_#? ,max=45"`
	// the UserDN is the CN name for the user.
	// required: true
	UserDN string `json:"userDN" validate:"required,excludesall=!_#? ,max=20"`
	// the Org is the organization name.
	// required: true
	Org string `json:"org" validate:"required,excludesall=!_#? ,max=15"`
}

// swagger:model UserCheck
type UserCheck struct {
	// the GroupDn is the CN name for the Group.
	// required: true
	GroupDn string `json:"groupDN" validate:"required,excludesall=!_#? ,max=45"`
	// the UserDN is the CN name for the user.
	// required: true
	UserDN string `json:"userDN" validate:"required,excludesall=!_#? ,max=20"`
	// the Org is the organization name.
	// required: true
	Org string `json:"org" validate:"required,excludesall=!_#? ,max=15"`
}

// swagger:model BulkModifyGroup
type BulkModifyGroup struct {
	// the GroupDn is the CN name for the Group.
	// required: true
	GroupDn string `json:"groupDN" validate:"required,excludesall=!_#? ,max=45"`
	// the UserDN is the CN name for the user.
	// required: true
	UserDN []string `json:"userDN" validate:"required,dive,excludesall=!_#? ,max=20"`
	// the Org is the organization name.
	// required: true
	Org string `json:"org" validate:"required,excludesall=!_#? ,max=15"`
}

// swagger:model GroupCreate
type GroupCreate struct {
	// The request body
	// in: body
	// Group name for creation.
	// required: true
	Group string `json:"group" validate:"required,excludesall=!_#? ,max=45"`
	// the Org is the organization name.
	// required: true
	Org string `json:"org" validate:"required,excludesall=!_#? ,max=15"`
	// OU name for group to be placed in the active directory.
	// required: true
	OU string `json:"ou" validate:"required,excludesall=!_#? ,max=15"`
}

// swagger:model GroupDelete
type GroupDelete struct {
	// Group name for creation.
	// required: true
	Group string `json:"group" validate:"required,excludesall=!_#? ,max=45"`
	// the Org is the organization name.
	// required: true
	Org string `json:"org" validate:"required,excludesall=!_#? ,max=15"`
}

func (p *ModifyGroup) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (p *UserCheck) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (p *UserCreate) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (p *GroupCreate) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (p *GroupDelete) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (p *BulkModifyGroup) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func AddUserToOU(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data UserCreate
	_ = json.NewDecoder(r.Body).Decode(&data)
	err := data.Validate()
	if err != nil {
		w.WriteHeader(401)
		hclog.Default().Error("validation error", err)
		payload := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": err.Error()})
		if payload != nil {

			return
		}
		return
	}
	reqData, status := ldapfunc.AddNewUser(data.Username, data.Password, data.OU)
	if status == "failed" {
		w.WriteHeader(400)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}
	} else if reqData == "User Created Successful" {
		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}

	} else if status == "Error" {
		w.WriteHeader(405)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": reqData})
		if err != nil {
			return
		}

	} else {
		w.WriteHeader(404)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": "Unknown Error"})
		if err != nil {
			return
		}
	}
}

func AddToGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data ModifyGroup
	_ = json.NewDecoder(r.Body).Decode(&data)
	err := data.Validate()
	if err != nil {
		w.WriteHeader(401)
		hclog.Default().Error("validation error", err)
		payload := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": err.Error()})
		if payload != nil {

			return
		}
		return
	}
	log.Printf("AddToGroup: Adding user:%s to group:%s and org:%s", data.UserDN, data.GroupDn, data.Org)
	reqData, status := ldapfunc.AddUsertoGroup(data.UserDN, data.Org+"-"+data.GroupDn)
	if reqData == "Unable to add user in the Group" {
		w.WriteHeader(400)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}
	} else if reqData == "User Added Successful" {
		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}

	} else if reqData == "User does not exist in the Directory" {
		w.WriteHeader(402)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}

	} else if reqData == "Group does not exist in the Directory" {
		w.WriteHeader(403)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}

	} else if status == "Error" {
		w.WriteHeader(405)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": reqData})
		if err != nil {
			return
		}

	} else {
		w.WriteHeader(404)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": "Unknown Error"})
		if err != nil {
			return
		}
	}

}

func CheckToGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data ModifyGroup
	_ = json.NewDecoder(r.Body).Decode(&data)
	err := data.Validate()
	if err != nil {
		w.WriteHeader(401)
		hclog.Default().Error("validation error", err)
		payload := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": err.Error()})
		if payload != nil {

			return
		}
		return
	}
	log.Printf("CheckUserToGroup: Check user:%s to group:%s and org:%s", data.UserDN, data.GroupDn, data.Org)
	reqData, status := ldapfunc.CheckUsertoGroup(data.UserDN, data.Org+"-"+data.GroupDn)
	if reqData == "false" {
		w.WriteHeader(400)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": "User is not Present in defined Group"})
		if err != nil {
			return
		}
	} else if reqData == "true" {
		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": "User is Present is the defined Group"})
		if err != nil {
			return
		}

	} else if reqData == "Group does not exist in the Directory" {
		w.WriteHeader(403)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}

	} else if status == "Error" {
		w.WriteHeader(405)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": reqData})
		if err != nil {
			return
		}

	} else {
		w.WriteHeader(404)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": "Unknown Error"})
		if err != nil {
			return
		}
	}
}

func AddGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data GroupCreate
	_ = json.NewDecoder(r.Body).Decode(&data)
	err := data.Validate()
	if err != nil {
		w.WriteHeader(401)
		hclog.Default().Error("validation error", err)
		payload := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": err.Error()})
		if payload != nil {

			return
		}
		return
	}
	log.Printf("AddGroup: Adding group:%s and org:%s", data.Group, data.Org)
	reqData, status := ldapfunc.AddNewGroup(data.Org+"-"+data.Group, data.OU)
	if status == "failed" {
		w.WriteHeader(400)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}
	} else if status == "Error" {
		w.WriteHeader(405)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": reqData})
		if err != nil {
			return
		}

	} else {
		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}

	}

}

func RemoveGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data GroupDelete
	_ = json.NewDecoder(r.Body).Decode(&data)
	err := data.Validate()
	if err != nil {
		w.WriteHeader(401)
		hclog.Default().Error("validation error", err)
		payload := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": err.Error()})
		if payload != nil {

			return
		}
		return
	}
	log.Printf("RemoveGroup: Removing group:%s and org:%s", data.Group, data.Org)
	reqData, status := ldapfunc.DelGroup(data.Org + "-" + data.Group)
	if reqData == "Group does not exist in the Directory" {
		w.WriteHeader(403)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}
	} else if reqData == "Team Group has been removed" {
		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}

	} else if status == "Error" {
		w.WriteHeader(405)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": reqData})
		if err != nil {
			return
		}

	} else {
		w.WriteHeader(404)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": "Unknown Error"})
		if err != nil {
			return
		}
	}

}

func RemoveFromGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data ModifyGroup
	_ = json.NewDecoder(r.Body).Decode(&data)
	err := data.Validate()
	if err != nil {
		w.WriteHeader(401)
		hclog.Default().Error("validation error", err)
		payload := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": err.Error()})
		if payload != nil {

			return
		}
		return
	}
	log.Printf("RemoveFromGroup: Removing user:%s from group:%s and org:%s", data.UserDN, data.GroupDn, data.Org)
	reqData, status := ldapfunc.RemoveUserFromGroup(data.UserDN, data.Org+"-"+data.GroupDn)
	if reqData == "User does not exist in the Directory" {
		w.WriteHeader(402)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}
	} else if reqData == "User removed from the Group" {
		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}

	} else if reqData == "Group does not exist in the Directory" {
		w.WriteHeader(403)
		err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
		if err != nil {
			return
		}
	} else if status == "Error" {
		w.WriteHeader(405)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": reqData})
		if err != nil {
			return
		}

	} else {
		w.WriteHeader(404)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": "Unknown Error"})
		if err != nil {
			return
		}
	}

}

func AddBulkToGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data BulkModifyGroup
	_ = json.NewDecoder(r.Body).Decode(&data)
	err := data.Validate()
	if err != nil {
		w.WriteHeader(401)
		hclog.Default().Error("validation error", err)
		payload := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": err.Error()})
		if payload != nil {

			return
		}
		return
	}
	var passed []string
	var failed []string
	var response []string
	for i := 0; i < len(data.UserDN); i++ {
		log.Printf("AddBulkToGroup: Adding user:%s to group:%s and org:%s", data.UserDN[i], data.GroupDn, data.Org)
		res, status := ldapfunc.AddUsertoGroup(data.UserDN[i], data.Org+"-"+data.GroupDn)
		if status == "failed" {
			failed = append(failed, data.UserDN[i])
			response = append(response, res)
		} else {
			passed = append(passed, data.UserDN[i])
			response = append(response, res)
		}

	}

	payload := map[string]interface{}{
		"status": "ok",
		"data": map[string]interface{}{"passed": passed,
			"failed":   failed,
			"response": response,
		},
		"passed": len(passed),
		"failed": len(failed),
	}
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		return
	}

}

func RemoveBulkFromGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data BulkModifyGroup
	_ = json.NewDecoder(r.Body).Decode(&data)
	err := data.Validate()
	if err != nil {
		w.WriteHeader(401)
		hclog.Default().Error("validation error", err)
		payload := json.NewEncoder(w).Encode(map[string]string{"status": "failed", "data": err.Error()})
		if payload != nil {

			return
		}
		return
	}
	for i := 0; i < len(data.UserDN); i++ {
		log.Printf("RemoveBulkFromGroup: Remove user:%s from group:%s and org:%s", data.UserDN[i], data.GroupDn, data.Org)
		reqData, status := ldapfunc.RemoveUserFromGroup(data.UserDN[i], data.Org+"-"+data.GroupDn)
		if status == "failed" {
			w.WriteHeader(400)
			err := json.NewEncoder(w).Encode(map[string]string{"status": status, "data": reqData})
			if err != nil {
				return
			}
		}

	}
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "data": "users removed from Team"})
	if err != nil {
		return
	}

}
