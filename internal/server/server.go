package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

// Server объединяет логгер и http-сервер.
type Server struct {
	logger *log.Logger
	server *http.Server
}

// New создаёт http-роутер, регистрирует в нём хендлеры и возвращает настроенный экземпляр сервера.
func New(logger *log.Logger) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/upload", handlers.Upload)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		logger: logger,
		server: srv,
	}
}

// Start запускает сервер и блокируется до его остановки или ошибки.
func (s *Server) Start() error {
	s.logger.Printf("starting server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}
