package router

import (
	"github.com/tespo/buddha/handlers"
	"github.com/tespo/satya/v2/types"
)

//
// LambdaRoutes are all routes that aws lambda
// will be calling
//
var LambdaRoutes = []types.Route{
	{
		Name:        "Dispenser Dispensed",
		Method:      "POST",
		Pattern:     "/dispenser/dispensed",
		HandlerFunc: handlers.DispenserDispensed,
	},
	{
		Name:        "Dispenser Inserted",
		Method:      "POST",
		Pattern:     "/dispenser/inserted",
		HandlerFunc: handlers.PodInserted,
	},
	{
		Name:        "Dispenser Connected",
		Method:      "POST",
		Pattern:     "/dispenser/connected",
		HandlerFunc: handlers.DispenserConnected,
	},
	{
		Name:        "Dispenser Disconnected",
		Method:      "POST",
		Pattern:     "/dispenser/disconnected",
		HandlerFunc: handlers.DispenserDisconnected,
	},
	{
		Name:        "Update User By External ID",
		Method:      "PUT",
		Pattern:     "/user/{external_id}",
		HandlerFunc: handlers.PutUsersByExternalID,
	},
}
