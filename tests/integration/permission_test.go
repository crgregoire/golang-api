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

var permission = types.Permission{
	ID:      uuid.NewV4(),
	Slug:    "/integration-test-permission",
	Actions: []byte("[\"GET\"]"),
	Meta:    []byte("{\"Name\":\"Permission Integration Test\"}"),
}

func TestCreatePermission(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	data, err := json.Marshal(permission)
	if err != nil {
		tests.Error(err)
		return
	}
	response, err := http.Post(os.Getenv("TESTING_URL")+"/permissions", "application/json", bytes.NewBuffer(data))
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testPermissions := types.Permission{}
	if err := json.Unmarshal(body, &testPermissions); err != nil {
		tests.Error(err)
		return
	}
	if testPermissions.ID != permission.ID {
		tests.Error(errors.New("Returned the wrong permission"))
		return
	}
	permission.ID = testPermissions.ID
}

func TestGetPermissionByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/permissions/" + permission.ID.String())
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testPermissions := types.Permission{}
	if err := json.Unmarshal(body, &testPermissions); err != nil {
		tests.Error(err)
		return
	}
	if testPermissions.ID != permission.ID {
		tests.Error(errors.New("Returned the wrong permission"))
		return
	}
}

func TestGetPermissions(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/permissions")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil && body != nil {
		tests.Error(err)
		return
	}
	testPermissions := types.PaginatedResponse{}.Data
	if err := json.Unmarshal(body, &testPermissions); err != nil {
		tests.Error(err)
		return
	}
	if testPermissions == nil {
		tests.Error(errors.New("Failed to return paginated permissions"))
		return
	}
}

func TestPutPermissionsByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	permission.Meta = []byte("{\"Updated?\":\"YES\"}")
	data, err := json.Marshal(permission)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/permissions/"+permission.ID.String(), bytes.NewBuffer(data))
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
	testPermissions := types.Permission{}
	if err := json.Unmarshal(body, &testPermissions); err != nil {
		tests.Error(err)
		return
	}
	if testPermissions.Meta == nil {
		tests.Error(errors.New("Wrong permission, failed to update"))
		return
	}
	permission.ID = testPermissions.ID
}

func TestDeletePermissionsByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", os.Getenv("TESTING_URL")+"/permissions/"+permission.ID.String(), nil)
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
	testPermissions := types.Permission{}
	if err := json.Unmarshal(body, &testPermissions); err != nil {
		tests.Error(err)
		return
	}
}
