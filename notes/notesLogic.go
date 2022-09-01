package notes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type LogicImpl struct {
	DB           *sql.DB
	Logger       *zap.Logger
	CalendarHost string
	Client       *http.Client
}

func (li *LogicImpl) GetAllNotes(ctx context.Context) ([]Note, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "getAllQuery",
		tracer.SpanType("db"),
		tracer.ServiceName("notes"),
		tracer.ResourceName("GetAllNotes"),
	)

	var err error
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	rows, err := li.DB.QueryContext(ctx, "SELECT * FROM notes ")
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	var notes []Note

	for rows.Next() {
		newNote := Note{}
		err = rows.Scan(&newNote.ID, &newNote.Description)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		notes = append(notes, newNote)
	}
	return notes, nil
}

func (li *LogicImpl) GetNote(ctx context.Context, id string) (Note, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "getQuery",
		tracer.SpanType("db"),
		tracer.ServiceName("notes"),
		tracer.ResourceName("GetNote"),
	)

	var err error
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	row := li.DB.QueryRowContext(ctx, "SELECT * FROM notes where id = ?", id)

	if err = row.Err(); err != nil {
		return Note{}, fmt.Errorf("query failed: %w", err)
	}
	note := Note{}
	err = row.Scan(&note.ID, &note.Description)
	if err != nil {
		return Note{}, fmt.Errorf("scan failed: %w", err)
	}

	return note, nil
}

func (li *LogicImpl) CreateNote(ctx context.Context, desc string, addDate bool) (Note, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "createQuery",
		tracer.SpanType("db"),
		tracer.ServiceName("notes"),
		tracer.ResourceName("createNote"),
	)
	var err error
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	if addDate {
		var date string
		date, err = li.getCalendarInfo(ctx)
		if err != nil {
			return Note{}, fmt.Errorf("getCalendarInfo failed: %w", err)
		}
		desc = desc + " with date " + date
	}

	newNote, err := li.DB.ExecContext(ctx, "INSERT INTO notes(description) VALUES (?) RETURNING id;")
	if err != nil {
		return Note{}, fmt.Errorf("execute query failed: %w", err)
	}

	newNoteId, err := newNote.LastInsertId()
	if err != nil {
		return Note{}, fmt.Errorf("lastInsertId failed: %w", err)
	}

	return Note{
		ID:          strconv.FormatInt(newNoteId, 10),
		Description: desc,
	}, nil
}

func (li *LogicImpl) UpdateNote(ctx context.Context, id string, newDescription string) (Note, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "updateQuery",
		tracer.SpanType("db"),
		tracer.ServiceName("notes"),
		tracer.ResourceName("updateNote"),
	)
	var err error
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	_, err = li.DB.ExecContext(ctx, "UPDATE notes set description = ? where id = ?")
	if err != nil {
		return Note{}, fmt.Errorf("execute query failed: %w", err)
	}

	return Note{id, newDescription}, nil
}

func (li *LogicImpl) DeleteNote(ctx context.Context, id string) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "deleteQuery",
		tracer.SpanType("db"),
		tracer.ServiceName("notes"),
		tracer.ResourceName("deleteNote"),
	)
	var err error
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	_, err = li.DB.ExecContext(ctx, "DELETE FROM notes where id = ?")
	if err != nil {
		return fmt.Errorf("query exec failed: %w", err)
	}

	return nil
}

func (li *LogicImpl) getCalendarInfo(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+li.CalendarHost+":9090/calendar", nil)
	if err != nil {
		return "", fmt.Errorf("request creation failed: %w", err)
	}

	resp, err := li.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}

	defer resp.Body.Close()
	m := json.NewDecoder(resp.Body)
	var date string
	err = m.Decode(&date)
	if err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}
	return date, nil
}
