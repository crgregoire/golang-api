package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/role.go

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/tespo/buddha/db"
	"github.com/tespo/buddha/util"
	"github.com/tespo/satya/v2/types"
)

//
// GetRoles is the GET method for a roles' roles
//
func GetRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = util.SetDBPagination(db, r)

	var roles types.Roles
	if err := roles.Get(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, roles)
}

//
// GetRolesByID is the GET method for a roles' roles by ID
//
func GetRolesByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["role_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var role types.Role
	if err := role.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(role); err != nil {
		panic(err)
	}
}

//
// PostRoles is the POST method for a roles' roles
//
func PostRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var role types.Role
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&role)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := role.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(role); err != nil {
		panic(err)
	}
}

//
// PutRolesByID is the PUT method for a roles' roles by ID
//
func PutRolesByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roleID := mux.Vars(r)["role_id"]
	if roleID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No Role ID supplied"))
		return
	}
	var role types.Role
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&role)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if role.ID.String() != roleID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update role ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := role.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(role); err != nil {
		panic(err)
	}
}

//
// DeleteRolesByID is the DELETE method for a roles' roles by ID
//
func DeleteRolesByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["role_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var role types.Role
	if err := role.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(role); err != nil {
		panic(err)
	}
}

//
// AddPermissionToRoleByID associates a permission to a role using the id
//
func AddPermissionToRoleByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roleID := mux.Vars(r)["role_id"]
	permissionID := mux.Vars(r)["permission_id"]
	if roleID == "" || permissionID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var role types.Role
	if err := role.GetByID(db, uuid.FromStringOrNil(roleID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err := role.AddPermissionToRoleByID(db, uuid.FromStringOrNil(permissionID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(role); err != nil {
		panic(err)
	}
}

//
// DeletePermissionFromRole removes a permission from a role
//
func DeletePermissionFromRole(w http.ResponseWriter, r *http.Request) {
	roleID := mux.Vars(r)["role_id"]
	permissionID := mux.Vars(r)["permission_id"]
	if roleID == "" || permissionID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	permission := types.Permission{}
	if err := permission.GetByID(db, uuid.FromStringOrNil(permissionID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	var role types.Role
	if err := role.GetByID(db, uuid.FromStringOrNil(roleID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err := db.Model(&role).Association("Permissions").Delete(permission).Error; err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, role)
}

//
// GetRoleWithPermissions returns a role with the attached permissions
//
func GetRoleWithPermissions(w http.ResponseWriter, r *http.Request) {
	roleID := mux.Vars(r)["role_id"]
	if roleID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var role types.Role
	if err := role.GetByID(db.Preload("Permissions"), uuid.FromStringOrNil(roleID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, role)
}
