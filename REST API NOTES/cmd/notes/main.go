package main

import (
	"log/slog"
	"net/http"
	"notes/internal/config"
	"notes/internal/http-server/handlers"
	"notes/internal/http-server/middlewares/cors"
	"notes/internal/logger/sl"
	"notes/internal/storage/postgres"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := sl.SetupLogger(cfg.Env)

	log.Info("starting notes service", slog.String("env", cfg.Env))

	storage, err := postgres.New(cfg.DBConfig)
	if err != nil {
		log.Error("failed to connect DB", sl.Err(err))
		os.Exit(1)
	}
	defer storage.Close()

	log.Info("database connected")

	app := handlers.App{Storage: storage, Log: log}
	mux := http.NewServeMux()
	mux.HandleFunc("/notes/", app.HandleNoteRequest)
	mux.HandleFunc("/notes", app.HandleAllNoteGET)

	handler := cors.WithCORS(mux)

	log.Info("launching server", slog.String("address", cfg.HTTPServer.Addres))
	srv := &http.Server{
		Addr:    cfg.HTTPServer.Addres,
		Handler: handler,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Error("server stopped")
}
