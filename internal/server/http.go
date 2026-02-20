package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/helixdb/helixdb/internal/config"
	"github.com/helixdb/helixdb/internal/storage"
)

type Server struct {
	engine *storage.Engine
	config config.Config
	mux    *http.ServeMux
}

func New(engine *storage.Engine, cfg config.Config) *Server {
	s := &Server{
		engine: engine,
		config: cfg,
		mux:    http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	handler := s.withMiddleware(s.mux)

	log.Printf("[INFO] HelixDB server starting on %s", addr)
	log.Printf("[INFO] Data file: %s", s.config.Storage.DataFile)
	log.Printf("[INFO] WAL directory: %s", s.config.Storage.WALDirectory)

	return http.ListenAndServe(addr, handler)
}

func (s *Server) Shutdown() error {
	return s.engine.Close()
}
