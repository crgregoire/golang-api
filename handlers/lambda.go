package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/tespo/buddha/db"
	"github.com/tespo/buddha/util"
	"github.com/tespo/satya/v2/types"
)

//
// DispenserDispensed handles the lambda message
// for Dispenser Dispensed
//
func DispenserDispensed(w http.ResponseWriter, r *http.Request) {
	lambdaMessage := types.LambdaMessage{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lambdaMessage); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	barcode := types.Barcode{
		Code: lambdaMessage.Payload.Pod.Barcode,
	}
	pod := types.Pod{}
	regimen := types.Regimen{}
	insertion := types.Insertion{}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	connections := types.Connections{}
	dispensers, err := connections.GetAccountDispensers(db, uuid.FromStringOrNil(lambdaMessage.Payload.Customer.ID))
	if err != nil {
		util.ErrorResponder(w, http.StatusNotFound, err)
		return
	}
	dispenser := types.Dispenser{}
	for _, disp := range dispensers {
		if disp.Serial == lambdaMessage.Payload.Dispenser.Serial {
			dispenser = disp
		}
	}
	if dispenser.Serial == "" {
		util.ErrorResponder(w, http.StatusNotFound, errors.New("dispenser cannot be found"))
		return
	}
	barcodeFound := false
	if barcode.Code != "" {
		if err := barcode.GetOneByQuery(db, "code = ?", barcode.Code); err != nil {
			if err.Error() != "record not found" {
				util.ErrorResponder(w, http.StatusNotFound, err)
				return
			}
		} else {
			barcodeFound = true
		}
	}
	if barcodeFound {
		if err := pod.GetByID(db, barcode.PodID); err != nil {
			util.ErrorResponder(w, http.StatusNotFound, err)
			return
		}
		if err := insertion.GetByQuery(db, "dispenser_id = ? AND barcode_id = ?", dispenser.ID, barcode.ID); err != nil {
			util.ErrorResponder(w, http.StatusNotFound, err)
			return
		}
		if err := regimen.GetByID(db, insertion.RegimenID); err != nil {
			util.ErrorResponder(w, http.StatusNotFound, err)
			return
		}
	}
	if regimen.ID.String() == uuid.Nil.String() {
		//GORM apparently inverts ordering... thus asc == desc and desc == asc
		if err := insertion.GetByQuery(db.Order("created_at asc"), "dispenser_id = ?", dispenser.ID); err != nil {
			util.ErrorResponder(w, http.StatusNotFound, err)
			return
		}
		if err := regimen.GetByID(db, insertion.RegimenID); err != nil {
			util.ErrorResponder(w, http.StatusNotFound, err)
			return
		}
	}

	newUsage := types.Usage{
		RegimenID:   regimen.ID,
		DispenserID: dispenser.ID,
		UserID:      regimen.UserID,
		BarcodeID:   &barcode.ID,
		Servings:    uint(lambdaMessage.Payload.Pod.ServingsRemaining),
		Flags:       uint(lambdaMessage.Payload.Pod.Flags),
	}

	if err := newUsage.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	regimen.LastReportedServingsRemaining = uint(lambdaMessage.Payload.Pod.ServingsRemaining)
	if err := regimen.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	meta := map[string]map[string]string{
		"pcb":        map[string]string{"version": lambdaMessage.Payload.Dispenser.ControllerFirmwareVersion},
		"wifi":       map[string]string{"version": lambdaMessage.Payload.Dispenser.WifiFirmwareVersion},
		"controller": map[string]string{"version": lambdaMessage.Payload.Dispenser.PcbFirmwareVersion},
	}
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	dispenser.Meta = metaBytes
	dispenser.Update(db)

	util.JSONResponder(w, map[string]string{"status": "success"})
}

//
// PodInserted handles the lambda message
// for Pod Inserted
//
func PodInserted(w http.ResponseWriter, r *http.Request) {
	lambdaMessage := types.LambdaMessage{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lambdaMessage); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	pod := types.Pod{}
	insertion := types.Insertion{}
	account := types.Account{}
	barcode := types.Barcode{}
	regimen := types.Regimen{}
	connections := types.Connections{}
	dispensers, err := connections.GetAccountDispensers(db, uuid.FromStringOrNil(lambdaMessage.Payload.Customer.ID))
	if err != nil {
		util.ErrorResponder(w, http.StatusNotFound, err)
		return
	}
	var serial = lambdaMessage.Payload.Dispenser.Serial
	dispenser := types.Dispenser{}
	for _, disp := range dispensers {
		if disp.Serial == serial {
			dispenser = disp
			break
		}
	}
	if dispenser.Serial == "" {
		util.ErrorResponder(w, http.StatusNotFound, fmt.Errorf("Could not find dispenser with serial %v", serial))
		return
	}
	newRegimen := false
	var userID uuid.UUID
	if err := account.GetByID(db, uuid.FromStringOrNil(lambdaMessage.Payload.Customer.ID)); err != nil {
		util.ErrorResponder(w, http.StatusNotFound, err)
		return
	}
	if err := account.GetUsers(db); err != nil {
		util.ErrorResponder(w, http.StatusNotFound, err)
		return
	}

	if len(account.Users) == 1 {
		userID = account.Users[0].ID
	}

	if lambdaMessage.Payload.Pod.Barcode != "" {
		if err := barcode.GetOneByQuery(db, "code = ?", lambdaMessage.Payload.Pod.Barcode); err != nil {
			if err.Error() != "record not found" {
				util.ErrorResponder(w, http.StatusNotFound, err)
				return
			}
			newRegimen = true
		}
		if err := pod.GetByID(db, barcode.PodID); err != nil {
			if err.Error() != "record not found" {
				util.ErrorResponder(w, http.StatusNotFound, err)
				return
			}
		}
		err := regimen.GetOneByQuery(db, "pod_id = ? AND account_id = ?", pod.ID, account.ID)
		if err != nil {
			if err.Error() != "record not found" {
				util.ErrorResponder(w, http.StatusNotFound, err)
				return
			}
			newRegimen = true
		}
	}
	if lambdaMessage.Payload.Pod.Barcode == "" || newRegimen {
		regimen = types.Regimen{
			PodID:                         &pod.ID,
			AccountID:                     account.ID,
			UserID:                        &userID,
			LastReportedServingsRemaining: uint(lambdaMessage.Payload.Pod.ServingsRemaining),
		}
		if err := regimen.Create(db); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		regimen.LastReportedServingsRemaining = uint(lambdaMessage.Payload.Pod.ServingsRemaining)
		if err := regimen.Update(db); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
	}

	insertion = types.Insertion{
		RegimenID:   regimen.ID,
		DispenserID: dispenser.ID,
		BarcodeID:   &barcode.ID,
		Flags:       uint(lambdaMessage.Payload.Pod.Flags),
		LabelTall:   barcode.LabelTall,
		LabelWide:   barcode.LabelWide,
	}
	if err := insertion.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	meta := map[string]map[string]string{
		"pcb":        map[string]string{"version": lambdaMessage.Payload.Dispenser.ControllerFirmwareVersion},
		"wifi":       map[string]string{"version": lambdaMessage.Payload.Dispenser.WifiFirmwareVersion},
		"controller": map[string]string{"version": lambdaMessage.Payload.Dispenser.PcbFirmwareVersion},
	}
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	dispenser.Meta = metaBytes
	dispenser.Update(db)

	util.JSONResponder(w, map[string]string{"status": "success"})
}

//
// DispenserConnected handles when the lambda function
// for dispensers connected
//
func DispenserConnected(w http.ResponseWriter, r *http.Request) {
	lambdaMessage := types.LambdaMessage{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lambdaMessage); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	account := types.Account{
		ID: uuid.FromStringOrNil(lambdaMessage.Payload.Customer.ID),
	}
	if err := account.GetConnections(db); err != nil {
		if err.Error() != "record not found" {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
	}
	if len(account.Connections) > 0 {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("account already has connection"))
		return
	}
	dispenser := types.Dispenser{}
	if err := dispenser.GetOneByQuery(db, "serial = ?", lambdaMessage.Payload.Dispenser.Serial); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	meta := map[string]map[string]string{
		"pcb":        map[string]string{"version": lambdaMessage.Payload.Dispenser.ControllerFirmwareVersion},
		"wifi":       map[string]string{"version": lambdaMessage.Payload.Dispenser.WifiFirmwareVersion},
		"controller": map[string]string{"version": lambdaMessage.Payload.Dispenser.PcbFirmwareVersion},
	}
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	dispenser.Meta = metaBytes
	dispenser.Update(db)
	connection := types.Connection{}
	if err := connection.GetOneByQuery(db, "dispenser_id = ?", dispenser.ID); err != nil {
		if err.Error() != "record not found" {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
		connection.DispenserID = dispenser.ID
		connection.AccountID = uuid.FromStringOrNil(lambdaMessage.Payload.Customer.ID)
		connection.ConnectedAt = time.Now()
		if err := connection.Create(db); err != nil {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
		util.JSONResponder(w, map[string]string{"status": "Success"})
		return
	}
	util.ErrorResponder(w, http.StatusBadRequest, errors.New("Connection already exists"))
	return
}

//
// DispenserDisconnected handles when the lambda function
// for dispensers being disconnected
//
func DispenserDisconnected(w http.ResponseWriter, r *http.Request) {
	lambdaMessage := types.LambdaMessage{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lambdaMessage); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	account := types.Account{
		ID: uuid.FromStringOrNil(lambdaMessage.Payload.Customer.ID),
	}
	if err := account.GetConnections(db); err != nil {
		if err.Error() != "record not found" {
			util.ErrorResponder(w, http.StatusInternalServerError, err)
			return
		}
	}
	dispenser := types.Dispenser{}
	if err := dispenser.GetOneByQuery(db, "serial = ?", lambdaMessage.Payload.Dispenser.Serial); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	connection := types.Connection{}
	if err := connection.GetOneByQuery(db, "dispenser_id = ?", dispenser.ID); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	now := time.Now()
	connection.DisconnectedAt = &now
	connection.Update(db)
	if err := connection.Delete(db, connection.ID); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}
	util.JSONResponder(w, map[string]string{"status": "success"})
	return
}