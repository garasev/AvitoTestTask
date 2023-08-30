package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"log/slog"

	_ "github.com/lib/pq"

	"github.com/garasev/AvitoTestTask/config"
	"github.com/garasev/AvitoTestTask/internal/handler"
	"github.com/garasev/AvitoTestTask/internal/repository/postgresql"
	"github.com/garasev/AvitoTestTask/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.GetConfig()

	log := setupLogger(*cfg)
	log.Info("App was started")
	log.Debug("Debug messages are enabled")

	db, err := ConnectDB(*cfg)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("DB was connected")
	repo := postgresql.NewPostgresRepo(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(*service, *log)

	router := Router(*handler)
	log.Info("Router was initialized")

	log.Info("Server is starting...")
	srv := &http.Server{
		Addr:    cfg.HTTPServer.Port,
		Handler: router,
	}
	err = srv.ListenAndServe()
	log.Error(err.Error())

}

func Router(handler handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Get("/{id}", handler.GetSlug)

	return r
}

func setupLogger(cfg config.Config) *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.Level(config.LogLevels[cfg.Logger.Level]),
		}),
	)
}

func ConnectDB(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DataBase.Username,
		cfg.DataBase.Password,
		cfg.DataBase.Host,
		cfg.DataBase.Port,
		cfg.DataBase.Database,
	)

	var db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
