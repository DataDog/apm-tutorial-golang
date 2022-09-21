package notes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type LogicImpl struct {
	DB           *sql.DB
	Logger       *zap.Logger
	CalendarHost string
	Client       *http.Client
}

func (li *LogicImpl) GetAllNotes(ctx context.Context) ([]Note, error) {
	rows, err := li.DB.QueryContext(ctx, "SELECT * FROM notes;")
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	notes := make([]Note, 0)

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
	row := li.DB.QueryRowContext(ctx, "SELECT * FROM notes WHERE id = ?;", id)

	if err := row.Err(); err != nil {
		return Note{}, fmt.Errorf("query failed: %w", err)
	}
	note := Note{}
	err := row.Scan(&note.ID, &note.Description)
	if err != nil {
		return Note{}, fmt.Errorf("scan failed: %w", err)
	}

	return note, nil
}

func (li *LogicImpl) CreateNote(ctx context.Context, desc string, addDate bool) (Note, error) {
	if addDate {
		var date string
		date, err := li.getCalendarInfo(ctx)
		if err != nil {
			return Note{}, fmt.Errorf("getCalendarInfo failed: %w", err)
		}
		desc = desc + " with date " + date
	}

	newNote, err := li.DB.ExecContext(ctx, "INSERT INTO notes(description) VALUES (?) RETURNING id;", desc)
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
	_, err := li.DB.ExecContext(ctx, "UPDATE notes set description = ? WHERE id = ?;", newDescription, id)
	if err != nil {
		return Note{}, fmt.Errorf("execute query failed: %w", err)
	}

	return Note{id, newDescription}, nil
}

func (li *LogicImpl) DeleteNote(ctx context.Context, id string) error {
	_, err := li.DB.ExecContext(ctx, "DELETE FROM notes WHERE id = ?;", id)
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
	var date string

	err = json.NewDecoder(resp.Body).Decode(&date)
	if err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}
	return date, nil
}
