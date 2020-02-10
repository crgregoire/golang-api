package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/tespo/satya/v2/types"
)

var connection = types.Connection{
	ID:          uuid.NewV4(),
	AccountID:   uuid.FromStringOrNil("d8e4c5dc-9767-41bd-b802-060e80d83867"),
	DispenserID: uuid.FromStringOrNil("9f717337-dba7-415c-9daf-c607df526d14"),
	Meta:        []byte("{\"Name\":\"Connection Integration Test\"}"),
	ConnectedAt: time.Now(),
}

func TestCreateConnection(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	data, err := json.Marshal(connection)
	if err != nil {
		tests.Error(err)
		return
	}
	response, err := http.Post(os.Getenv("TESTING_URL")+"/connections", "application/json", bytes.NewBuffer(data))
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testConnection := types.Connection{}
	if err := json.Unmarshal(body, &testConnection); err != nil {
		tests.Error(err)
		return
	}
	if testConnection.ID != connection.ID {
		tests.Error(errors.New("Returned the wrong connection"))
		return
	}
	connection.ID = testConnection.ID
}

func TestGetConnectionByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/connections/" + connection.ID.String())
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testConnection := types.Connection{}
	if err := json.Unmarshal(body, &testConnection); err != nil {
		tests.Error(err)
		return
	}
	if testConnection.ID != connection.ID {
		tests.Error(errors.New("Returned the wrong connection"))
		return
	}
}

func TestGetConnections(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/connections")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil && body != nil {
		tests.Error(err)
		return
	}
	testConnections := types.PaginatedResponse{}.Data
	if err := json.Unmarshal(body, &testConnections); err != nil {
		tests.Error(err)
		return
	}
	if testConnections == nil {
		tests.Error(errors.New("Failed to return paginated connections"))
		return
	}
}

func TestPutConnectionsByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	connection.Meta = []byte("{\"Updated?\":\"YES\"}")
	data, err := json.Marshal(connection)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/connections/"+connection.ID.String(), bytes.NewBuffer(data))
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
	testConnection := types.Connection{}
	if err := json.Unmarshal(body, &testConnection); err != nil {
		tests.Error(err)
		return
	}
	if testConnection.Meta == nil {
		tests.Error(errors.New("Wrong connection, failed to update"))
		return
	}
	connection.ID = testConnection.ID
}

func TestDeleteConnectionsByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", os.Getenv("TESTING_URL")+"/connections/"+connection.ID.String(), nil)
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
	testConnection := types.Connection{}
	if err := json.Unmarshal(body, &testConnection); err != nil {
		tests.Error(err)
		return
	}
}
