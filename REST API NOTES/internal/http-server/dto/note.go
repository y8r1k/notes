package dto

import (
	"notes/internal/domain"
	"time"
)

type NoteResponse struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Created     time.Time  `json:"created"`
	Changed     *time.Time `json:"changed,omitempty"`
}

type CreateNotePOSTRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ChangeNotePUTRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (note CreateNotePOSTRequest) CreateRequestToDomain() domain.Note {
	return domain.Note{
		Title:       note.Title,
		Description: note.Description,
	}
}

func (note ChangeNotePUTRequest) ChangeRequestToDomain(id int64) domain.Note {
	return domain.Note{
		ID:          id,
		Title:       note.Title,
		Description: note.Description,
	}
}

func DomainToResponse(n domain.Note) NoteResponse {
	return NoteResponse{
		ID:          n.ID,
		Title:       n.Title,
		Description: n.Description,
		Created:     n.Created,
		Changed:     n.Changed,
	}
}
