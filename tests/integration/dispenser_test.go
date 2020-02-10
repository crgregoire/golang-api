package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/tespo/satya/v2/types"
)

var dispenser = types.Dispenser{
	Serial:  "dispenser-11218" + randomString(5),
	Name:    "Integration testing dispenser",
	Network: "Aardvark 5G",
}

func TestCreateDispenser(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	data, err := json.Marshal(dispenser)
	if err != nil {
		tests.Error(err)
		return
	}
	response, err := http.Post(os.Getenv("TESTING_URL")+"/dispensers", "application/json", bytes.NewBuffer(data))
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testDispenser := types.Dispenser{}
	if err := json.Unmarshal(body, &testDispenser); err != nil {
		tests.Error(err)
		return
	}
	if testDispenser.Serial != dispenser.Serial {
		tests.Error(errors.New("Returned the wrong dispenser"))
		return
	}
	dispenser.ID = testDispenser.ID
}

func TestGetDispenserByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/dispensers/" + dispenser.ID.String())
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testDispenser := types.Dispenser{}
	if err := json.Unmarshal(body, &testDispenser); err != nil {
		tests.Error(err)
		return
	}
	if testDispenser.ID != dispenser.ID {
		tests.Error(errors.New("Returned the wrong dispenser"))
		return
	}
}

func TestGetDispensers(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/dispensers")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil && body != nil {
		tests.Error(err)
		return
	}
	testDispensers := types.PaginatedResponse{}.Data
	if err := json.Unmarshal(body, &testDispensers); err != nil {
		tests.Error(err)
		return
	}
	if testDispensers == nil {
		tests.Error(errors.New("Failed to return paginated dispensers"))
		return
	}
}

func TestPutDispensersByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	dispenser.Meta = []byte("{\"Updated?\":\"YES - TEST\"}")
	data, err := json.Marshal(dispenser)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/dispensers/"+dispenser.ID.String(), bytes.NewBuffer(data))
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
	testDispenser := types.Dispenser{}
	if err := json.Unmarshal(body, &testDispenser); err != nil {
		tests.Error(err)
		return
	}
	if testDispenser.Meta == nil {
		tests.Error(errors.New("Wrong dispenser, failed to update"))
		return
	}
	dispenser.ID = testDispenser.ID
}

func TestDeleteDispensersByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", os.Getenv("TESTING_URL")+"/dispensers/"+dispenser.ID.String(), nil)
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
	testDispenser := types.Dispenser{}
	if err := json.Unmarshal(body, &testDispenser); err != nil {
		tests.Error(err)
		return
	}
}

func randomString(len int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}
