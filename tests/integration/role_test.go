package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/tespo/satya/v2/types"
)

var role = types.Role{
	ID:   uuid.NewV4(),
	Name: "Integration Test Role",
}

func TestCreateRole(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	data, err := json.Marshal(role)
	if err != nil {
		tests.Error(err)
		return
	}
	response, err := http.Post(os.Getenv("TESTING_URL")+"/role", "application/json", bytes.NewBuffer(data))
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testRole := types.Role{}
	json.Unmarshal(body, &testRole)
	if err := json.Unmarshal(body, &testRole); err != nil {
		tests.Error(err)
		return
	}
	if testRole.Name != role.Name {
		tests.Error(errors.New("Returned the wrong role"))
		return
	}
	role.ID = testRole.ID
}

func TestGetRole(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/roles")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	paginatedResponse := types.PaginatedResponse{}
	if err := json.Unmarshal(body, &paginatedResponse); err != nil {
		tests.Error(err)
		return
	}
	roleData, err := json.Marshal(paginatedResponse.Data)
	if err != nil {
		tests.Error(err)
		return
	}
	testRoles := types.Roles{}
	if err := json.Unmarshal(roleData, &testRoles); err != nil {
		tests.Error(err)
		return
	}
}

func TestGetRoleByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/roles/" + role.ID.String())
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testRole := types.Role{}
	if err := json.Unmarshal(body, &testRole); err != nil {
		tests.Error(err)
		return
	}
	if testRole.ID != role.ID {
		tests.Error(errors.New("Returned the wrong role based on ID"))
		return
	}
}

func TestPutRolesByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	role.Meta = []byte("{\"Updated?\":\"YES\"}")
	data, err := json.Marshal(role)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/roles/"+role.ID.String(), bytes.NewBuffer(data))
	response, err := client.Do(request)
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testRole := types.Role{}
	if err := json.Unmarshal(body, &testRole); err != nil {
		tests.Error(err)
		return
	}
	if testRole.Meta == nil {
		tests.Error(errors.New("Wrong role, failed to update"))
		return
	}
	role.ID = testRole.ID
}

func TestDeleteRolesByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", os.Getenv("TESTING_URL")+"/roles/"+role.ID.String(), nil)
	response, err := client.Do(request)
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testRole := types.Role{}
	if err := json.Unmarshal(body, &testRole); err != nil {
		tests.Error(err)
		return
	}
}
