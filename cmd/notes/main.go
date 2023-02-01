package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"
	//sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	//chitrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi"
	//httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	//"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/datadog/apm_tutorial_golang/notes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//tracer.Start()
	//defer tracer.Stop()

	logger, _ := zap.NewDevelopment()
	logger.Debug("Starting from port 8080")

	db := setupDB(logger)
	defer db.Close()

	client := http.DefaultClient
	// Creates span with resource name equal to http Method and path
	//client = httptrace.WrapClient(client, httptrace.RTWithResourceNamer(func(req *http.Request) string {
	//	return fmt.Sprintf("%s %s", req.Method, req.URL.Path)
	//}))

	host, found := os.LookupEnv("CALENDAR_HOST")
	if !found || host == "" {
		host = "localhost"
	}

	logic := &notes.LogicImpl{
		DB:           db,
		Logger:       logger,
		Client:       client,
		CalendarHost: host,
	}

	nr := notes.Router{
		Logger: logger,
		Logic:  logic,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//r.Use(chitrace.Middleware(chitrace.WithServiceName("notes")))
	r.Mount("/", nr.Register())

	log.Fatal(http.ListenAndServe(":8080", r))
}

func setupDB(logger *zap.Logger) *sql.DB {
	//sqltrace.Register("sqlite3", &sqlite3.SQLiteDriver{}, sqltrace.WithServiceName("db"))
	//db, err := sqltrace.Open("sqlite3", "file::memory:?cache=shared")
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		logger.Fatal("error setting up database", zap.Error(err))
	}

	sts := ` DROP TABLE IF EXISTS notes;
			CREATE TABLE notes(id INTEGER PRIMARY KEY, description TEXT);`
	_, err = db.Exec(sts)
	if err != nil {
		logger.Fatal("error creating schema", zap.Error(err))
	}
	return db
}
