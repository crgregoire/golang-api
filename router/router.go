package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tespo/buddha/auth"
	"github.com/tespo/buddha/handlers"
	"github.com/tespo/buddha/util"
)

//
// CreateRouter builds the endpoints
//
func CreateRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	for _, route := range LambdaRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(auth.LambdaRouterAuthenticationWrapper(route.Pattern, route.Method, util.SentryWrapper(route.HandlerFunc)))
	}

	for scope, route := range ImplicitRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(auth.ImplicitRouterAuthenticationWrapper(scope, util.SentryWrapper(route.HandlerFunc)))
	}

	for _, route := range ExplicitRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(auth.ExplicitRouterAuthenticationWrapper(route.Pattern, route.Method, util.SentryWrapper(route.HandlerFunc)))
	}

	for provider, routes := range VoiceCommandRoutes {
		for _, route := range routes {
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(auth.AuthenticateVoiceRequest(provider, util.SentryWrapper(route.HandlerFunc)))
		}
	}

	router.Methods("GET").Path("/").Name("Status Check").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200); w.Write([]byte("ok")) })

	return router
}
