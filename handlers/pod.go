package handlers

// Handlers for https://github.com/tespo/satya/v2/blob/develop/types/pod.go

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
// GetPods is the GET method for a pods' pods
//
func GetPods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = util.SetDBPagination(db, r)

	var pods types.Pods
	if err := pods.Get(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	util.PaginationResponder(w, r, pods)
}

//
// GetPodsByID is the GET method for a pods' pods by ID
//
func GetPodsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["pod_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var pod types.Pod
	if err := pod.GetByID(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(pod); err != nil {
		panic(err)
	}
}

//
// PostPods is the POST method for a pods' pods
//
func PostPods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var pod types.Pod
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pod)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := pod.Create(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(pod); err != nil {
		panic(err)
	}
}

//
// PutPodsByID is the PUT method for a pods' pods by ID
//
func PutPodsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	podID := mux.Vars(r)["pod_id"]
	if podID == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}
	var pod types.Pod
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pod)
	if err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}

	if pod.ID.String() != podID {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("Cannot update pod ID"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := pod.Update(db); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(pod); err != nil {
		panic(err)
	}
}

//
// DeletePodsByID is the DELETE method for a pods' pods by ID
//
func DeletePodsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["pod_id"]
	if id == "" {
		util.ErrorResponder(w, http.StatusBadRequest, errors.New("No ID supplied"))
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var pod types.Pod
	if err := pod.Delete(db, uuid.FromStringOrNil(id)); err != nil {
		util.ErrorResponder(w, http.StatusInternalServerError, err)
		return
	}

	if err = json.NewEncoder(w).Encode(pod); err != nil {
		panic(err)
	}
}
