package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"notes/internal/config"
	"notes/internal/domain"
	"notes/internal/storage"
	"os"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Close() {
	fmt.Println("Closing database connecting")
	s.db.Close()
}

func New(cfg config.DBConfig) (*Storage, error) {
	const op = "storage.postgres.New"
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s",
		cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	// Open connecting
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: error to connecting to database: %v", op, err)
	}

	// Check connecting
	err = db.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: ping database is failed: %v\n", op, err)
	}

	// Prepare query to create table
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS notes(
			id serial PRIMARY KEY,
			title varchar(20) NOT NULL,
			description text NOT NULL,
			created timestamp NOT NULL DEFAULT current_timestamp,
			changed timestamp
		)
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare query to create table 'notes': %v", op, err)
	}

	// Execute query
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to exec query to create table 'notes': %v", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GETAllNotes() ([]domain.Note, error) {
	const op = "storage.postgres.GETAllNotes"

	// Prepare query to getting data
	rows, err := s.db.Query(`
		SELECT * FROM notes
	`)

	if err != nil {
		return []domain.Note{}, fmt.Errorf("%s: failed to prepare query to selecting: %v", op, err)
	}

	// Getting data
	notes := []domain.Note{}
	for rows.Next() {
		var note NoteDB
		err := rows.Scan(&note.ID, &note.Title, &note.Description, &note.Created, &note.Changed)
		if err != nil {
			return []domain.Note{}, fmt.Errorf("%s: failed to executing statement: %v", op, err)
		}
		notes = append(notes, note.DBToDomain())
	}

	return notes, nil
}

/*
Пока нет реализации на фронте: получение отдельной записи
*/
func (s *Storage) GETNote(id int64) (domain.Note, error) {
	const op = "storage.postgres.GETNote"

	// Prepare query to getting data
	stmt, err := s.db.Prepare(`
		SELECT *
			FROM notes
			WHERE id = $1
	`)

	if err != nil {
		return domain.Note{}, fmt.Errorf("%s: failed to getting data from db: %v", op, err)
	}

	defer stmt.Close()

	// Getting data
	var note NoteDB
	err = stmt.QueryRow(id).Scan(&note.ID, &note.Title, &note.Description, &note.Created, &note.Changed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Note{}, storage.ErrNoteNotFound
		}
		return domain.Note{}, fmt.Errorf("%s: failed to executing statement: %w", op, err)
	}

	return note.DBToDomain(), nil
}

/*
Точечно не обработана ошибка, существования в таблице записи с данным первичным ключем
*/
func (s *Storage) POSTNote(note domain.Note) (domain.Note, error) {
	const op = "storage.postgres.POSTNote"

	// Prepare query to adding data
	stmt, err := s.db.Prepare(`
		INSERT INTO notes(title, description) VALUES($1, $2)
		RETURNING id
	`)

	if err != nil {
		return domain.Note{}, fmt.Errorf("%s: failed to prepare query to inserting: %w", op, err)
	}

	defer stmt.Close()

	// Adding data and getting ID of added data
	var id int64
	err = stmt.QueryRow(note.Title, note.Description).Scan(&id)
	if err != nil {
		return domain.Note{}, fmt.Errorf("%s: failed to executing statement: %v", op, err)
	}

	return s.GETNote(id)
}

func (s *Storage) DELETENote(id int64) error {
	const op = "storage.postgres.DELETENote"

	// Prepare query to deleting data
	stmt, err := s.db.Prepare(`
		DELETE FROM notes
			WHERE id = $1
	`)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare query to deleting: %v", op, err)
	}

	defer stmt.Close()

	// Deleting data
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: failed to executing statement: %w", op, err)
	}

	return nil
}

func (s *Storage) PUTNote(note domain.Note) (domain.Note, error) {
	const op = "storage.postgres.PUTNote"

	// Prepare query to getting data
	stmt, err := s.db.Prepare(`
		UPDATE notes
			SET title = $1,
				description = $2,
				changed = NOW()
			WHERE id = $3
	`)
	if err != nil {
		return domain.Note{}, fmt.Errorf("%s: : failed to prepare query to selecting: %v", op, err)
	}

	// Changing data
	_, err = stmt.Exec(note.Title, note.Description, note.ID)
	if err != nil {
		return domain.Note{}, fmt.Errorf("%s: failed to executing statement: %v", op, err)
	}

	return s.GETNote(note.ID)
}
