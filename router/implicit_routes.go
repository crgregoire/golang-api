package router

import (
	"github.com/tespo/buddha/handlers"
	"github.com/tespo/satya/v2/types"
)

//
// ImplicitRoutes are all routes that implicitly
// gather data based on the user/account in the
// jwt token
//
var ImplicitRoutes = map[string]types.Route{
	"account.info": {
		Name:        "Get Account",
		Method:      "GET",
		Pattern:     "/account",
		HandlerFunc: handlers.GetAccount,
	},
	"account.update": {
		Name:        "Put Account",
		Method:      "PUT",
		Pattern:     "/account",
		HandlerFunc: handlers.PutAccount,
	},
	"account.users": {
		Name:        "Get Users",
		Method:      "GET",
		Pattern:     "/account/users",
		HandlerFunc: handlers.GetAccountUsers,
	},
	"account.user": {
		Name:        "Get Users By ID",
		Method:      "GET",
		Pattern:     "/account/users/{user_id}",
		HandlerFunc: handlers.GetAccountUsersByUserID,
	},
	"account.update.user": {
		Name:        "Put Users By ID",
		Method:      "PUT",
		Pattern:     "/account/users/{user_id}",
		HandlerFunc: handlers.PutAccountUsersByUserID,
	},
	"account.create.users": {
		Name:        "Create Account User",
		Method:      "POST",
		Pattern:     "/account/users",
		HandlerFunc: handlers.CreateAccountUser,
	},
	"account.delete.users": {
		Name:        "Delete Account User",
		Method:      "DELETE",
		Pattern:     "/account/users/{user_id}",
		HandlerFunc: handlers.DeleteAccountUserByID,
	},
	"user.info": {
		Name:        "Get User",
		Method:      "GET",
		Pattern:     "/user",
		HandlerFunc: handlers.GetSelfUser,
	},
	"user.update": {
		Name:        "Put User",
		Method:      "PUT",
		Pattern:     "/user",
		HandlerFunc: handlers.PutSelfUser,
	},
	"account.connections": {
		Name:        "Get Connection",
		Method:      "GET",
		Pattern:     "/account/connections",
		HandlerFunc: handlers.GetAccountConnections,
	},
	"account.connection": {
		Name:        "Get Connection By ID",
		Method:      "GET",
		Pattern:     "/account/connection/{connection_id}",
		HandlerFunc: handlers.GetAccountConnectionByID,
	},
	"account.update.connection": {
		Name:        "Put Connection By ID",
		Method:      "PUT",
		Pattern:     "/account/connection/{connection_id}",
		HandlerFunc: handlers.PutAccountConnectionsByID,
	},
	"account.delete.connection": {
		Name:        "Delete Connection",
		Method:      "DELETE",
		Pattern:     "/account/connection/{connection_id}",
		HandlerFunc: handlers.DeleteAccountConnectionsByID,
	},
	"account.dispensers": {
		Name:        "Get Dispenser",
		Method:      "GET",
		Pattern:     "/account/dispensers",
		HandlerFunc: handlers.GetAccountDispensers,
	},
	"account.dispenser": {
		Name:        "Get Dispenser By ID",
		Method:      "GET",
		Pattern:     "/account/dispensers/{dispenser_id}",
		HandlerFunc: handlers.GetDispenserByID,
	},
	"account.update.dispenser": {
		Name:        "Put Dispenser By ID",
		Method:      "PUT",
		Pattern:     "/account/dispensers/{dispenser_id}",
		HandlerFunc: handlers.PutDispenserByID,
	},
	"account.delete.dispenser": {
		Name:        "Delete Dispenser",
		Method:      "DELETE",
		Pattern:     "/account/dispensers",
		HandlerFunc: handlers.DeleteDispenser,
	},
	"account.usages": {
		Name:        "Get Account Usages",
		Method:      "GET",
		Pattern:     "/account/usages",
		HandlerFunc: handlers.GetAccountUsages,
	},
	"account.usage": {
		Name:        "Get Account Usages By ID",
		Method:      "GET",
		Pattern:     "/account/usages/{usage_id}",
		HandlerFunc: handlers.GetAccountUsageByID,
	},
	"account.update.usage": {
		Name:        "Update Account Usages By ID",
		Method:      "PUT",
		Pattern:     "/account/usages/{usage_id}",
		HandlerFunc: handlers.PutAccountUsageByID,
	},
	"user.usages": {
		Name:        "Get User Usages",
		Method:      "GET",
		Pattern:     "/user/usages",
		HandlerFunc: handlers.GetUserUsages,
	},
	"user.usage": {
		Name:        "Get User Usages By ID",
		Method:      "GET",
		Pattern:     "/user/usages/{usage_id}",
		HandlerFunc: handlers.GetUserUsageByID,
	},
	"user.update.usage": {
		Name:        "Update User Usages By ID",
		Method:      "PUT",
		Pattern:     "/user/usages/{usage_id}",
		HandlerFunc: handlers.PutUserUsageByID,
	},
	"account.regimens": {
		Name:        "Get Regimens",
		Method:      "GET",
		Pattern:     "/account/regimens",
		HandlerFunc: handlers.GetAccountRegimens,
	},
	"user.regimens": {
		Name:        "Get Regimens",
		Method:      "GET",
		Pattern:     "/user/regimens",
		HandlerFunc: handlers.GetUserRegimens,
	},
	"account.regimen": {
		Name:        "Get Regimens By ID",
		Method:      "GET",
		Pattern:     "/account/regimens/{regimen_id}",
		HandlerFunc: handlers.GetAccountRegimensByID,
	},
	"user.regimen": {
		Name:        "Get Regimens By ID",
		Method:      "GET",
		Pattern:     "/user/regimens/{regimen_id}",
		HandlerFunc: handlers.GetUserRegimensByID,
	},
	"account.update.regimen": {
		Name:        "Put Regimens By ID",
		Method:      "PUT",
		Pattern:     "/account/regimens/{regimen_id}",
		HandlerFunc: handlers.PutAccountRegimensByID,
	},
	"account.delete.regimen": {
		Name:        "Delete Regimens By ID",
		Method:      "DELETE",
		Pattern:     "/account/regimens/{regimen_id}",
		HandlerFunc: handlers.DeleteAccountRegimenByID,
	},
	"user.reminder": {
		Name:        "Get Reminder",
		Method:      "GET",
		Pattern:     "/regimens/{regimen_id}/reminders",
		HandlerFunc: handlers.GetUserRemindersByRegimenID,
	},
	"user.reminders": {
		Name:        "Get Reminder",
		Method:      "GET",
		Pattern:     "/user/reminders",
		HandlerFunc: handlers.GetReminders,
	},
	"user.create.reminder": {
		Name:        "Post Reminder",
		Method:      "POST",
		Pattern:     "/regimens/{regimen_id}/reminders",
		HandlerFunc: handlers.PostReminder,
	},
	"user.update.reminder": {
		Name:        "Put Reminder By ID",
		Method:      "PUT",
		Pattern:     "/regimens/{regimen_id}/reminders/{reminder_id}",
		HandlerFunc: handlers.PutReminderByID,
	},
	"user.delete.reminder": {
		Name:        "Delete Reminder By ID",
		Method:      "DELETE",
		Pattern:     "/regimens/{regimen_id}/reminders/{reminder_id}",
		HandlerFunc: handlers.DeleteReminderByID,
	},
	"account.invitation": {
		Name:        "Get Invitation By ID",
		Method:      "GET",
		Pattern:     "/invitation/{invitation_id}",
		HandlerFunc: handlers.GetInvitationByID,
	},
	"account.invitations": {
		Name:        "Get All Invitations",
		Method:      "GET",
		Pattern:     "/invitation",
		HandlerFunc: handlers.GetInvitations,
	},
	"account.invitation.create": {
		Name:        "Create Invitation",
		Method:      "POST",
		Pattern:     "/invitation",
		HandlerFunc: handlers.PostInvitation,
	},
	"account.invitation.delete": {
		Name:        "Delete Invitation By ID",
		Method:      "DELETE",
		Pattern:     "/invitation/{invitation_id}",
		HandlerFunc: handlers.DeleteInvitation,
	},
	"account.invitation.accept": {
		Name:        "Get Invitation By ID",
		Method:      "GET",
		Pattern:     "/invitation/{invitation_id}/accept",
		HandlerFunc: handlers.AcceptInvitation,
	},
}
