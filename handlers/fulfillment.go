package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/context"
	"github.com/tespo/buddha/db"
	"github.com/tespo/buddha/util"
	"github.com/tespo/satya/v2/types"
)

//
// GoogleFulfillment handles the voice commands webhook calls
//
func GoogleFulfillment(w http.ResponseWriter, r *http.Request) {
	request := types.GoogleHomeRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}
	response := types.GoogleHomeResponse{}
	switch request.Inputs[0].Intent {
	case "action.devices.EXECUTE":
		go googleDispensePod(r, request)
	case "Dispense":
		// err = dispensePod(r, request)
	case "action.devices.SYNC":
		response.Payload = types.ResponsePayload{
			AgentUserID: "test",
			Devices: []types.GoogleHomeDevice{
				{
					ID:     "tespo dispenser",
					Type:   "action.devices.types.MICROWAVE",
					Traits: []string{"action.devices.traits.StartStop"},
					Name: types.GoogleHomeDeviceName{
						DefaultNames: []string{"Tespo Connect Dispenser", "Vitamin Dispenser"},
						Name:         "Dispenser",
						Nicknames:    []string{"Kitchen Dispenser"},
					},
					WillReportState: false,
					Attributes: types.GoogleHomeDeviceAttributes{
						Pausable: false,
						AvailableZones: []string{
							"Kitchen",
							"Bathroom",
						},
					},
				},
			},
		}
	}
	response.RequestID = request.RequestID
	util.JSONResponder(w, response)
}

func googleDispensePod(r *http.Request, request types.GoogleHomeRequest) error {
	id, ok := context.GetOk(r, "user_id")
	if !ok {
		return errors.New("Cannot process token claims")
	}
	user := types.User{
		ID: uuid.FromStringOrNil(id.(string)),
	}
	connections := types.Connections{}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := user.GetByID(db, user.ID); err != nil {
		return err
	}
	dispensers, err := connections.GetAccountDispensers(db, user.AccountID)
	if err != nil {
		return err
	}
	dispenser := dispensers[0]
	payload := types.Payload{
		Customer: types.PayloadCustomer{
			ID: user.AccountID.String(),
		},
		Dispenser: types.PayloadDispenser{
			Serial: dispenser.Serial,
			Name:   dispenser.Name,
		},
	}
	_, err = util.TriggerLambda(payload, "DispenserDispense")
	if err != nil {
		return err
	}
	return nil
}

//
// AlexaFulfillment will handle the different intent
// requests that come in from Alexa
//
func AlexaFulfillment(w http.ResponseWriter, r *http.Request) {
	var request types.AlexaRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		util.ErrorResponder(w, http.StatusBadRequest, err)
		return
	}
	var response types.AlexaResponse
	switch request.Directive.Header.Name {
	case "Discover":
		response = alexaDiscoverResponse(r, request)
	case "TurnOn":
		response = *alexaDispensePod(r, request)
	}
	util.JSONResponder(w, &response)
}

func alexaDiscoverResponse(r *http.Request, request types.AlexaRequest) types.AlexaResponse {
	response := types.AlexaResponse{}
	response.Event = types.AlexaDiscoverResponseEvent{
		Header: types.AlexaHeader{
			MessageID:      request.Directive.Header.MessageID,
			Namespace:      "Alexa.Discovery",
			Name:           "Discover.Response",
			PayloadVersion: "3",
		},
		Payload: types.AlexaDiscoverResponseEventPayload{
			Endpoints: []types.AlexaDiscoverResponseEventPayloadEndpoint{
				{
					EndpointID:        "dispense321",
					FriendlyName:      "Dispenser",
					Description:       "Tespo Connect Dispenser",
					ManufacturerName:  "Tespo",
					DisplayCategories: []string{"OTHER"},
					Cookie:            map[string]interface{}{},
					Capabilities: []types.AlexaEndpointCapabilities{
						{
							Type:      "AlexaInterface",
							Interface: "Alexa.PowerController",
							Version:   "3",
							Properties: types.AlexaCapabilityProperites{
								ProactivelyReported: true,
								Retrievable:         true,
								Supported: []map[string]string{
									{
										"name": "powerState",
									},
								},
							},
						},
					},
					AdditionalAttributes: types.AlexaAdditionalAttributes{
						Manufacturer:     "Tepspo",
						Model:            "Model 2",
						SerialNumber:     "123",
						FirmwareVersion:  "1",
						SoftwareVersion:  "1",
						CustomIdentifier: "1",
					},
				},
			},
		},
	}
	data, _ := json.Marshal(response)
	fmt.Printf("%s", data)
	return response
}

func alexaDispensePod(r *http.Request, request types.AlexaRequest) *types.AlexaResponse {
	id, ok := context.GetOk(r, "user_id")
	if !ok {
		return nil
	}
	user := types.User{
		ID: uuid.FromStringOrNil(id.(string)),
	}
	connections := types.Connections{}
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := user.GetByID(db, user.ID); err != nil {
		return nil
	}
	dispensers, err := connections.GetAccountDispensers(db, user.AccountID)
	if err != nil {
		return nil
	}
	dispenser := dispensers[0]
	response := types.AlexaResponse{
		Context: types.AlexaContext{
			Properties: []types.AlexaHeader{
				{
					Namespace:    "Alexa.ToggleController",
					Name:         "toggleState",
					Value:        "ON",
					TimeOfSample: time.Now(),
				},
			},
		},
		Event: types.AlexaDiscoverResponseEvent{
			Header: types.AlexaHeader{
				Namespace:      "Alexa",
				Name:           "Response",
				PayloadVersion: "3",
				MessageID:      request.Directive.Header.MessageID,
			},
			Payload: types.AlexaDiscoverResponseEventPayload{},
		},
	}
	payload := types.Payload{
		Customer: types.PayloadCustomer{
			ID: user.AccountID.String(),
		},
		Dispenser: types.PayloadDispenser{
			Serial: dispenser.Serial,
			Name:   dispenser.Name,
		},
	}
	_, err = util.TriggerLambda(payload, "DispenserDispense")
	if err != nil {
		panic(err)
	}
	return &response
}
