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

var barcodeID = uuid.FromStringOrNil("62012ee1-f3c2-4a4b-a0ab-8686e3d173e4")

var insertion = types.Insertion{
	ID:          uuid.NewV4(),
	DispenserID: uuid.FromStringOrNil("142201c2-0c5f-4650-8c99-fc233412e030"),
	RegimenID:   uuid.FromStringOrNil("8ba3049b-17a1-4eae-b2bc-db7d18596d28"),
	BarcodeID:   &barcodeID,
	Flags:       5,
	Servings:    31,
	Meta:        []byte("{\"Name\":\"Insertion Integration Test\"}"),
}

func TestCreateInsertion(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	data, err := json.Marshal(insertion)
	if err != nil {
		tests.Error(err)
		return
	}
	response, err := http.Post(os.Getenv("TESTING_URL")+"/insertions", "application/json", bytes.NewBuffer(data))
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testInsertion := types.Insertion{}
	if err := json.Unmarshal(body, &testInsertion); err != nil {
		tests.Error(err)
		return
	}
	if testInsertion.ID != insertion.ID {
		tests.Error(errors.New("Returned the wrong insertion"))
		return
	}
	insertion.ID = testInsertion.ID
}

func TestGetInsertionByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/insertions/" + insertion.ID.String())
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testInsertion := types.Insertion{}
	if err := json.Unmarshal(body, &testInsertion); err != nil {
		tests.Error(err)
		return
	}
	if testInsertion.ID != insertion.ID {
		tests.Error(errors.New("Returned the wrong insertion"))
		return
	}
}

func TestGetInsertions(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/insertions")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil && body != nil {
		tests.Error(err)
		return
	}
	testInsertions := types.PaginatedResponse{}.Data
	if err := json.Unmarshal(body, &testInsertions); err != nil {
		tests.Error(err)
		return
	}
	if testInsertions == nil {
		tests.Error(errors.New("Failed to return paginated insertions"))
		return
	}
}

func TestPutInsertionsByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	data, err := json.Marshal(insertion)
	if err != nil {
		tests.Error(err)
		return
	}
	response, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/insertions/"+insertion.ID.String(), bytes.NewBuffer(data))
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testInsertion := types.Insertion{}
	if err := json.Unmarshal(body, &testInsertion); err != nil {
		tests.Error(err)
		return
	}
	if testInsertion.ID != insertion.ID {
		tests.Error(errors.New("Could not update Insertions"))
		return
	}
}

func TestDeleteInsertionsByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", os.Getenv("TESTING_URL")+"/insertions/"+insertion.ID.String(), nil)
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
	testInsertion := types.Insertion{}
	if err := json.Unmarshal(body, &testInsertion); err != nil {
		tests.Error(err)
		return
	}
}
