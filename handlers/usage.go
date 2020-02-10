package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/usage.go

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tespo/satya/v2/scoping"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/tespo/buddha/db"
	"github.com/tespo/buddha/util"
	"github.com/tespo/satya/v2/types"
)

//
// GetUserUsages is the GET method for a users usages
//
func GetUserUsages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := context.GetOk(r, "user_id")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	scopedFields, ok := context.GetOk(r, "scoped_fields")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var usages types.Usages
	if err := usages.GetByQuery(db, "user_id = ?", uuid.FromStringOrNil(userID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	for i, usage := range usages {
		usages[i] = scoping.FilterByScopes(scopedFields.([]string), usage).(types.Usage)
	}

	util.PaginationResponder(w, r, usages)
}

//
// GetUserUsageByID is the GET method for a user usages
//
func GetUserUsageByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["usage_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	userID, ok := context.GetOk(r, "user_id")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	scopedFields, ok := context.GetOk(r, "scoped_fields")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var usages types.Usages
	if err := usages.GetByQuery(db, "id = ? AND user_id = ?", uuid.FromStringOrNil(id), uuid.FromStringOrNil(userID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	for i, usage := range usages {
		usages[i] = scoping.FilterByScopes(scopedFields.([]string), usage).(types.Usage)
	}

	if err = json.NewEncoder(w).Encode(usages); err != nil {
		panic(err)
	}
}

//
// PutUserUsageByID is the GET method for a user usages
//
func PutUserUsageByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["usage_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	userID, ok := context.GetOk(r, "user_id")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	var newUsageObject types.Usage
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUsageObject)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	scopedFields, ok := context.GetOk(r, "scoped_fields")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var usages types.Usages
	if err := usages.GetByQuery(db, "id = ? AND user_id = ?", uuid.FromStringOrNil(id), uuid.FromStringOrNil(userID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err := newUsageObject.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	for i, usage := range usages {
		usages[i] = scoping.FilterByScopes(scopedFields.([]string), usage).(types.Usage)
	}

	if err = json.NewEncoder(w).Encode(usages); err != nil {
		panic(err)
	}
}

//
// GetAccountUsages is the GET method for a account usages
//
func GetAccountUsages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := context.GetOk(r, "account_id")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	scopedFields, ok := context.GetOk(r, "scoped_fields")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var usages types.Usages
	if err := usages.GetByQuery(db, "account_id = ?", uuid.FromStringOrNil(userID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	for i, usage := range usages {
		usages[i] = scoping.FilterByScopes(scopedFields.([]string), usage).(types.Usage)
	}

	util.PaginationResponder(w, r, usages)
}

//
// GetAccountUsageByID is the GET method for a account usages
//
func GetAccountUsageByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["usage_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	scopedFields, ok := context.GetOk(r, "scoped_fields")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var usage types.Usage
	if err := usage.GetByQuery(db, "id = ?", uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), usage))
}

//
// PutAccountUsageByID is the GET method for a account usages
//
func PutAccountUsageByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["usage_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var newUsageObject types.Usage
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUsageObject)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	scopedFields, ok := context.GetOk(r, "scoped_fields")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var usage types.Usage
	if err := usage.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if newUsageObject.ID != usage.ID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("cannot update usage id"))
		return
	}

	if err := newUsageObject.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), newUsageObject))
}

//
// GetUsages is the GET method for a usages' usages
//
func GetUsages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = util.SetDBPagination(db, r)

	var usages types.Usages
	if err := usages.Get(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, usages)
}

//
// GetUsagesByID is the GET method for a usages' usages by ID
//
func GetUsagesByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["usage_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var usage types.Usage
	if err := usage.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(usage); err != nil {
		panic(err)
	}
}

//
// PostUsages is the POST method for a usages' usages
//
func PostUsages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var usage types.Usage
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&usage)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := usage.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(usage); err != nil {
		panic(err)
	}
}

//
// PutUsagesByID is the PUT method for a usages' usages by ID
//
func PutUsagesByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usageID := mux.Vars(r)["usage_id"]
	if usageID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No Usage ID supplied"))
		return
	}
	var usage types.Usage
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&usage)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if usage.ID.String() != usageID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update usage ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := usage.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(usage); err != nil {
		panic(err)
	}
}

//
// DeleteUsagesByID is the DELETE method for a usages' usages by ID
//
func DeleteUsagesByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["usage_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var usage types.Usage
	if err := usage.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(usage); err != nil {
		panic(err)
	}
}
