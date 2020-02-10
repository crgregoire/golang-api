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

var pod = types.Pod{
	ID:        uuid.NewV4(),
	Name:      "Integration Testing Pod",
	Slug:      "integration-testing-pod",
	Color:     "#000000",
	Cells:     31,
	LabelTall: "gettespo.com/labeltall",
	LabelWide: "gettespo.com/labelwide",
	Meta:      nil,
}

func TestCreatePod(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	data, err := json.Marshal(pod)
	if err != nil {
		tests.Error(err)
		return
	}
	response, err := http.Post(os.Getenv("TESTING_URL")+"/pods", "application/json", bytes.NewBuffer(data))
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testPod := types.Pod{}
	if err := json.Unmarshal(body, &testPod); err != nil {
		tests.Error(err)
		return
	}
	if testPod.Name != pod.Name {
		tests.Error(errors.New("Returned the wrong pod"))
		return
	}
	pod.ID = testPod.ID
}

func TestGetPod(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/pods/" + pod.ID.String())
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testPod := types.Pod{}
	if err := json.Unmarshal(body, &testPod); err != nil {
		tests.Error(err)
		return
	}
	if testPod.ID != pod.ID {
		tests.Error(errors.New("Returned the wrong pod"))
		return
	}
}

func TestGetPods(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/pods")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil && body != nil {
		tests.Error(err)
		return
	}
	testPods := types.PaginatedResponse{}.Data
	if err := json.Unmarshal(body, &testPods); err != nil {
		tests.Error(err)
		return
	}
	if testPods == nil {
		tests.Error(errors.New("Failed to return paginated pods"))
		return
	}
}

func TestPutPodsByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	pod.Meta = []byte("{\"Updated?\":\"YES\"}")
	data, err := json.Marshal(pod)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/pods/"+pod.ID.String(), bytes.NewBuffer(data))
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
	testPod := types.Pod{}
	if err := json.Unmarshal(body, &testPod); err != nil {
		tests.Error(err)
		return
	}
	if testPod.Meta == nil {
		tests.Error(errors.New("Wrong pod, failed to update"))
		return
	}
	pod.ID = testPod.ID
}

func TestDeletePodsByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", os.Getenv("TESTING_URL")+"/pods/"+pod.ID.String(), nil)
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
	testPod := types.Pod{}
	if err := json.Unmarshal(body, &testPod); err != nil {
		tests.Error(err)
		return
	}
}
