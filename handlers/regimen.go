package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/regimen.go

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
// GetAccountRegimens is the GET method for an account's regimens
//
func GetAccountRegimens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

	regimens := types.Regimens{}

	if err := regimens.GetAccountRegimens(db, uuid.FromStringOrNil(accountID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	for i, regimen := range regimens {
		regimens[i] = scoping.FilterByScopes(scopedFields.([]string), regimen).(types.Regimen)
	}

	util.PaginationResponder(w, r, regimens)
}

//
// GetUserRegimens is the GET method for an account's regimens
//
func GetUserRegimens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := context.GetOk(r, "user_id")
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

	user := types.User{
		ID: uuid.FromStringOrNil(userID.(string)),
	}
	if err := user.GetRegimens(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	for i, regimen := range user.Regimens {
		user.Regimens[i] = scoping.FilterByScopes(scopedFields.([]string), regimen).(types.Regimen)
	}
	util.PaginationResponder(w, r, user.Regimens)
}

//
// GetAccountRegimensByID is the GET method for an account's regimens by ID
//
func GetAccountRegimensByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	regimenID := mux.Vars(r)["regimen_id"]
	if regimenID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
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
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var regimen types.Regimen
	if err := regimen.GetAccountRegimenByID(db, uuid.FromStringOrNil(regimenID), uuid.FromStringOrNil(accountID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), regimen))
}

//
// GetUserRegimensByID is the GET method for an account's regimens by ID
//
func GetUserRegimensByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	regimenID := mux.Vars(r)["regimen_id"]
	if regimenID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	userID, ok := context.GetOk(r, "user_id")
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

	var regimen types.Regimen
	if err := regimen.GetUserRegimenByID(db, uuid.FromStringOrNil(regimenID), uuid.FromStringOrNil(userID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), regimen))
}

//
// PutAccountRegimensByID is the PUT method for an account's regimens by ID
//
func PutAccountRegimensByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	regimenID := mux.Vars(r)["regimen_id"]
	if regimenID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	scopedFields, ok := context.GetOk(r, "scoped_fields")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No scopes"))
		return
	}
	var regimen types.Regimen
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&regimen)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if regimen.ID.String() != regimenID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update regimen ID"))
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

	currentRegimen := types.Regimen{}
	if err := currentRegimen.GetByID(db, uuid.FromStringOrNil(regimenID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	regimenExists := true
	existingRegimen := types.Regimen{}
	if currentRegimen.PodID == nil && regimen.PodID != nil {
		if err := existingRegimen.GetOneByQuery(db, "account_id = ? and pod_id = ?", accountID, regimen.PodID); err != nil {
			if err.Error() != "record not found" {
				util.ErrorResponder(w, http.StatusInternalServerError, err)
				return
			}
			regimenExists = false
		}
	}

	if !regimenExists || !(currentRegimen.PodID == nil && regimen.PodID != nil) {
		account := types.Account{
			ID: uuid.FromStringOrNil(accountID.(string)),
		}
		if err := account.UpdateAccountRegimenByID(db, regimen); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
		util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), regimen))
	}

	if err := db.Model(&existingRegimen).Association("Usages").Append(currentRegimen.Usages).Error; err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	existingRegimen.LastReportedServingsRemaining = currentRegimen.LastReportedServingsRemaining
	if err := existingRegimen.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	deleteRegimen := types.Regimen{}
	if err := deleteRegimen.Delete(db, uuid.FromStringOrNil(regimenID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), existingRegimen))
}

//
// DeleteAccountRegimenByID is the DELETE method for an account's regimens by ID
//
func DeleteAccountRegimenByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	regimenID := mux.Vars(r)["regimen_id"]
	if regimenID == "" {
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

	var regimen types.Regimen
	if err := regimen.DeleteAccountRegimenByID(db, uuid.FromStringOrNil(regimenID), uuid.FromStringOrNil(accountID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, map[string]string{"status": "success"})
}

//
// GetRegimen is the GET method for a regimen
//
func GetRegimen(w http.ResponseWriter, r *http.Request) {
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = util.SetDBPagination(db, r)

	var regimens types.Regimens
	if err := regimens.Get(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, regimens)
}

//
// GetRegimenByID is the GET method for a regimen by ID
//
func GetRegimenByID(w http.ResponseWriter, r *http.Request) {

	regimenID := mux.Vars(r)["regimen_id"]
	if regimenID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var regimen types.Regimen
	if err := regimen.GetByID(db, uuid.FromStringOrNil(regimenID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(regimen); err != nil {
		panic(err)
	}
}

//
// PutRegimenByID is the PUT method for a regimen by ID
//
func PutRegimenByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	regimenID := mux.Vars(r)["regimen_id"]
	if regimenID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	var regimen types.Regimen
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&regimen)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if regimen.ID.String() != regimenID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update regimen ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := regimen.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(regimen); err != nil {
		panic(err)
	}
}

//
// DeleteRegimenByID is the DELETE method for a regimen by ID
//
func DeleteRegimenByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["regimen_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var regimen types.Regimen
	if err := regimen.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(regimen); err != nil {
		panic(err)
	}

}
