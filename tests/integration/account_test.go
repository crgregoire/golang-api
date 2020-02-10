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

var account = types.Account{
	ID:   uuid.NewV4(),
	Name: "Integration testing account",
}

var account2 = types.Account{
	ID:   uuid.FromStringOrNil("22b5123d-9cee-4701-b15b-8c9078142666"),
	Name: "Billy Bob Thornton",
	Meta: []byte("{\"Updated?\":\"No\"}"),
}

func TestCreateAccount(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	data, err := json.Marshal(account)
	if err != nil {
		tests.Error(err)
		return
	}
	response, err := http.Post(os.Getenv("TESTING_URL")+"/accounts", "application/json", bytes.NewBuffer(data))
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testAccount := types.Account{}
	if err := json.Unmarshal(body, &testAccount); err != nil {
		tests.Error(err)
		return
	}
	if testAccount.Name != account.Name {
		tests.Error(errors.New("Returned the wrong account"))
		return
	}
	account.ID = testAccount.ID
}

func TestGetAccountsByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/accounts/" + account2.ID.String())
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testAccount := types.Account{}
	if err := json.Unmarshal(body, &testAccount); err != nil {
		tests.Error(err)
		return
	}
	if testAccount.ID != account2.ID {
		tests.Error(errors.New("Returned the wrong account"))
		return
	}
}

func TestGetAccounts(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/accounts")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil && body != nil {
		tests.Error(err)
		return
	}
	testAccount := types.PaginatedResponse{}.Data
	if err := json.Unmarshal(body, &testAccount); err != nil {
		tests.Error(err)
		return
	}
	if testAccount == nil {
		tests.Error(errors.New("Failed to return paginated accounts"))
		return
	}
}

func TestPutAccountsByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	account2.Meta = []byte("{\"Updated?\":\"YES\"}")
	data, err := json.Marshal(account2)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/accounts/"+account2.ID.String(), bytes.NewBuffer(data))
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
	testAccount := types.Account{}
	if err := json.Unmarshal(body, &testAccount); err != nil {
		tests.Error(err)
		return
	}
	if testAccount.Meta == nil {
		tests.Error(errors.New("Wrong account, failed to update"))
		return
	}
	account2.ID = testAccount.ID
}

func TestDeleteAccountsByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", os.Getenv("TESTING_URL")+"/accounts/"+account.ID.String(), nil)
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
	testAccount := types.Account{}
	if err := json.Unmarshal(body, &testAccount); err != nil {
		tests.Error(err)
		return
	}
}
