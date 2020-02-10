package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/user.go

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
// GetSelfUser is the GET method for a users' users
//
func GetSelfUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := context.GetOk(r, "user_id")
	if !ok {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
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
	if err := user.GetByID(db, user.ID); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), user))
}

//
// PutSelfUser is the PUT method for a users' users
//
func PutSelfUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := context.GetOk(r, "user_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
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
	var user types.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	user.ID = uuid.FromStringOrNil(userID.(string))
	user.AccountID = uuid.FromStringOrNil(accountID.(string))
	if err := user.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), user))
}

//
// GetUsers is the GET method for a users' users
//
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = util.SetDBPagination(db, r)

	var users types.Users
	if err := users.Get(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, users)
}

//
// GetUserByID is the GET method for a users' users by ID
//
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["user_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user types.User
	if err := user.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

//
// PostUsers is the POST method for a users' users
//
func PostUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user types.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := user.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

//
// PutUsersByID is the PUT method for users by ID
//
func PutUsersByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No User ID supplied"))
		return
	}
	var user types.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if user.ID.String() != userID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update user ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := user.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

//
// DeleteUsersByID is the DELETE method for a users' users by ID
//
func DeleteUsersByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["user_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user types.User
	if err := user.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

//
// GetUsersByAccountID is the GET method for a user by account for developers
//
func GetUsersByAccountID(w http.ResponseWriter, r *http.Request) {
	accountID := mux.Vars(r)["account_id"]
	if accountID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No Account ID supplied"))
		return
	}

	var account types.Account

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := account.GetByID(db, uuid.FromStringOrNil(accountID)); err != nil {
		util.ErrorResponder(w, http.StatusNoContent, err)
		return
	}
	if err := account.GetUsers(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, account.Users)
}

//
// GetAccountUsers is the GET method for a user by account for developers
//
func GetAccountUsers(w http.ResponseWriter, r *http.Request) {
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
	account := types.Account{
		ID: uuid.FromStringOrNil(accountID.(string)),
	}
	if err := account.GetUsers(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	for i, user := range account.Users {
		account.Users[i] = scoping.FilterByScopes(scopedFields.([]string), user).(types.User)
	}

	util.PaginationResponder(w, r, account.Users)
}

//
// GetAccountUsersByUserID is the GET method for a user by account for developers
//
func GetAccountUsersByUserID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
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
	account := types.Account{
		ID: uuid.FromStringOrNil(accountID.(string)),
	}
	if err := account.GetUserByID(db, uuid.FromStringOrNil(userID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), account.Users[0]))
}

//
// CreateAccountUser is the POST method for users by account for developers
//
func CreateAccountUser(w http.ResponseWriter, r *http.Request) {
	var user types.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
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

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	account := types.Account{
		ID: uuid.FromStringOrNil(accountID.(string)),
	}
	if err := account.UpdateAccountUserByID(db, user); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), account.Users[0]))
}

//
// PutAccountUsersByUserID is the POST method for users by account for developers
//
func PutAccountUsersByUserID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	var user types.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
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
	if user.ID.String() != userID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update user ID"))
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
	if err := account.UpdateAccountUserByID(db, user); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), account.Users[0]))
}

//
// DeleteAccountUserByID is the DELETE method for users by accoun
//
func DeleteAccountUserByID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	accountID, ok := context.GetOk(r, "account_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}
	requestedUserID, ok := context.GetOk(r, "user_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}

	if userID == requestedUserID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("cannot delete yourself"))
		return
	}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user types.User
	if err := user.GetByID(db, uuid.FromStringOrNil(requestedUserID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	if !user.Owner {
		util.ErrorResponder(w, http.StatusUnauthorized, errors.New("owner operation only"))
		return
	}
	deleteUser := types.User{}
	if err := deleteUser.GetByQuery(db, "id = ? AND account_id = ?", uuid.FromStringOrNil(userID), uuid.FromStringOrNil(accountID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusNotFound, err)
		return
	}

	userRegimens := types.Regimens{}
	if err := userRegimens.GetByQuery(db, "user_id = ?", uuid.FromStringOrNil(userID)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	for _, regimen := range userRegimens {
		userUsages := types.Usages{}
		if err := userUsages.GetByQuery(db, "regimen_id = ?", uuid.FromStringOrNil(userID)); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
		for _, usage := range userUsages {
			usage.UserID = nil
			if err := usage.Update(db); err != nil {
				util.ErrorResponder(w, http.StatusInternalServerError, err)
				return
			}
		}
		regimen.UserID = nil
		regimen.User = types.User{}
		if err := regimen.Update(db); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
	}

	if err := deleteUser.Delete(db, deleteUser.ID); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	util.JSONResponder(w, map[string]string{"Status": "Success"})
}

//
// PostUsersByAccountID is the PUT method for users by account for developers
//
func PostUsersByAccountID(w http.ResponseWriter, r *http.Request) {
	accountID := mux.Vars(r)["account_id"]
	if accountID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	var user types.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	user.AccountID = uuid.FromStringOrNil(accountID)
	if err := user.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

//
// PutUsersByAccountIDAndUserID is the PUT method for users by account for developers
//
func PutUsersByAccountIDAndUserID(w http.ResponseWriter, r *http.Request) {

	accountID := mux.Vars(r)["account_id"]
	if accountID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No Account ID supplied"))
		return
	}
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No User ID supplied"))
		return
	}

	var user types.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if user.ID.String() != userID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update user ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	user.ID = uuid.FromStringOrNil(userID)
	user.AccountID = uuid.FromStringOrNil(accountID)
	if err := user.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

//
// PutUsersByExternalID is the PUT method for users by their external (wordpress) ID
//
func PutUsersByExternalID(w http.ResponseWriter, r *http.Request) {

	externalID := mux.Vars(r)["external_id"]
	if externalID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No external ID supplied"))
		return
	}

	var user, userUpdates types.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userUpdates)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if queryErr := user.GetByQuery(db, "external_id = ?", externalID); queryErr != nil {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Error getting user by external ID"))
		return
	}
	if userUpdates.FirstName != "" {
		user.FirstName = userUpdates.FirstName
	}
	if userUpdates.LastName != "" {
		user.LastName = userUpdates.LastName
	}
	if userUpdates.Email != "" {
		user.Email = userUpdates.Email
	}

	if err := user.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

//
// DeleteUsersByAccountIDAndUserID is the DELETE method for users by account for developers
//
func DeleteUsersByAccountIDAndUserID(w http.ResponseWriter, r *http.Request) {
	accountID := mux.Vars(r)["account_id"]
	if accountID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No Account ID supplied"))
		return
	}
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No User ID supplied"))
		return
	}

	var user types.User
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	user.ID = uuid.FromStringOrNil(userID)
	user.AccountID = uuid.FromStringOrNil(accountID)
	if err := user.Delete(db, user.ID); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}
