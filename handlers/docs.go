package handlers

// A HealthCheck is an success response that is used to verify if the server is up.
// swagger:response healthcheckok
type HealthCheck struct {
	// The success message
	// in: body
	Body struct {

		// Example: ok
		Status string `json:"status"`
		// Example: Api Server is Running
		Data string `json:"data"`
	}
}

// swagger:route GET / health-checks
//
// Health check.
//
// This will show if the server is up.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Responses:
//       200: healthcheckok

// LdapApiDoc for creating groups and adding user to the groups
// swagger:model AddNewGroupRes
type AddNewGroupRes struct {
	// The success message
	// in: body

	// Example: ok
	Status string `json:"status"`
	// Example: Add/Remove a Team/User
	Data string `json:"data"`
}

// swagger:model AddNewGroupErrRes
type AddNewGroupErrRes struct {
	// The success message
	// in: body

	// Example: failed
	Status string `json:"status"`
	// Example: LDAP connection issues or LDAP unable to created user invalid group or ou
	Data string `json:"data"`
}

// swagger:model UserCheckRes
type UserCheckRes struct {
	// The success message
	// in: body

	// Example: ok
	Status string `json:"status"`
	// Example: User is Present in defined the Group
	Data string `json:"data"`
}

// swagger:model UserCheckErrRes
type UserCheckErrRes struct {
	// The success message
	// in: body

	// Example: failed
	Status string `json:"status"`
	// Example: User is not Present in defined the Group
	Data string `json:"data"`
}

// swagger:operation POST /api/add/group Group GroupCreate
//
// ADD new group.
// ---
// produces:
// - application/json
// parameters:
// - name: server_id
//   in: body
//   schema:
//	   "$ref": "#/definitions/GroupCreate"
//
// responses:
//   '200':
//     description: A success response that verifies the creation of the new group in AD.
//     schema:
//       "$ref": "#/definitions/AddNewGroupRes"
//   '400':
//     description: A failure response with error message.
//     schema:
//       "$ref": "#/definitions/AddNewGroupErrRes"

// swagger:operation POST /api/remove/group Group GroupDelete
//
// Remove group.
// ---
// produces:
// - application/json
// parameters:
// - name: server_id
//   in: body
//   schema:
//	   "$ref": "#/definitions/GroupDelete"
//
// responses:
//   '200':
//     description: A success response that verifies the creation of the new group in AD.
//     schema:
//       "$ref": "#/definitions/AddNewGroupRes"
//   '400':
//     description: A failure response with error message.
//     schema:
//       "$ref": "#/definitions/AddNewGroupErrRes"

// swagger:operation POST /api/check/usertogroup Group UserToGroup
//
// Check user exist in defined group.
// ---
// produces:
// - application/json
// parameters:
// - name: server_id
//   in: body
//   schema:
//	   "$ref": "#/definitions/UserCheck"
//
// responses:
//   '200':
//     description: A success response that verifies the creation of the new group in AD.
//     schema:
//       "$ref": "#/definitions/UserCheckRes"
//   '400':
//     description: A failure response with error message.
//     schema:
//       "$ref": "#/definitions/UserCheckErrRes"

// swagger:operation POST /api/add/usertogroup Group UserToGroup
//
// ADD user to a group.
// ---
// produces:
// - application/json
// parameters:
// - name: server_id
//   in: body
//   schema:
//	   "$ref": "#/definitions/ModifyGroup"
//
// responses:
//   '200':
//     description: A success response that verifies the creation of the new group in AD.
//     schema:
//       "$ref": "#/definitions/AddNewGroupRes"
//   '400':
//     description: A failure response with error message.
//     schema:
//       "$ref": "#/definitions/AddNewGroupErrRes"

// swagger:operation POST /api/remove/userfromgroup Group RemoveUserFromGroup
//
// Remove user from a group.
// ---
// produces:
// - application/json
// parameters:
// - name: server_id
//   in: body
//   schema:
//	   "$ref": "#/definitions/ModifyGroup"
//
// responses:
//   '200':
//     description: A success response that verifies the creation of the new group in AD.
//     schema:
//       "$ref": "#/definitions/AddNewGroupRes"
//   '400':
//     description: A failure response with error message.
//     schema:
//       "$ref": "#/definitions/AddNewGroupErrRes"

// swagger:operation POST /api/add/bulkusertogroup BulkGroup AddUsersToGroup
//
// Add users from a group.
// ---
// produces:
// - application/json
// parameters:
// - name: server_id
//   in: body
//   schema:
//	   "$ref": "#/definitions/BulkModifyGroup"
//
// responses:
//   '200':
//     description: A success response that verifies the creation of the new group in AD.
//     schema:
//       "$ref": "#/definitions/AddNewGroupRes"
//   '400':
//     description: A failure response with error message.
//     schema:
//       "$ref": "#/definitions/AddNewGroupErrRes"

// swagger:operation POST /api/remove/bulkusersfromgroup BulkGroup RemoveUsersToGroup
//
// Remove users from a group.
// ---
// produces:
// - application/json
// parameters:
// - name: server_id
//   in: body
//   schema:
//	   "$ref": "#/definitions/BulkModifyGroup"
//
// responses:
//   '200':
//     description: A success response that verifies the creation of the new group in AD.
//     schema:
//       "$ref": "#/definitions/AddNewGroupRes"
//   '400':
//     description: A failure response with error message.
//     schema:
//       "$ref": "#/definitions/AddNewGroupErrRes"
