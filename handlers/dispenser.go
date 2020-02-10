package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/dispenser.go

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/tespo/buddha/db"
	"github.com/tespo/buddha/util"
	"github.com/tespo/satya/v2/scoping"
	"github.com/tespo/satya/v2/types"
)

//
// GetAccountDispensers is the GET method for an account's dispensers
//
func GetAccountDispensers(w http.ResponseWriter, r *http.Request) {
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
	db = util.SetDBPagination(db, r)

	connections := types.Connections{}

	dispensers, err := connections.GetAccountDispensers(db, uuid.FromStringOrNil(accountID.(string)))
	if err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, dispensers)

}

//
// GetDispenserByID is the GET method for a dispenser by ID
//
func GetDispenserByID(w http.ResponseWriter, r *http.Request) {
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
	dispenserID := mux.Vars(r)["dispenser_id"]
	if dispenserID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	connection := types.Connection{}
	if err := connection.GetAccountDispenserByID(db, uuid.FromStringOrNil(accountID.(string)), uuid.FromStringOrNil(dispenserID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), connection.Dispensers[0]))

}

//
// PutDispenserByID is the PUT method for a dispenser by ID
//
func PutDispenserByID(w http.ResponseWriter, r *http.Request) {
	dispenserID := mux.Vars(r)["dispenser_id"]
	if dispenserID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No Dispenser ID supplied"))
		return
	}

	var dispenser types.Dispenser
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dispenser)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if dispenser.ID.String() != dispenserID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update dispenser ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := dispenser.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.NewEncoder(w).Encode(dispenser); err != nil {
		panic(err)
	}
}

//
// DeleteDispenser is the DELETE method for a dispenser
//
func DeleteDispenser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	connections := types.Connections{}
	dispensers, err := connections.GetAccountDispensers(db, uuid.FromStringOrNil(accountID.(string)))
	if err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	var dispenser types.Dispenser
	dispenser.GetByID(db, dispensers[0].ID)
	if err := dispenser.Delete(db, dispenser.ID); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, map[string]string{"status": "success"})
}

// Dev routes

//
// GetDispensers is the GET method for a dispensers' dispensers
//
func GetDispensers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var dispensers types.Dispensers
	if err := dispensers.Get(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, dispensers)
}

//
// GetDispensersByID is the GET method for a dispensers' dispensers by ID
//
func GetDispensersByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["dispenser_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var dispenser types.Dispenser
	if err := dispenser.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(dispenser); err != nil {
		panic(err)
	}
}

//
// PostDispensers is the POST method for a dispensers' dispensers
//
func PostDispensers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dispenser types.Dispenser
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dispenser)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := dispenser.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(dispenser); err != nil {
		panic(err)
	}
}

//
// PutDispensersByID is the PUT method for a dispensers' dispensers by ID
//
func PutDispensersByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dispenserID := mux.Vars(r)["dispenser_id"]
	if dispenserID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	var dispenser types.Dispenser
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dispenser)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if dispenser.ID.String() != dispenserID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update dispenser ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := dispenser.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(dispenser); err != nil {
		panic(err)
	}
}

//
// DeleteDispensersByID is the DELETE method for a dispensers' dispensers by ID
//
func DeleteDispensersByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["dispenser_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var dispenser types.Dispenser
	if err := dispenser.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(dispenser); err != nil {
		panic(err)
	}
}
