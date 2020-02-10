package integration

import (
	"flag"
	"net/http"
	"os"
	"testing"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/tespo/buddha/router"
)

func CreateRouter() *mux.Router {

	testRouter := mux.NewRouter().StrictSlash(true)

	for _, route := range router.ExplicitRoutes {
		testRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(testAccountRouteWrapper(route.HandlerFunc))
	}

	for _, route := range router.ImplicitRoutes {
		testRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(testAccountRouteWrapper(route.HandlerFunc))
	}

	return testRouter
}

func testAccountRouteWrapper(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		context.Set(request, "account_id", "d8e4c5dc-9767-41bd-b802-060e80d83867")
		context.Set(request, "user_id", "8c8aa229-3959-4a40-bbe6-67c2eeace5cb")
		context.Set(request, "scoped_fields", []string{"user.*"})
		context.Set(request, "scoped_fields", []string{"account.*"})
		next(writer, request)
	}
}

func TestMain(m *testing.M) {
	os.Setenv("TESTING_URL", "http://localhost:5555")
	os.Setenv("DB_HOST", "")
	os.Setenv("DB_NAME", "tespo_docker")
	os.Setenv("DB_USER", "root")

	flag.Parse()

	if !testing.Short() {
		srv := &http.Server{
			Handler: CreateRouter(),
			Addr:    ":5555",
		}
		go srv.ListenAndServe()
	}

	os.Exit(m.Run())
}
