package notes

import (
	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func (n Note) getID() string {
	return n.ID
}

func (n Note) getDescription() string {
	return n.Description
}

func (n *Note) setID(ID string) {
	n.ID = ID
}

func (n *Note) setDescription(description string) {
	n.Description = description
}

func NewNote(description string, id string) Note {
	return Note{ID: id, Description: description}
}
