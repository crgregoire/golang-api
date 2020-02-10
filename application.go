package main

import (
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/subosito/gotenv"
	"github.com/tespo/buddha/router"
)

func init() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}
	gotenv.Load("./config/" + env + ".env")
}

func main() {

	sentry.Init(sentry.ClientOptions{
		Dsn:         os.Getenv("SENTRY_DSN"),
		Environment: os.Getenv("GO_ENV"),
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	r := router.CreateRouter()

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + port,
	}

	log.Fatal(srv.ListenAndServe())
}
