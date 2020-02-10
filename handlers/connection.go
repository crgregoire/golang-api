package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/connection.go

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/tespo/buddha/db"
	"github.com/tespo/buddha/util"
	"github.com/tespo/satya/v2/scoping"
	"github.com/tespo/satya/v2/types"
)

//
// GetConnections is the GET method for a dispensers' connections
//
func GetConnections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db = util.SetDBPagination(db, r)

	var connections types.Connections
	if err := connections.Get(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, connections)
}

//
// GetConnectionsByID is the GET method for a dispensers' connections by ID
//
func GetConnectionsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["connection_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var connection types.Connection
	if err := connection.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(connection); err != nil {
		panic(err)
	}
}

//
// PostConnections is the POST method for a dispensers' connections
//
func PostConnections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var connection types.Connection
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&connection)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := connection.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(connection); err != nil {
		panic(err)
	}
}

//
// PutConnectionsByID is the PUT method for a dispensers' connections by ID
//
func PutConnectionsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	connectionID := mux.Vars(r)["connection_id"]
	if connectionID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	var connection types.Connection
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&connection)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if connection.ID.String() != connectionID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update connection ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := connection.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(connection); err != nil {
		panic(err)
	}
}

//
// DeleteConnectionsByID is the DELETE method for a dispensers' connections by ID
//
func DeleteConnectionsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["connection_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var connection types.Connection

	if err := connection.GetOneByQuery(db, "id = ?", id); err != nil {
		util.ErrorResponder(w, http.StatusNotFound, err)
		return
	}
	now := time.Now()
	connection.DisconnectedAt = &now
	connection.Update(db)
	if err := connection.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(connection); err != nil {
		panic(err)
	}
}

//
// GetAccountConnections is the GET method for an account's connections
//
func GetAccountConnections(w http.ResponseWriter, r *http.Request) {
	accountID, ok := context.GetOk(r, "account_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	account := types.Account{
		ID: uuid.FromStringOrNil(accountID.(string)),
	}
	if err := account.GetConnections(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, account.Connections)
}

//
// GetAccountConnectionByID is the GET method for a connection by ID
//
func GetAccountConnectionByID(w http.ResponseWriter, r *http.Request) {
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
	connectionID := mux.Vars(r)["connection_id"]
	if connectionID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	account := types.Account{
		ID: uuid.FromStringOrNil(accountID.(string)),
	}
	if err := account.GetConnectionByID(db, uuid.FromStringOrNil(connectionID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), account.Connections[0]))

}

//
// PutAccountConnectionsByID is the POST method for connections by account for developers
//
func PutAccountConnectionsByID(w http.ResponseWriter, r *http.Request) {
	connectionID := mux.Vars(r)["connection_id"]
	if connectionID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	var connection types.Connection
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&connection)
	if err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
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

	if connection.ID.String() != connectionID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update connection ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	account := types.Account{
		ID: uuid.FromStringOrNil(accountID.(string)),
	}
	if err := account.UpdateAccountConnectionByID(db, connection); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), account.Connections[0]))
}

//
// DeleteAccountConnectionsByID is the POST method for connections by account for developers
//
func DeleteAccountConnectionsByID(w http.ResponseWriter, r *http.Request) {
	connectionID := mux.Vars(r)["connection_id"]
	if connectionID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	accountID, ok := context.GetOk(r, "account_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	account := types.Account{
		ID: uuid.FromStringOrNil(accountID.(string)),
	}

	connection := types.Connection{}
	if err := connection.GetOneByQuery(db, "account_id = ?", account.ID); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	now := time.Now()
	connection.DeletedAt = &now
	connection.DisconnectedAt = &now
	if err := account.UpdateAccountConnectionByID(db, connection); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	util.JSONResponder(w, map[string]string{"status": "success"})
}
