package postgres

import (
	"database/sql"
	"notes/internal/domain"
	"time"
)

type NoteDB struct {
	ID          int64
	Title       string
	Description string
	Created     time.Time
	Changed     sql.NullTime
}

func (noteDB NoteDB) DBToDomain() domain.Note {
	var changed *time.Time
	if noteDB.Changed.Valid {
		changed = &noteDB.Changed.Time
	}
	return domain.Note{
		ID:          noteDB.ID,
		Title:       noteDB.Title,
		Description: noteDB.Description,
		Created:     noteDB.Created,
		Changed:     changed,
	}
}

// Преобразование из домена в бд случается только в 2 случаях:
// добавление новой записи
// изменение записи
// поэтому возможно время при изменении можно добавлять самому
func FromDomain(n domain.Note) NoteDB {
	return NoteDB{
		ID:          n.ID,
		Title:       n.Title,
		Description: n.Description,
		Created:     n.Created,
		Changed:     sql.NullTime{},
	}
}
