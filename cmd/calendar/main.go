package main

import (
	"log"
	"net/http"
	"time"

	"github.com/datadog/apm_tutorial_golang/calendar"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	chitrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {

	tracer.Start()
	defer tracer.Stop()

	log.Printf("Starting from port 9090")

	r := chi.NewRouter()

	r.Use(middleware.Timeout(2500 * time.Millisecond))
	r.Use(middleware.Logger)
	r.Use(chitrace.Middleware(chitrace.WithServiceName("calendar")))

	r.Get("/calendar", calendar.GetDate)

	log.Fatal(http.ListenAndServe(":9090", r))

}
