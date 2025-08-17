package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notes/internal/http-server/dto"
	"notes/internal/storage/postgres"
	"os"
	"strconv"
	"strings"
)

type App struct {
	Storage *postgres.Storage
}

func (a *App) HandleAllNoteGET(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http-server.handlers.HandleAllNoteGET"

	// Getting data from db
	notesDomain, err := a.Storage.GETAllNotes()
	if err != nil {
		// Здесь по хорошему нужно добавить обработку
		// Пока сделан общий вариант
		fmt.Fprintf(os.Stderr, "%s: database error: %v\n", op, err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Converting data to response
	notesToResponse := make([]dto.NoteResponse, len(notesDomain))
	for i, note := range notesDomain {
		notesToResponse[i] = dto.DomainToResponse(note)
	}

	// Marshaling data
	b, err := json.Marshal(notesToResponse)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to marshaling %v\n", op, err)
		return
	}
	w.Write(b)
}

func (a *App) handleNoteGET(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http-server.handlers.handleNoteGET"

	// Getting id from url
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 3 || parts[1] != "notes" {
		fmt.Fprintf(os.Stderr, "%s: url params error %q:", op, r.URL.Path)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: parse error %q: %v\n", op, idStr, err)
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	// Getting data from db
	noteDomain, err := a.Storage.GETNote(id)
	if err != nil {
		// Здесь по хорошему нужно добавить обработку
		// Ошибка запроса: такого id не сущесвует
		// Пока сделан общий вариант
		fmt.Fprintf(os.Stderr, "%s: database error: %v\n", op, err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Converting data to response
	noteToResponse := dto.DomainToResponse(noteDomain)

	// Marshaling data
	b, err := json.Marshal(noteToResponse)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to marshaling %v\n", op, err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (a *App) handleNotePOST(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http-server.handlers.HandleAllNoteGET"

	// Decode data from body
	defer r.Body.Close()

	var newNote dto.CreateNotePOSTRequest

	if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to decode data: %v\n", op, err)
		http.Error(w, "Bad JSON data", http.StatusBadRequest)
		return
	}

	// Getting data from db
	noteDomain, err := a.Storage.POSTNote(newNote.CreateRequestToDomain())
	if err != nil {
		// Здесь по хорошему нужно добавить обработку
		// Пока сделан общий вариант
		fmt.Fprintf(os.Stderr, "%s: database error: %v\n", op, err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Converting data to response
	noteToResponse := dto.DomainToResponse(noteDomain)

	// Marshaling data
	resultNoteByte, err := json.Marshal(noteToResponse)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to adding note: %v\n", op, err)
		http.Error(w, "Bad JSON data", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(resultNoteByte)
}

func (a *App) handleNotePUT(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http-server.handlers.handleNotePUT"

	// Getting id from url
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 3 || parts[1] != "notes" {
		fmt.Fprintf(os.Stderr, "%s: url params error %q:", op, r.URL.Path)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: parse error %q: %v\n", op, idStr, err)
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	// Decode data from body
	defer r.Body.Close()

	var changingNote dto.ChangeNotePUTRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&changingNote); err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to decode data: %v\n", op, err)
		http.Error(w, "Bad JSON data", http.StatusBadRequest)
		return
	}

	// Getting data from db
	noteDomain, err := a.Storage.PUTNote(changingNote.ChangeRequestToDomain(id))
	if err != nil {
		// Здесь по хорошему нужно добавить обработку
		// Пока сделан общий вариант
		fmt.Fprintf(os.Stderr, "%s: database error: %v\n", op, err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Converting data to response
	noteToResponse := dto.DomainToResponse(noteDomain)

	// Marshaling data
	resultNoteByte, err := json.Marshal(noteToResponse)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to marshaling data: %v\n", op, err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.Write(resultNoteByte)
}

func (a *App) handleNoteDELETE(w http.ResponseWriter, r *http.Request) {
	const op = "internal.http-server.handlers.handleNoteGET"

	// Getting id from url
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 3 || parts[1] != "notes" {
		fmt.Fprintf(os.Stderr, "%s: url params error %q:", op, r.URL.Path)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: parse error %q: %v\n", op, idStr, err)
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	// Getting data from db
	err = a.Storage.DELETENote(id)
	if err != nil {
		// Здесь по хорошему нужно добавить обработку
		// Ошибка запроса: такого id не сущесвует
		// Пока сделан общий вариант
		fmt.Fprintf(os.Stderr, "%s: database error: %v\n", op, err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	message := "Successful deletion"
	result := []byte(message)

	// Можно было добавить статус, но оставим так
	// w.WriteHeader(http.StatusNoContent)
	w.Write(result)
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
