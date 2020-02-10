package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/account.go

import (
	"encoding/json"
	"errors"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/tespo/buddha/db"
	"github.com/tespo/buddha/util"
	"github.com/tespo/satya/v2/scoping"
	"github.com/tespo/satya/v2/types"
)

// SELF HANDLERS

//
// GetAccount is the Get method for Account
//
func GetAccount(w http.ResponseWriter, r *http.Request) {
	accountID, ok := context.GetOk(r, "account_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}
	scopedFields, ok := context.GetOk(r, "scoped_fields")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No scopes"))
		return
	}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	account := types.Account{}
	if err := account.GetByID(db, uuid.FromStringOrNil(accountID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), account))
}

//
// PutAccount is the Put method for Account
//
func PutAccount(w http.ResponseWriter, r *http.Request) {
	accountID, ok := context.GetOk(r, "account_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}
	scopedFields, ok := context.GetOk(r, "scoped_fields")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No scopes"))
		return
	}
	var account types.Account
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if account.ID.String() != accountID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update account ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	account.ID = uuid.FromStringOrNil(accountID.(string))
	if err := account.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), account))
}

// DEVELOPER HANDLERS

//
// GetAccountByID is the Get method for Account By the ID for Developers
//
func GetAccountByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["account_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var account types.Account
	if err := account.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(account); err != nil {
		panic(err)
	}
}

//
// GetAccounts is the Get method for all Accounts for Developers
//
func GetAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db = util.SetDBPagination(db, r)

	var accounts types.Accounts
	if err := accounts.Get(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, accounts)
}

//
// PostAccount is the POST method for Accounts for Developers
//
func PostAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var account types.Account
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := account.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(account); err != nil {
		panic(err)
	}
}

//
// PutAccountByID is the PUT method for Accounts by IDs for Developers
//
func PutAccountByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	accountID := mux.Vars(r)["account_id"]
	if accountID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	var account types.Account
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if account.ID.String() != accountID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update account ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := account.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(account); err != nil {
		panic(err)
	}
}

//
// DeleteAccountByID is the DELETE method for Accounts by IDs for Developers
//
func DeleteAccountByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["account_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var account types.Account
	if err := account.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(account); err != nil {
		panic(err)
	}
}
