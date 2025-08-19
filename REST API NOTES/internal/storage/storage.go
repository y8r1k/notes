package storage

import (
	"errors"
	"notes/internal/domain"
)

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrNoteExists   = errors.New("note exists")
)

type NoteStore interface {
	GETAllNotes() ([]domain.Note, error)
	GETNote(id int64) (domain.Note, error)
	POSTNote(note domain.Note) (domain.Note, error)
	PUTNote(note domain.Note) (domain.Note, error)
	DELETENote(id int64) error
}
