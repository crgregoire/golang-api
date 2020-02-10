package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
// GetInvitationByID will get an invitation by id
//
func GetInvitationByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["invitation_id"]
	if id == "" {
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
	invitation := types.Invitation{}
	if err := invitation.GetOneByQuery(db, "id = ? and account_id = ?", uuid.FromStringOrNil(id), accountID.(uuid.UUID)); err != nil {
		util.ErrorResponder(w, http.StatusNotFound, err)
		return
	}
	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), invitation))
}

//
// GetInvitations will get invitations
//
func GetInvitations(w http.ResponseWriter, r *http.Request) {
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
	invitations := types.Invitations{}
	if err := invitations.GetByQuery(db, "account_id = ?", uuid.FromStringOrNil(accountID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusNotFound, err)
		return
	}
	for i, invitation := range invitations {
		invitations[i] = scoping.FilterByScopes(scopedFields.([]string), invitation).(types.Invitation)
	}

	util.PaginationResponder(w, r, invitations)
}

//
// PostInvitation will create invitation
//
func PostInvitation(w http.ResponseWriter, r *http.Request) {
	var invitation types.Invitation
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&invitation)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
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
	invitation.AccountID = uuid.FromStringOrNil(accountID.(string))
	invitation.ExpiresAt = time.Now().Add(48 * time.Hour)
	invitation.Code = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	account := types.Account{}
	if err := account.GetByID(db, uuid.FromStringOrNil(accountID.(string))); err != nil {
		util.ErrorResponder(w, http.StatusNotFound, err)
		return
	}
	existingInvite := types.Invitation{}
	if err := existingInvite.GetOneByQuery(db, "account_id = ? and email = ?", invitation.AccountID, strings.ToLower(invitation.Email)); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
		if err := invitation.Create(db); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		invitation.ID = existingInvite.ID
		if err := invitation.Update(db); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
	}

	if err := util.SendInviteEmail(strings.ToLower(invitation.Email), account.Name, invitation.Code); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		panic(err)
	}

	util.JSONResponder(w, scoping.FilterByScopes(scopedFields.([]string), invitation))
}

//
// DeleteInvitation will delete invitation
//
func DeleteInvitation(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["invitation_id"]
	if id == "" {
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
	invitation := types.Invitation{
		ID:        uuid.FromStringOrNil(id),
		AccountID: uuid.FromStringOrNil(accountID.(string)),
	}
	if err := invitation.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	util.JSONResponder(w, map[string]string{"status": "success"})
}

//
// AcceptInvitation will accept invitation
//
func AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["invitation_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
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
	owner, ok := context.GetOk(r, "owner")
	if !ok {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("Cannot process token claims"))
		return
	}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	invitation := types.Invitation{}
	if err := invitation.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	acceptUser := types.User{
		ID: uuid.FromStringOrNil(userID.(string)),
	}
	if err := acceptUser.GetUserWithAllData(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if strings.ToLower(acceptUser.Email) != strings.ToLower(invitation.Email) || time.Since(invitation.ExpiresAt) > 0 {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("user cannot accept this invitation"))
		return
	}
	now := time.Now()
	if owner.(bool) {
		accountWithUsers := types.Account{
			ID: uuid.FromStringOrNil(accountID.(string)),
		}
		if err := accountWithUsers.GetUsers(db); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
		for _, forUser := range accountWithUsers.Users {
			forUser.AccountID = invitation.AccountID
			accountUser := &forUser
			if err := accountUser.Update(db); err != nil {
				util.ErrorResponder(w, http.StatusInternalServerError, err)
				return
			}
			if err := accountUser.GetUserWithAllData(db); err != nil {
				util.ErrorResponder(w, http.StatusInternalServerError, err)
				return
			}
			now := time.Now()
			// Delete old user data
			for _, regimen := range accountUser.Regimens {
				for _, reminder := range regimen.Reminders {
					reminder.Meta = []byte("{\"delete_cause\":\"user " + userID.(string) + " accepted invitation\"}")
					reminder.DeletedAt = &now
					if err := reminder.Update(db); err != nil {
						util.ErrorResponder(w, http.StatusInternalServerError, err)
						return
					}
				}
				for _, usage := range regimen.Usages {
					usage.Meta = []byte("{\"delete_cause\":\"user " + userID.(string) + " accepted invitation\"}")
					usage.DeletedAt = &now
					if err := usage.Update(db); err != nil {
						util.ErrorResponder(w, http.StatusInternalServerError, err)
						return
					}
				}
				regimen.Meta = []byte("{\"delete_cause\":\"user " + userID.(string) + " accepted invitation\"}")
				regimen.DeletedAt = &now
				if err := regimen.Update(db); err != nil {
					util.ErrorResponder(w, http.StatusInternalServerError, err)
					return
				}
			}
		}
		accountWithOutUsers := types.Account{}
		if err := accountWithOutUsers.GetByID(db, uuid.FromStringOrNil(accountID.(string))); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
		accountWithOutUsers.Meta = []byte("{\"delete_cause\":\"user " + userID.(string) + " accepted invitation\"}")
		accountWithOutUsers.DeletedAt = &now
		if err := accountWithOutUsers.Update(db); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
		util.JSONResponder(w, map[string]string{"status": "success"})
		return
	}

	acceptUser.AccountID = invitation.AccountID
	if err := acceptUser.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	// Delete old user data
	for _, regimen := range acceptUser.Regimens {
		for _, reminder := range regimen.Reminders {
			reminder.Meta = []byte("{\"delete_cause\":\"user " + userID.(string) + " accepted invitation\"}")
			reminder.DeletedAt = &now
			if err := reminder.Update(db); err != nil {
				util.ErrorResponder(w, http.StatusInternalServerError, err)
				return
			}
		}
		for _, usage := range regimen.Usages {
			usage.Meta = []byte("{\"delete_cause\":\"user " + userID.(string) + " accepted invitation\"}")
			usage.DeletedAt = &now
			if err := usage.Update(db); err != nil {
				util.ErrorResponder(w, http.StatusInternalServerError, err)
				return
			}
		}
		regimen.Meta = []byte("{\"delete_cause\":\"user " + userID.(string) + " accepted invitation\"}")
		regimen.DeletedAt = &now
		if err := regimen.Update(db); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
	}

	if err := invitation.Delete(db, invitation.ID); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, errors.New("invite accepted but could not delete invitation: "+err.Error()))
		return
	}

	util.JSONResponder(w, map[string]string{"status": "success"})
}
