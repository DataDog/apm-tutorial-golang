package notes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func NoteRouter() chi.Router {

	r := chi.NewRouter()
	r.Get("/notes", getAllNotes)                // GET /notes
	r.Post("/notes", createNote)                // POST /notes
	r.Get("/notes/{noteID}", getNoteByID)       // GET /articles/123
	r.Put("/notes/{noteID}", updateNoteByID)    // PUT /articles/123
	r.Delete("/notes/{noteID}", deleteNoteByID) // DELETE /articles/123

	return r
}

func getAllNotes(w http.ResponseWriter, r *http.Request) {

	span, ctx := tracer.StartSpanFromContext(r.Context(), "getNote")
	defer span.Finish()

	response, _ := json.Marshal(getNote("all", ctx))

	doLongRunningProcess(span)
	anotherProcess(span)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func getNoteByID(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartSpanFromContext(r.Context(), "postNote")
	defer span.Finish()

	response, _ := json.Marshal(getNote(chi.URLParam(r, "noteID"), ctx))
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func createNote(w http.ResponseWriter, r *http.Request) {

	span, ctx := tracer.StartSpanFromContext(r.Context(), "postNote")
	defer span.Finish()

	//set ID with increasing levels of ID
	desc := r.URL.Query().Get("desc")

	if r.URL.Query().Get("add_date") != "" && strings.EqualFold(r.URL.Query().Get("add_date"), "y") {
		host, found := os.LookupEnv("CALENDAR_HOST")
		if !found || host == "" {
			host = "localhost"
		}

		req, err := http.NewRequest("GET", "http://"+host+":9090/calendar", nil)
		checkErr(err)

		req = req.WithContext(ctx)
		// Inject the span Context in the Request headers
		err = tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(req.Header))
		checkErr(err)

		resp, err := http.DefaultClient.Do(req)
		checkErr(err)

		body, err := ioutil.ReadAll(resp.Body)
		checkErr(err)

		desc = desc + " with date " + string(body)
	}

	//create note with desc and insert into table
	testNote := Note{
		ID:          addNote(desc, ctx),
		Description: desc,
	}

	response, _ := json.Marshal(testNote)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func updateNoteByID(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartSpanFromContext(r.Context(), "postNote")
	defer span.Finish()

	response, _ := json.Marshal(updateNote(chi.URLParam(r, "noteID"), r.URL.Query().Get("desc"), ctx))
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func deleteNoteByID(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartSpanFromContext(r.Context(), "postNote")
	defer span.Finish()

	response, _ := json.Marshal(deleteNote(chi.URLParam(r, "noteID"), ctx))
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
