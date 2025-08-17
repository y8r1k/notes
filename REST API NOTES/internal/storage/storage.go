package storage

import (
	"errors"
)

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrNoteExists   = errors.New("note exists")
)
