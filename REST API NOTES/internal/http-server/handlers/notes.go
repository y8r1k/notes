package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"notes/internal/http-server/dto"
	"notes/internal/logger/sl"
	"notes/internal/storage"
)

type App struct {
	storage storage.NoteStore
	log     *slog.Logger
}

func NewApp(storage storage.NoteStore, log *slog.Logger) *App {
	return &App{storage: storage, log: log}
}

func (a *App) HandleAllNoteGET(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http-server.handlers.HandleAllNoteGET"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Debug("incoming request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	// Get data from db
	notesDomain, err := a.storage.GETAllNotes()
	if err != nil {
		log.Error("db: get notes failed",
			sl.Err(err),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert data to response
	notesToResponse := make([]dto.NoteResponse, len(notesDomain))
	for i, note := range notesDomain {
		notesToResponse[i] = dto.DomainToResponse(note)
	}

	// Write response
	writeJSON(w, http.StatusOK, notesToResponse)
}

func (a *App) handleNoteGET(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http-server.handlers.handleNoteGET"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Debug("incoming request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	// Get id from url
	id, err := parseID(r)
	if err != nil {
		if errors.Is(err, ErrInvalidURLID) {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get data from db
	noteDomain, err := a.storage.GETNote(id)
	if err != nil {
		// Здесь по хорошему нужно добавить обработку ошибки ErrNoteNotFound
		log.Error("db: get note failed",
			slog.Int64("id", id),
			sl.Err(err),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert data to response
	noteToResponse := dto.DomainToResponse(noteDomain)

	// Write response
	writeJSON(w, http.StatusOK, noteToResponse)
}

func (a *App) handleNotePOST(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http-server.handlers.handleNotePOST"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Debug("incoming request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	// Decode data from body
	var newNote dto.CreateNotePOSTRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
		log.Error("failed to decode data",
			sl.Err(err),
		)
		http.Error(w, "Bad JSON data", http.StatusBadRequest)
		return
	}

	// Validate newNote
	if err := newNote.Validate(); err != nil {
		log.Warn("validation failed", sl.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create note-record in db
	noteDomain, err := a.storage.POSTNote(newNote.CreateRequestToDomain())
	if err != nil {
		// здесь можно отдельно обрабатывать ErrNoteExists
		log.Error("db: create note failed", sl.Err(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Info("note created",
		slog.Int64("id", noteDomain.ID),
	)

	// Convert data to response
	noteToResponse := dto.DomainToResponse(noteDomain)

	// Write response
	writeJSON(w, http.StatusCreated, noteToResponse)

}

func (a *App) handleNotePUT(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http-server.handlers.handleNotePUT"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Debug("incoming request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	// Get id from url
	id, err := parseID(r)
	if err != nil {
		if errors.Is(err, ErrInvalidURLID) {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Decode data from body
	var changingNote dto.ChangeNotePUTRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&changingNote); err != nil {
		log.Error("failed to decode data",
			sl.Err(err),
		)
		http.Error(w, "Bad JSON data", http.StatusBadRequest)
		return
	}

	if err := changingNote.Validate(); err != nil {
		log.Warn("validation failed", sl.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Change data in db
	noteDomain, err := a.storage.PUTNote(changingNote.ChangeRequestToDomain(id))
	if err != nil {
		//  здесь можно отдельно обрабатывать ErrNoteNotFound и ErrNoteExists
		log.Error("db: change note failed",
			slog.Int64("id", id),
			sl.Err(err),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Info("note changed",
		slog.Int64("id", noteDomain.ID),
	)

	// Convert data to response
	noteToResponse := dto.DomainToResponse(noteDomain)

	// Write response
	writeJSON(w, http.StatusOK, noteToResponse)
}

func (a *App) handleNoteDELETE(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http-server.handlers.handleNoteDELETE"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Debug("incoming request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	// Get id from url
	id, err := parseID(r)
	if err != nil {
		if errors.Is(err, ErrInvalidURLID) {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Delete data from db
	err = a.storage.DELETENote(id)
	if err != nil {
		// здесь можно отдельно обрабатывать ErrNoteNotFound
		log.Error("db: delete note failed",
			slog.Int64("id", id),
			sl.Err(err),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Info("note deleted",
		slog.Int64("id", id),
	)

	w.WriteHeader(http.StatusNoContent)
}

func (a *App) HandleNoteRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.handleNoteGET(w, r)
	case http.MethodPost:
		a.handleNotePOST(w, r)
	case http.MethodPut:
		a.handleNotePUT(w, r)
	case http.MethodDelete:
		a.handleNoteDELETE(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
