package router

import (
	"github.com/tespo/buddha/handlers"
	"github.com/tespo/satya/v2/types"
)

//
// VoiceCommandRoutes are routes for Google Home
// Alexa, and Siri
//
var VoiceCommandRoutes = map[string][]types.Route{
	"google": {
		{
			Name:        "Fulfillment for google home voice commands",
			Method:      "POST",
			Pattern:     "/google/fulfillment",
			HandlerFunc: handlers.GoogleFulfillment,
		},
	},
	"alexa": {
		{
			Name:        "Fulfillment for alexa voice commands",
			Method:      "POST",
			Pattern:     "/alexa/fulfillment",
			HandlerFunc: handlers.AlexaFulfillment,
		},
	},
}
