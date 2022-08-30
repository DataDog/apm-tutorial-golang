package notes

import (
	"context"
	"fmt"

	"github.com/mattn/go-sqlite3"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func getNote(id string, parentContext context.Context) []Note {
	sqltrace.Register("sqlite3", &sqlite3.SQLiteDriver{}, sqltrace.WithServiceName("notes"))
	span, ctx := tracer.StartSpanFromContext(parentContext, "getQuery",
		tracer.SpanType("db"),
		tracer.ServiceName("notes"),
		tracer.ResourceName("getNoteByID"),
	)

	db, err := sqltrace.Open("sqlite3", "file::memory:?cache=shared")
	checkErr(err)
	defer db.Close()

	if id == "all" {
		rows, err := db.QueryContext(ctx, "SELECT * FROM notes ")
		checkErr(err)
		defer rows.Close()

		notes := make([]Note, 0)

		for rows.Next() {
			newNote := Note{}
			err = rows.Scan(&newNote.ID, &newNote.Description)
			checkErr(err)

			notes = append(notes, newNote)
		}
		span.Finish(tracer.WithError(err))
		return notes

	} else {
		rows, err := db.QueryContext(ctx, "SELECT * FROM notes where id = ?", id)
		checkErr(err)
		defer rows.Close()

		notes := make([]Note, 0)

		for rows.Next() {
			newNote := Note{}
			err = rows.Scan(&newNote.ID, &newNote.Description)
			checkErr(err)

			notes = append(notes, newNote)
		}

		span.Finish(tracer.WithError(err))
		return notes

	}
}

func addNote(desc string, parentContext context.Context) string {
	sqltrace.Register("sqlite3", &sqlite3.SQLiteDriver{}, sqltrace.WithServiceName("notes"))
	span, ctx := tracer.StartSpanFromContext(parentContext, "createQuery",
		tracer.SpanType("db"),
		tracer.ServiceName("notes"),
		tracer.ResourceName("createNote"),
	)

	db, err := sqltrace.Open("sqlite3", "file::memory:?cache=shared")
	checkErr(err)
	defer db.Close()

	stmt, err := db.PrepareContext(ctx, "INSERT INTO notes(description) VALUES (?) RETURNING id;")
	checkErr(err)

	newNote, err := stmt.ExecContext(ctx, desc)
	checkErr(err)
	defer stmt.Close()

	newNoteId, err := newNote.LastInsertId()
	checkErr(err)

	span.Finish(tracer.WithError(err))
	return fmt.Sprint(newNoteId)
}

func updateNote(id string, desc string, parentContext context.Context) Note {
	sqltrace.Register("sqlite3", &sqlite3.SQLiteDriver{}, sqltrace.WithServiceName("notes"))

	span, ctx := tracer.StartSpanFromContext(parentContext, "updateQuery",
		tracer.SpanType("db"),
		tracer.ServiceName("notes"),
		tracer.ResourceName("updateNote"),
	)
	db, err := sqltrace.Open("sqlite3", "file::memory:?cache=shared")
	checkErr(err)
	defer db.Close()

	stmt, err := db.PrepareContext(ctx, "UPDATE notes set description = ? where id = ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, desc, id)
	checkErr(err)

	span.Finish(tracer.WithError(err))
	return Note{id, desc}
}

func deleteNote(id string, parentContext context.Context) string {
	sqltrace.Register("sqlite3", &sqlite3.SQLiteDriver{}, sqltrace.WithServiceName("notes"))

	span, ctx := tracer.StartSpanFromContext(parentContext, "deleteQuery",
		tracer.SpanType("db"),
		tracer.ServiceName("notes"),
		tracer.ResourceName("deleteNote"),
	)
	db, err := sqltrace.Open("sqlite3", "file::memory:?cache=shared")
	checkErr(err)
	defer db.Close()

	stmt, err := db.PrepareContext(ctx, "DELETE FROM notes where id = ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	checkErr(err)

	span.Finish(tracer.WithError(err))
	return "Deleted"
}
