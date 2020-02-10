package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/permission.go

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
// GetPermissions is the GET method for a permissions' permissions
//
func GetPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = util.SetDBPagination(db, r)

	var permissions types.Permissions
	if err := permissions.Get(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, permissions)
}

//
// GetPermissionsByID is the GET method for a permissions' permissions by ID
//
func GetPermissionsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["permission_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var permission types.Permission
	if err := permission.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(permission); err != nil {
		panic(err)
	}
}

//
// PostPermissions is the POST method for a permissions' permissions
//
func PostPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var permission types.Permission
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&permission)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := permission.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(permission); err != nil {
		panic(err)
	}
}

//
// PutPermissionsByID is the PUT method for a permissions' permissions by ID
//
func PutPermissionsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	permissionID := mux.Vars(r)["permission_id"]
	if permissionID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	var permission types.Permission
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&permission)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if permission.ID.String() != permissionID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update permission ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := permission.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(permission); err != nil {
		panic(err)
	}
}

//
// DeletePermissionsByID is the DELETE method for a permissions' permissions by ID
//
func DeletePermissionsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["permission_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var permission types.Permission
	if err := permission.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(permission); err != nil {
		panic(err)
	}
}
