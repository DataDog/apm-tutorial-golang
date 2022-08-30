package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	notes "notes.com"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	chitrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"

	"github.com/mattn/go-sqlite3"
)

func main() {

	tracer.Start()
	defer tracer.Stop()

	log.Printf("Starting from port 8080")

	sqltrace.Register("sqlite3", &sqlite3.SQLiteDriver{}, sqltrace.WithServiceName("notes"))
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sts := ` DROP TABLE IF EXISTS notes;
			CREATE TABLE notes(id INTEGER PRIMARY KEY, description TEXT);`
	_, err = db.Exec(sts)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Timeout(2500 * time.Millisecond))
	r.Use(middleware.Logger)
	r.Use(chitrace.Middleware(chitrace.WithServiceName("notes")))

	r.Mount("/", notes.NoteRouter())

	log.Fatal(http.ListenAndServe(":8080", r))

}
