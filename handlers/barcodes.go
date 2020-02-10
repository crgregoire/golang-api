package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/barcode.go

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
// GetBarcodes is the GET method for barcodes
//
func GetBarcodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = util.SetDBPagination(db, r)

	var barcodes types.Barcodes
	if err := barcodes.Get(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, barcodes)
}

//
// GetBarcodesByID is the GET method for a barcode by ID
//
func GetBarcodesByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["barcode_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var barcode types.Barcode
	if err := barcode.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(barcode); err != nil {
		panic(err)
	}
}

//
// PostBarcodes is the POST method for a barcode
//
func PostBarcodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var barcode types.Barcode
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&barcode)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := barcode.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(barcode); err != nil {
		panic(err)
	}
}

//
// PutBarcodesByID is the PUT method for a barcodes' barcodes by ID
//
func PutBarcodesByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	barcodeID := mux.Vars(r)["barcode_id"]
	if barcodeID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	var barcode types.Barcode
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&barcode)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if barcode.ID.String() != barcodeID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update barcode ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := barcode.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(barcode); err != nil {
		panic(err)
	}
}

//
// PutBarcodesByCode is the PUT method for a barcode by
// the code on the record
//
func PutBarcodesByCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	code := mux.Vars(r)["code"]
	if code == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	var barcode types.Barcode
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&barcode)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	existingBarcode := types.Barcode{}
	if err := existingBarcode.GetOneByQuery(db, "code = ?", code); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	barcode.ID = existingBarcode.ID

	if err := barcode.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(barcode); err != nil {
		panic(err)
	}
}

//
// DeleteBarcodesByID is the DELETE method for a barcodes' barcodes by ID
//
func DeleteBarcodesByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["barcode_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var barcode types.Barcode
	if err := barcode.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(barcode); err != nil {
		panic(err)
	}
}
