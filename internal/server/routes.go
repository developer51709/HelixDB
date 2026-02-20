package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/developer51709/helixdb/internal/storage"
)

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/", s.handleRoot)
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/collections", s.handleListCollections)
	s.mux.HandleFunc("/collections/", s.handleCollections)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name":    "HelixDB",
		"version": "0.1.0",
		"status":  "running",
		"uptime":  time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "healthy",
	})
}

func (s *Server) handleListCollections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	collections := s.engine.ListCollections()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"collections": collections,
	})
}

func (s *Server) handleCollections(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/collections/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "collection name required"})
		return
	}

	collectionName := parts[0]

	if len(parts) >= 2 && parts[1] == "query" {
		s.handleQuery(w, r, collectionName)
		return
	}

	if len(parts) == 2 && parts[1] != "" {
		docID := parts[1]
		switch r.Method {
		case http.MethodGet:
			s.handleGetDocument(w, r, collectionName, docID)
		case http.MethodDelete:
			s.handleDeleteDocument(w, r, collectionName, docID)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
		return
	}

	switch r.Method {
	case http.MethodPost:
		s.handleCreateDocument(w, r, collectionName)
	case http.MethodGet:
		s.handleListDocuments(w, r, collectionName)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) handleCreateDocument(w http.ResponseWriter, r *http.Request, collection string) {
	var body struct {
		ID   string                 `json:"id"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	if body.Data == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "data field is required"})
		return
	}

	id := body.ID
	if id == "" {
		id = generateID()
	}

	doc, err := s.engine.InsertDocument(collection, id, body.Data)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, doc)
}

func (s *Server) handleGetDocument(w http.ResponseWriter, r *http.Request, collection, id string) {
	doc, exists := s.engine.GetDocument(collection, id)
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "document not found"})
		return
	}
	writeJSON(w, http.StatusOK, doc)
}

func (s *Server) handleDeleteDocument(w http.ResponseWriter, r *http.Request, collection, id string) {
	if deleted := s.engine.DeleteDocument(collection, id); !deleted {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "document not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) handleListDocuments(w http.ResponseWriter, r *http.Request, collection string) {
	docs := s.engine.QueryDocuments(collection, nil, 0)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"collection": collection,
		"count":      len(docs),
		"documents":  docs,
	})
}

func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request, collection string) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var body struct {
		Filter map[string]interface{} `json:"filter"`
		Limit  int                    `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	docs := s.engine.QueryDocuments(collection, body.Filter, body.Limit)

	result := make([]storage.Document, 0, len(docs))
	for _, d := range docs {
		result = append(result, *d)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"collection": collection,
		"count":      len(result),
		"documents":  result,
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
