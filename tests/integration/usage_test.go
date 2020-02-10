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

var usage = types.Usage{
	ID:          uuid.NewV4(),
	RegimenID:   uuid.FromStringOrNil("8ba3049b-17a1-4eae-b2bc-db7d18596d28"),
	DispenserID: uuid.FromStringOrNil("142201c2-0c5f-4650-8c99-fc233412e030"),
}

var usage2 = types.Usage{
	ID:          uuid.NewV4(),
	RegimenID:   uuid.FromStringOrNil("8ba3049b-17a1-4eae-b2bc-db7d18596d28"),
	DispenserID: uuid.FromStringOrNil("142201c2-0c5f-4650-8c99-fc233412e030"),
	Meta:        []byte("{\"Updated?\":\"NO\"}"),
}

func TestCreateUsage(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	data, err := json.Marshal(usage)
	if err != nil {
		tests.Error(err)
		return
	}
	response, err := http.Post(os.Getenv("TESTING_URL")+"/usages", "application/json", bytes.NewBuffer(data))
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testUsage := types.Usage{}
	if err := json.Unmarshal(body, &testUsage); err != nil {
		tests.Error(err)
		return
	}
	if testUsage.ID != usage.ID {
		tests.Error(errors.New("Returned the wrong usage"))
		return
	}
	usage.ID = testUsage.ID
}

func TestGetUsagesByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/usages/" + usage.ID.String())
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testUsage := types.Usage{}
	if err := json.Unmarshal(body, &testUsage); err != nil {
		tests.Error(err)
		return
	}
	if testUsage.ID != usage.ID {
		tests.Error(errors.New("Returned the wrong usage"))
		return
	}
}

func TestGetUsages(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/usages")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil && body != nil {
		tests.Error(err)
		return
	}
	testUsages := types.PaginatedResponse{}.Data
	if err := json.Unmarshal(body, &testUsages); err != nil {
		tests.Error(err)
		return
	}
	if testUsages == nil {
		tests.Error(errors.New("Failed to return paginated usages"))
		return
	}
}

func TestPutUsagesByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	usage.Meta = []byte("{\"Updated?\":\"YES\"}")
	data, err := json.Marshal(usage)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/usages/"+usage.ID.String(), bytes.NewBuffer(data))
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
	testUsage := types.Usage{}
	if err := json.Unmarshal(body, &testUsage); err != nil {
		tests.Error(err)
		return
	}
	if testUsage.Meta == nil {
		tests.Error(errors.New("Wrong usage, failed to update"))
		return
	}
	usage.ID = testUsage.ID
}

func TestDeleteUsagesByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", os.Getenv("TESTING_URL")+"/usages/"+usage.ID.String(), nil)
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
	testUsage := types.Usage{}
	if err := json.Unmarshal(body, &testUsage); err != nil {
		tests.Error(err)
		return
	}
}

func TestGetUserUsages(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/user/usages/")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testUsage := types.Usage{}
	if err := json.Unmarshal(body, &testUsage); err != nil {
		tests.Error(err)
		return
	}
	if testUsage.ID == usage.ID {
		tests.Error(errors.New("Returned the wrong usage"))
		return
	}
}

func TestGetUserUsagesByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/user/usages/" + usage.ID.String())
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testUsage := types.Usages{}
	if err := json.Unmarshal(body, &testUsage); err != nil {
		tests.Error(err)
		return
	}
	if testUsage == nil {
		tests.Error(errors.New("Returned the wrong usage"))
		return
	}
}

func TestPutUserUsagesByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	usage2.Meta = []byte("{\"Updated?\":\"YES\"}")
	data, err := json.Marshal(usage2)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/user/usages/"+usage2.ID.String(), bytes.NewBuffer(data))
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
	testUsage := types.Usages{}
	if err := json.Unmarshal(body, &testUsage); err != nil {
		tests.Error(err)
		return
	}
	if testUsage == nil {
		tests.Error(errors.New("Wrong usage, failed to update"))
		return
	}
}
