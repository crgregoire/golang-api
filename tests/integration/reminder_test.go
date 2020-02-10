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

var reminder = types.Reminder{
	ID:        uuid.NewV4(),
	UserID:    uuid.FromStringOrNil("8c8aa229-3959-4a40-bbe6-67c2eeace5cb"),
	RegimenID: uuid.FromStringOrNil("8ba3049b-17a1-4eae-b2bc-db7d18596d28"),
	Minute:    215,
}

var reminder2 = types.Reminder{
	ID:        uuid.FromStringOrNil("22b5123d-9cee-4701-b15b-8c9078142666"),
	UserID:    uuid.FromStringOrNil("8c8aa229-3959-4a40-bbe6-67c2eeace5cb"),
	RegimenID: uuid.FromStringOrNil("8ba3049b-17a1-4eae-b2bc-db7d18596d28"),
	Minute:    262,
	Meta:      []byte("{\"Updated?\":\"No\"}"),
}

func TestCreateReminder(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	data, err := json.Marshal(reminder)
	if err != nil {
		tests.Error(err)
		return
	}
	response, err := http.Post(os.Getenv("TESTING_URL")+"/regimens/"+reminder.RegimenID.String()+"/reminders", "application/json", bytes.NewBuffer(data))
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tests.Error(err)
		return
	}
	testReminder := types.Reminders{}
	if err := json.Unmarshal(body, &testReminder); err != nil {
		tests.Error(err)
		return
	}
	if testReminder == nil {
		tests.Error(errors.New("Failed to create reminder"))
		return
	}
}

func TestGetUserReminders(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/user/reminders")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil && body != nil {
		tests.Error(err)
		return
	}
	testReminders := types.PaginatedResponse{}.Data
	if err := json.Unmarshal(body, &testReminders); err != nil {
		tests.Error(err)
		return
	}
	if testReminders == nil {
		tests.Error(errors.New("Failed to return paginated reminders"))
		return
	}
}

func TestGetRemindersByRegimenID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	response, err := http.Get(os.Getenv("TESTING_URL") + "/regimens/" + reminder.RegimenID.String() + "/reminders")
	if err != nil {
		tests.Error(err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil && body != nil {
		tests.Error(err)
		return
	}
	testReminder := types.Reminders{}
	if testReminder == nil {
		tests.Error(errors.New("Failed to return reminders"))
		return
	}
}

func TestPutRemindersByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	reminder2.Meta = []byte("{\"Updated?\":\"YES\"}")
	data, err := json.Marshal(reminder2)
	if err != nil {
		tests.Error(err)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", os.Getenv("TESTING_URL")+"/regimens/"+reminder2.RegimenID.String()+"/reminders/"+reminder2.ID.String(), bytes.NewBuffer(data))
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
	testReminder := types.Reminders{}
	if err := json.Unmarshal(body, &testReminder); err != nil {
		tests.Error(err)
		return
	}
	if testReminder[0].Meta == nil {
		tests.Error(errors.New("Wrong reminder, failed to update"))
		return
	}
	reminder2.ID = testReminder[0].ID
}

func TestDeleteRemindersByID(tests *testing.T) {
	if testing.Short() {
		tests.Skip()
	}
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", os.Getenv("TESTING_URL")+"/regimens/"+reminder.RegimenID.String()+"/reminders/"+reminder.ID.String(), nil)
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
	testReminder := types.Reminders{}
	if err := json.Unmarshal(body, &testReminder); err != nil {
		tests.Error(err)
		return
	}
}
