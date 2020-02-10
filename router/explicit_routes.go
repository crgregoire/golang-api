package router

import (
	"github.com/tespo/buddha/handlers"
	"github.com/tespo/satya/v2/types"
)

//
// ExplicitRoutes are all routes that require
// permissions to explicitly access database resources
//
var ExplicitRoutes = []types.Route{
	{
		Name:        "Get Users",
		Method:      "GET",
		Pattern:     "/users",
		HandlerFunc: handlers.GetUsers,
	},
	{
		Name:        "Get User By ID",
		Method:      "GET",
		Pattern:     "/user/{user_id}",
		HandlerFunc: handlers.GetUserByID,
	},
	{
		Name:        "Post User",
		Method:      "POST",
		Pattern:     "/users",
		HandlerFunc: handlers.PostUsers,
	},
	{
		Name:        "Put User By ID",
		Method:      "PUT",
		Pattern:     "/user/{user_id}",
		HandlerFunc: handlers.PutUsersByID,
	},
	{
		Name:        "Delete User By ID",
		Method:      "DELETE",
		Pattern:     "/user/{user_id}",
		HandlerFunc: handlers.DeleteUsersByID,
	},
	// Dev paths start
	{
		Name:        "Get Accounts",
		Method:      "GET",
		Pattern:     "/accounts",
		HandlerFunc: handlers.GetAccounts,
	},
	{
		Name:        "Get Accounts By ID",
		Method:      "GET",
		Pattern:     "/accounts/{account_id}",
		HandlerFunc: handlers.GetAccountByID,
	},
	{
		Name:        "Post Accounts",
		Method:      "POST",
		Pattern:     "/accounts",
		HandlerFunc: handlers.PostAccount,
	},
	{
		Name:        "Put Accounts By ID",
		Method:      "PUT",
		Pattern:     "/accounts/{account_id}",
		HandlerFunc: handlers.PutAccountByID,
	},
	{
		Name:        "Delete Accounts By ID",
		Method:      "DELETE",
		Pattern:     "/accounts/{account_id}",
		HandlerFunc: handlers.DeleteAccountByID,
	},
	{
		Name:        "Get Users By Account ID",
		Method:      "GET",
		Pattern:     "/accounts/{account_id}/users",
		HandlerFunc: handlers.GetUsersByAccountID,
	},
	{
		Name:        "Post Users By Account ID",
		Method:      "POST",
		Pattern:     "/accounts/{account_id}/users",
		HandlerFunc: handlers.PostUsersByAccountID,
	},
	{
		Name:        "Put Users By Account ID And User ID",
		Method:      "PUT",
		Pattern:     "/accounts/{account_id}/users/{user_id}",
		HandlerFunc: handlers.PutUsersByAccountIDAndUserID,
	},
	{
		Name:        "Delete Users By Account ID And User ID",
		Method:      "DELETE",
		Pattern:     "/accounts/{account_id}/users/{user_id}",
		HandlerFunc: handlers.DeleteUsersByAccountIDAndUserID,
	},
	{
		Name:        "Get Pods",
		Method:      "GET",
		Pattern:     "/pods",
		HandlerFunc: handlers.GetPods,
	},
	{
		Name:        "Get Pods By ID",
		Method:      "GET",
		Pattern:     "/pods/{pod_id}",
		HandlerFunc: handlers.GetPodsByID,
	},
	{
		Name:        "Post Pods By ID",
		Method:      "POST",
		Pattern:     "/pods",
		HandlerFunc: handlers.PostPods,
	},
	{
		Name:        "Put Pods By ID",
		Method:      "PUT",
		Pattern:     "/pods/{pod_id}",
		HandlerFunc: handlers.PutPodsByID,
	},
	{
		Name:        "Delete Pods By ID",
		Method:      "DELETE",
		Pattern:     "/pods/{pod_id}",
		HandlerFunc: handlers.DeletePodsByID,
	},
	{
		Name:        "Get Barcodes",
		Method:      "GET",
		Pattern:     "/barcodes",
		HandlerFunc: handlers.GetBarcodes,
	},
	{
		Name:        "Get Barcodes By ID",
		Method:      "GET",
		Pattern:     "/barcodes/{barcode_id}",
		HandlerFunc: handlers.GetBarcodesByID,
	},
	{
		Name:        "Post Barcode",
		Method:      "POST",
		Pattern:     "/barcodes",
		HandlerFunc: handlers.PostBarcodes,
	},
	{
		Name:        "Put Barcodes By ID",
		Method:      "PUT",
		Pattern:     "/barcodes/code/{code}",
		HandlerFunc: handlers.PutBarcodesByCode,
	},
	{
		Name:        "Put Barcodes By ID",
		Method:      "PUT",
		Pattern:     "/barcodes/{barcode_id}",
		HandlerFunc: handlers.PutBarcodesByID,
	},
	{
		Name:        "Delete Barcodes By ID",
		Method:      "DELETE",
		Pattern:     "/barcodes/{barcode_id}",
		HandlerFunc: handlers.DeleteBarcodesByID,
	},
	{
		Name:        "Get Dispensers",
		Method:      "GET",
		Pattern:     "/dispensers",
		HandlerFunc: handlers.GetDispensers,
	},
	{
		Name:        "Get Dispensers By ID",
		Method:      "GET",
		Pattern:     "/dispensers/{dispenser_id}",
		HandlerFunc: handlers.GetDispensersByID,
	},
	{
		Name:        "Post Dispensers",
		Method:      "POST",
		Pattern:     "/dispensers",
		HandlerFunc: handlers.PostDispensers,
	},
	{
		Name:        "Put Dispensers By ID",
		Method:      "PUT",
		Pattern:     "/dispensers/{dispenser_id}",
		HandlerFunc: handlers.PutDispensersByID,
	},
	{
		Name:        "Delete Dispensers By ID",
		Method:      "DELETE",
		Pattern:     "/dispensers/{dispenser_id}",
		HandlerFunc: handlers.DeleteDispensersByID,
	},
	{
		Name:        "Get Connections",
		Method:      "GET",
		Pattern:     "/connections",
		HandlerFunc: handlers.GetConnections,
	},
	{
		Name:        "Get Connections By ID",
		Method:      "GET",
		Pattern:     "/connections/{connection_id}",
		HandlerFunc: handlers.GetConnectionsByID,
	},
	{
		Name:        "Post Connections",
		Method:      "POST",
		Pattern:     "/connections",
		HandlerFunc: handlers.PostConnections,
	},
	{
		Name:        "Put Connections By ID",
		Method:      "PUT",
		Pattern:     "/connections/{connection_id}",
		HandlerFunc: handlers.PutConnectionsByID,
	},
	{
		Name:        "Delete Connections By ID",
		Method:      "DELETE",
		Pattern:     "/connections/{connection_id}",
		HandlerFunc: handlers.DeleteConnectionsByID,
	},
	{
		Name:        "Get Insertions",
		Method:      "GET",
		Pattern:     "/insertions",
		HandlerFunc: handlers.GetInsertions,
	},
	{
		Name:        "Get Insertions By ID",
		Method:      "GET",
		Pattern:     "/insertions/{insertion_id}",
		HandlerFunc: handlers.GetInsertionByID,
	},
	{
		Name:        "Post Insertions",
		Method:      "POST",
		Pattern:     "/insertions",
		HandlerFunc: handlers.PostInsertion,
	},
	{
		Name:        "Put Insertions By ID",
		Method:      "PUT",
		Pattern:     "/insertions/{insertion_id}",
		HandlerFunc: handlers.PutInsertionByID,
	},
	{
		Name:        "Delete Insertions By ID",
		Method:      "DELETE",
		Pattern:     "/insertions/{insertion_id}",
		HandlerFunc: handlers.DeleteInsertionsByID,
	},
	{
		Name:        "Get Regimen",
		Method:      "GET",
		Pattern:     "/regimens",
		HandlerFunc: handlers.GetRegimen,
	},
	{
		Name:        "Get Regimen By ID",
		Method:      "GET",
		Pattern:     "/regimens/{regimen_id}",
		HandlerFunc: handlers.GetRegimenByID,
	},
	{
		Name:        "Put Regimen By ID",
		Method:      "PUT",
		Pattern:     "/regimens/{regimen_id}",
		HandlerFunc: handlers.PutRegimenByID,
	},
	{
		Name:        "Delete Regimen By ID",
		Method:      "DELETE",
		Pattern:     "/regimens/{regimen_id}",
		HandlerFunc: handlers.DeleteRegimenByID,
	},
	{
		Name:        "Get Usages",
		Method:      "GET",
		Pattern:     "/usages",
		HandlerFunc: handlers.GetUsages,
	},
	{
		Name:        "Get Usages By ID",
		Method:      "GET",
		Pattern:     "/usages/{usage_id}",
		HandlerFunc: handlers.GetUsagesByID,
	},
	{
		Name:        "Post Usages",
		Method:      "POST",
		Pattern:     "/usages",
		HandlerFunc: handlers.PostUsages,
	},
	{
		Name:        "Put Usages By ID",
		Method:      "PUT",
		Pattern:     "/usages/{usage_id}",
		HandlerFunc: handlers.PutUsagesByID,
	},
	{
		Name:        "Delete Usages By ID",
		Method:      "DELETE",
		Pattern:     "/usages/{usage_id}",
		HandlerFunc: handlers.DeleteUsagesByID,
	},
	{
		Name:        "Get Permissions",
		Method:      "GET",
		Pattern:     "/permissions",
		HandlerFunc: handlers.GetPermissions,
	},
	{
		Name:        "Get Permissions By ID",
		Method:      "GET",
		Pattern:     "/permissions/{permission_id}",
		HandlerFunc: handlers.GetPermissionsByID,
	},
	{
		Name:        "Post Permissions",
		Method:      "POST",
		Pattern:     "/permissions",
		HandlerFunc: handlers.PostPermissions,
	},
	{
		Name:        "Put Permissions By ID",
		Method:      "PUT",
		Pattern:     "/permissions/{permission_id}",
		HandlerFunc: handlers.PutPermissionsByID,
	},
	{
		Name:        "Delete Permissions By ID",
		Method:      "DELETE",
		Pattern:     "/permissions/{permission_id}",
		HandlerFunc: handlers.DeletePermissionsByID,
	},
	{
		Name:        "Get Roles",
		Method:      "GET",
		Pattern:     "/roles",
		HandlerFunc: handlers.GetRoles,
	},
	{
		Name:        "Get Roles By ID",
		Method:      "GET",
		Pattern:     "/roles/{role_id}",
		HandlerFunc: handlers.GetRolesByID,
	},
	{
		Name:        "Post Roles",
		Method:      "POST",
		Pattern:     "/role",
		HandlerFunc: handlers.PostRoles,
	},
	{
		Name:        "Put Roles By ID",
		Method:      "PUT",
		Pattern:     "/roles/{role_id}",
		HandlerFunc: handlers.PutRolesByID,
	},
	{
		Name:        "Delete Roles By ID",
		Method:      "DELETE",
		Pattern:     "/roles/{role_id}",
		HandlerFunc: handlers.DeleteRolesByID,
	},
	{
		Name:        "Add Permission to Role By ID",
		Method:      "GET",
		Pattern:     "/roles/{role_id}/permissions",
		HandlerFunc: handlers.GetRoleWithPermissions,
	},
	{
		Name:        "Add Permission to Role By ID",
		Method:      "DELETE",
		Pattern:     "/roles/{role_id}/permission/{permission_id}",
		HandlerFunc: handlers.DeletePermissionFromRole,
	},
	{
		Name:        "Add Permission to Role By ID",
		Method:      "PUT",
		Pattern:     "/roles/{role_id}/{permission_id}",
		HandlerFunc: handlers.AddPermissionToRoleByID,
	},
}
