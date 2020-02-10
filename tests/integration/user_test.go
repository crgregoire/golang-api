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

var user = types.User{
	ID:        uuid.NewV4(),
	AccountID: uuid.FromStringOrNil("d8e4c5dc-9767-41bd-b802-060e80d83867"),
	FirstName: "Integration",
	LastName:  "testing user",
}
var user2 = types.User{
	ID:        uuid.NewV4(),
	AccountID: uuid.FromStringOrNil("d8e4c5dc-9767-41bd-b802-060e80d83867"),
	FirstName: "Integration",
	LastName:  "testing user",
	Meta:      []byte("{\"Updated?\":\"No\"}"),
}

func TestCreateUser(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	data, err := json.Marshal(user)
	if err != nil {
		tests.Error(err)
		return
	}
	response, err := http.Post(os.Getenv("TESTING_URL")+"/users", "application/json", bytes.NewBuffer(data))
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testUser := types.User{}
	if err := json.Unmarshal(body, &testUser); err != nil {
		tests.Error(err)
		return
	}
	if testUser.ID != user.ID {
		tests.Error(errors.New("Returned the wrong user"))
		return
	}
	user.ID = testUser.ID
}

func TestGetUser(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/user")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testUser := types.User{}
	if err := json.Unmarshal(body, &testUser); err != nil {
		tests.Error(err)
		return
	}
	if testUser.ID != uuid.Nil {
		tests.Error(errors.New("User response is not scoped"))
		return
	}
}

func TestPutUser(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	user2.Meta = []byte("{\"Updated?\":\"Yes\"}")
	user2.CognitoID = uuid.NewV4()
	data, err := json.Marshal(user2)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/user", bytes.NewBuffer(data))
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
	testUser := types.User{}
	if err := json.Unmarshal(body, &testUser); err != nil {
		tests.Error(err)
		return
	}
	if testUser.Meta == nil {
		tests.Error(errors.New("Wrong account, failed to update"))
		return
	}
	user2.ID = testUser.ID
}

func TestGetUsers(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/users")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil && body != nil {
		tests.Error(err)
		return
	}
	testUsers := types.PaginatedResponse{}.Data
	if err := json.Unmarshal(body, &testUsers); err != nil {
		tests.Error(err)
		return
	}
	if testUsers == nil {
		tests.Error(errors.New("Failed to return paginated users"))
		return
	}
}

func TestPutUsersByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	user2.Meta = []byte("{\"Updated?\":\"YES\"}")
	data, err := json.Marshal(user2)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/user/"+user2.ID.String(), bytes.NewBuffer(data))
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
	testUser := types.User{}
	if err := json.Unmarshal(body, &testUser); err != nil {
		tests.Error(err)
		return
	}
	if testUser.Meta == nil {
		tests.Error(errors.New("Wrong user, failed to update"))
		return
	}
	user2.ID = testUser.ID
}

func TestPutUsersByAccountAndID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	user2.Meta = []byte("{\"Updated?\":\"YES\"}")
	data, err := json.Marshal(user2)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/account/users/"+user2.ID.String(), bytes.NewBuffer(data))
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
	testUser := types.User{}
	if err := json.Unmarshal(body, &testUser); err != nil {
		tests.Error(err)
		return
	}
	if testUser.Meta == nil {
		tests.Error(errors.New("Wrong user2, failed to update"))
		return
	}
	user2.ID = testUser.ID
}

func TestDeleteUsersByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", os.Getenv("TESTING_URL")+"/accounts/"+user.AccountID.String()+"/users/"+user.ID.String(), nil)
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
	testUser := types.User{}
	if err := json.Unmarshal(body, &testUser); err != nil {
		tests.Error(err)
		return
	}
}
