package main

import (
	"fmt"
	"net/http"
	"notes/internal/config"
	"notes/internal/http-server/handlers"
	"notes/internal/http-server/middlewares/cors"
	"notes/internal/storage/postgres"
	"os"
)

func main() {
	// Конфиг
	cfg := config.MustLoad()

	// Подключение Postgres
	storage, err := postgres.New(cfg.DBConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error of connecting DB: %v\n", err)
		os.Exit(1)
	}

	defer storage.Close()

	// Инициализация сервера
	app := handlers.App{Storage: storage}

	mux := http.NewServeMux()

	// Handlers
	mux.HandleFunc("/notes/", app.HandleNoteRequest)
	mux.HandleFunc("/notes", app.HandleAllNoteGET)

	// Middlewares
	handler := cors.WithCORS(mux)

	// Запуск сервера
	fmt.Println("Сервер запускается...")
	err = http.ListenAndServe(cfg.HTTPServer.Addres, handler)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
