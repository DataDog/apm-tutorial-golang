package notes

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/go-chi/chi"
)

type Logic interface {
	GetAllNotes(ctx context.Context) ([]Note, error)
	GetNote(ctx context.Context, id string) (Note, error)
	CreateNote(ctx context.Context, description string, addDate bool) (Note, error)
	UpdateNote(ctx context.Context, id string, newDescription string) (Note, error)
	DeleteNote(ctx context.Context, id string) error
}

type Router struct {
	Logger *zap.Logger
	Logic  Logic
}

func (nr *Router) Register() chi.Router {
	r := chi.NewRouter()
	r.Get("/notes", makeSpanMiddleware("GetAllNotes", nr.GetAllNotes))               // GET /notes
	r.Post("/notes", makeSpanMiddleware("CreateNote", nr.CreateNote))                // POST /notes
	r.Get("/notes/{noteID}", makeSpanMiddleware("GetNote", nr.GetNoteByID))          // GET /notes/123
	r.Put("/notes/{noteID}", makeSpanMiddleware("UpdateNote", nr.UpdateNoteByID))    // PUT /notes/123
	r.Delete("/notes/{noteID}", makeSpanMiddleware("DeleteNote", nr.DeleteNoteByID)) // DELETE /notes/123

	r.Post("/notes/quit", func(rw http.ResponseWriter, r *http.Request) {
		time.AfterFunc(1*time.Second, func() { os.Exit(0) })
		rw.Write([]byte("Goodbye\n"))
	}) //Quits program

	return r
}

func makeSpanMiddleware(name string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span, ctx := tracer.StartSpanFromContext(r.Context(), name)
		r = r.WithContext(ctx)
		defer span.Finish()
		h.ServeHTTP(w, r)
	}
}

func reportError(err error, category string, w http.ResponseWriter) {
	msg := struct {
		Category string `json:"category"`
		Message  string `json:"message"`
	}{
		Category: category,
		Message:  err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(msg)
}

func reportInputError(message string, w http.ResponseWriter) {
	msg := struct {
		Category string `json:"category"`
		Message  string `json:"message"`
	}{
		Category: message,
		Message:  "invalid input",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(msg)
}

func (nr *Router) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	notes, err := nr.Logic.GetAllNotes(ctx)
	if err != nil {
		reportError(err, "GetAllNotes", w)
		return
	}

	doLongRunningProcess(ctx)
	anotherProcess(ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		reportError(err, "GetAllNotes-Encode", w)
		return
	}
}

func (nr *Router) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "noteID")
	if strings.TrimSpace(id) == "" {
		reportInputError("noteID not specified", w)
		return
	}
	note, err := nr.Logic.GetNote(ctx, id)
	if err != nil {
		reportError(err, "GetNoteByID", w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		reportError(err, "GetNotes-Encode", w)
		return
	}
}

func (nr *Router) CreateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	desc := r.URL.Query().Get("desc")
	addDate := false

	if strings.EqualFold(r.URL.Query().Get("add_date"), "y") {
		addDate = true
	}

	note, err := nr.Logic.CreateNote(ctx, desc, addDate)
	if err != nil {
		reportError(err, "CreateNote", w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		reportError(err, "CreateNote-Encode", w)
		return
	}
}

func (nr *Router) UpdateNoteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "noteID")
	desc := r.URL.Query().Get("desc")
	note, err := nr.Logic.UpdateNote(ctx, id, desc)
	if err != nil {
		reportError(err, "UpdateNoteByID", w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		reportError(err, "UpdateNote-Encode", w)
		return
	}
}

func (nr *Router) DeleteNoteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "noteID")
	err := nr.Logic.DeleteNote(ctx, id)
	if err != nil {
		reportError(err, "DeleteNoteByID", w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode("Deleted")
	if err != nil {
		reportError(err, "DeleteNote-Encode", w)
		return
	}
}
