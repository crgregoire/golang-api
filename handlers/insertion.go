package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/insertion.go

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
// GetInsertions is the GET method for a insertions' insertions
//
func GetInsertions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = util.SetDBPagination(db, r)

	var insertions types.Insertions
	if err := insertions.Get(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, insertions)
}

//
// GetInsertionByID is the GET method for insertions by ID
//
func GetInsertionByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["insertion_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var insertion types.Insertion
	if err := insertion.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(insertion); err != nil {
		panic(err)
	}
}

//
// PostInsertion is the POST method for a insertions' insertions
//
func PostInsertion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var insertion types.Insertion
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&insertion)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := insertion.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(insertion); err != nil {
		panic(err)
	}
}

//
// PutInsertionByID is the PUT method for insertions
//
func PutInsertionByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	insertionID := mux.Vars(r)["insertion_id"]
	if insertionID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	var insertion types.Insertion
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&insertion)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if insertion.ID.String() != insertionID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update insertion ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := insertion.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(insertion); err != nil {
		panic(err)
	}
}

//
// DeleteInsertionsByID is the DELETE method for a insertions' insertions by ID
//
func DeleteInsertionsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["insertion_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var insertion types.Insertion
	if err := insertion.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(insertion); err != nil {
		panic(err)
	}
}
