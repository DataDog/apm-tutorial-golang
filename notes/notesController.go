package notes

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
	r.Get("/notes/{noteID}", makeSpanMiddleware("GetNode", nr.GetNoteByID))          // GET /articles/123
	r.Put("/notes/{noteID}", makeSpanMiddleware("UpdateNode", nr.UpdateNoteByID))    // PUT /articles/123
	r.Delete("/notes/{noteID}", makeSpanMiddleware("DeleteNote", nr.DeleteNoteByID)) // DELETE /articles/123

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

	response, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(response)
}

func reportInputError(message string, w http.ResponseWriter) {
	msg := struct {
		Category string `json:"category"`
		Message  string `json:"message"`
	}{
		Category: message,
		Message:  "invalid input",
	}

	response, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(response)
}

func (nr *Router) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	notes, err := nr.Logic.GetAllNotes(ctx)
	if err != nil {
		reportError(err, "GetAllNotes", w)
		return
	}
	response, err := json.Marshal(notes)
	if err != nil {
		reportError(err, "GetAllNotes-Marshal", w)
		return
	}

	doLongRunningProcess(ctx)
	anotherProcess(ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
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
	response, err := json.Marshal(note)
	if err != nil {
		reportError(err, "GetNoteByID-Marshal", w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
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

	response, err := json.Marshal(note)
	if err != nil {
		reportError(err, "CreateNote-Marshal", w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
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

	response, err := json.Marshal(note)
	if err != nil {
		reportError(err, "UpdateNoteByID-Marshal", w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (nr *Router) DeleteNoteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "noteID")
	err := nr.Logic.DeleteNote(ctx, id)
	if err != nil {
		reportError(err, "DeleteNoteByID", w)
		return
	}

	response, _ := json.Marshal("Deleted")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
