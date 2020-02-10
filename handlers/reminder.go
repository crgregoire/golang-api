package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/reminder.go

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
// GetUserRemindersByRegimenID is the GET method for reminders by regimen
//
func GetUserRemindersByRegimenID(w http.ResponseWriter, r *http.Request) {
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
	userUUID := uuid.FromStringOrNil(userID.(string))
	regimen := types.Regimen{
		ID:     uuid.FromStringOrNil(regimenID),
		UserID: &userUUID,
	}
	if err := regimen.GetReminders(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	for i, reminder := range regimen.Reminders {
		regimen.Reminders[i] = scoping.FilterByScopes(scopedFields.([]string), reminder).(types.Reminder)
	}
	util.PaginationResponder(w, r, regimen.Reminders)
}

//
// GetReminders is the GET method for all user reminders
//
func GetReminders(w http.ResponseWriter, r *http.Request) {
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
	userUUID := uuid.FromStringOrNil(userID.(string))

	reminders := types.Reminders{}
	if err := reminders.GetByQuery(db, "user_id = ?", userUUID); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	for i, reminder := range reminders {
		reminders[i] = scoping.FilterByScopes(scopedFields.([]string), reminder).(types.Reminder)
	}

	util.PaginationResponder(w, r, reminders)

}

//
// PostReminder is the POST method for reminders by regimen
//
func PostReminder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	regimenID := mux.Vars(r)["regimen_id"]
	if regimenID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	var reminder types.Reminder
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reminder)
	if err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	accountID, ok := context.GetOk(r, "account_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}

	userID, ok := context.GetOk(r, "user_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userUUID := uuid.FromStringOrNil(userID.(string))
	regimen := types.Regimen{
		ID:        uuid.FromStringOrNil(regimenID),
		UserID:    &userUUID,
		AccountID: uuid.FromStringOrNil(accountID.(string)),
	}

	if err := regimen.CreateReminder(db, reminder); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(regimen.Reminders); err != nil {
		panic(err)
	}
}

//
// PutReminderByID is the PUT method for reminders by regimen id
//
func PutReminderByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	regimenID := mux.Vars(r)["regimen_id"]
	reminderID := mux.Vars(r)["reminder_id"]
	if regimenID == "" || reminderID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	var reminder types.Reminder
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reminder)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if reminder.ID.String() != reminderID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update reminder ID"))
		return
	}

	accountID, ok := context.GetOk(r, "account_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}

	userID, ok := context.GetOk(r, "user_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userUUID := uuid.FromStringOrNil(userID.(string))
	regimen := types.Regimen{
		ID:        uuid.FromStringOrNil(regimenID),
		UserID:    &userUUID,
		AccountID: uuid.FromStringOrNil(accountID.(string)),
	}

	if err := regimen.UpdateReminder(db, reminder); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(regimen.Reminders); err != nil {
		panic(err)
	}
}

//
// DeleteReminderByID is the DELETE method for reminders by regimen id
//
func DeleteReminderByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	regimenID := mux.Vars(r)["regimen_id"]
	reminderID := mux.Vars(r)["reminder_id"]
	if regimenID == "" || reminderID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	accountID, ok := context.GetOk(r, "account_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}

	userID, ok := context.GetOk(r, "user_id")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userUUID := uuid.FromStringOrNil(userID.(string))

	regimen := types.Regimen{
		ID:        uuid.FromStringOrNil(regimenID),
		UserID:    &userUUID,
		AccountID: uuid.FromStringOrNil(accountID.(string)),
	}
	now := time.Now()
	reminder := types.Reminder{
		ID:        uuid.FromStringOrNil(reminderID),
		RegimenID: uuid.FromStringOrNil(regimenID),
		UserID:    uuid.FromStringOrNil(userID.(string)),
		DeletedAt: &now,
	}

	if err := regimen.UpdateReminder(db, reminder); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(regimen.Reminders); err != nil {
		panic(err)
	}
}
