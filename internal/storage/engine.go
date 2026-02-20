package storage

import (
        "crypto/sha256"
        "encoding/json"
        "fmt"
        "os"
        "path/filepath"
        "sync"
        "time"
)

type Document struct {
        ID        string                 `json:"id"`
        Data      map[string]interface{} `json:"data"`
        CreatedAt time.Time              `json:"createdAt"`
        UpdatedAt time.Time              `json:"updatedAt"`
        Checksum  string                 `json:"checksum"`
}

type Collection struct {
        Name      string               `json:"name"`
        Documents map[string]*Document `json:"documents"`
        mu        sync.RWMutex
}

type Engine struct {
        dataFile string
        walDir   string
        collections map[string]*Collection
        mu          sync.RWMutex
        wal         *WAL
}

func NewEngine(dataFile, walDir string) (*Engine, error) {
        if err := os.MkdirAll(filepath.Dir(dataFile), 0755); err != nil {
                return nil, fmt.Errorf("creating data directory: %w", err)
        }
        if err := os.MkdirAll(walDir, 0755); err != nil {
                return nil, fmt.Errorf("creating WAL directory: %w", err)
        }

        e := &Engine{
                dataFile:    dataFile,
                walDir:      walDir,
                collections: make(map[string]*Collection),
        }

        wal, err := NewWAL(walDir)
        if err != nil {
                return nil, fmt.Errorf("initializing WAL: %w", err)
        }
        e.wal = wal

        if err := e.loadFromDisk(); err != nil {
                fmt.Printf("[INFO] No existing data file found, starting fresh\n")
        }

        if err := e.recover(); err != nil {
                fmt.Printf("[WARN] Recovery encountered issues: %v\n", err)
        }

        return e, nil
}

func (e *Engine) GetCollection(name string) *Collection {
        e.mu.RLock()
        col, exists := e.collections[name]
        e.mu.RUnlock()
        if exists {
                return col
        }

        e.mu.Lock()
        defer e.mu.Unlock()
        if col, exists = e.collections[name]; exists {
                return col
        }
        col = &Collection{
                Name:      name,
                Documents: make(map[string]*Document),
        }
        e.collections[name] = col
        return col
}

func (e *Engine) ListCollections() []string {
        e.mu.RLock()
        defer e.mu.RUnlock()
        names := make([]string, 0, len(e.collections))
        for name := range e.collections {
                names = append(names, name)
        }
        return names
}

func (e *Engine) InsertDocument(collection string, id string, data map[string]interface{}) (*Document, error) {
        col := e.GetCollection(collection)

        now := time.Now().UTC()
        doc := &Document{
                ID:        id,
                Data:      data,
                CreatedAt: now,
                UpdatedAt: now,
        }
        doc.Checksum = computeChecksum(doc)

        entry := WALEntry{
                Operation:  "INSERT",
                Collection: collection,
                DocumentID: id,
                Data:       data,
                Timestamp:  now,
        }
        if err := e.wal.Write(entry); err != nil {
                return nil, fmt.Errorf("writing WAL: %w", err)
        }

        col.mu.Lock()
        col.Documents[id] = doc
        col.mu.Unlock()

        go func() { _ = e.saveToDisk() }()

        return doc, nil
}

func (e *Engine) GetDocument(collection, id string) (*Document, bool) {
        col := e.GetCollection(collection)
        col.mu.RLock()
        defer col.mu.RUnlock()
        doc, exists := col.Documents[id]
        return doc, exists
}

func (e *Engine) DeleteDocument(collection, id string) bool {
        col := e.GetCollection(collection)

        col.mu.Lock()
        _, exists := col.Documents[id]
        if exists {
                delete(col.Documents, id)
        }
        col.mu.Unlock()

        if exists {
                entry := WALEntry{
                        Operation:  "DELETE",
                        Collection: collection,
                        DocumentID: id,
                        Timestamp:  time.Now().UTC(),
                }
                _ = e.wal.Write(entry)
                go func() { _ = e.saveToDisk() }()
        }

        return exists
}

func (e *Engine) QueryDocuments(collection string, filter map[string]interface{}, limit int) []*Document {
        col := e.GetCollection(collection)
        col.mu.RLock()
        defer col.mu.RUnlock()

        var results []*Document
        for _, doc := range col.Documents {
                if matchesFilter(doc.Data, filter) {
                        results = append(results, doc)
                        if limit > 0 && len(results) >= limit {
                                break
                        }
                }
        }
        return results
}

func matchesFilter(data map[string]interface{}, filter map[string]interface{}) bool {
        if len(filter) == 0 {
                return true
        }
        for key, val := range filter {
                docVal, exists := data[key]
                if !exists {
                        return false
                }
                if fmt.Sprintf("%v", docVal) != fmt.Sprintf("%v", val) {
                        return false
                }
        }
        return true
}

func computeChecksum(doc *Document) string {
        data, _ := json.Marshal(doc.Data)
        h := sha256.Sum256(data)
        return fmt.Sprintf("%x", h[:8])
}

type diskData struct {
        Collections map[string]map[string]*Document `json:"collections"`
}

func (e *Engine) saveToDisk() error {
        e.mu.RLock()
        dd := diskData{Collections: make(map[string]map[string]*Document)}
        for name, col := range e.collections {
                col.mu.RLock()
                docs := make(map[string]*Document)
                for id, doc := range col.Documents {
                        docs[id] = doc
                }
                col.mu.RUnlock()
                dd.Collections[name] = docs
        }
        e.mu.RUnlock()

        data, err := json.MarshalIndent(dd, "", "  ")
        if err != nil {
                return err
        }

        tmpFile := e.dataFile + ".tmp"
        if err := os.WriteFile(tmpFile, data, 0644); err != nil {
                return err
        }
        return os.Rename(tmpFile, e.dataFile)
}

func (e *Engine) loadFromDisk() error {
        data, err := os.ReadFile(e.dataFile)
        if err != nil {
                return err
        }

        var dd diskData
        if err := json.Unmarshal(data, &dd); err != nil {
                return err
        }

        e.mu.Lock()
        defer e.mu.Unlock()
        for name, docs := range dd.Collections {
                col := &Collection{
                        Name:      name,
                        Documents: docs,
                }
                e.collections[name] = col
        }
        return nil
}

func (e *Engine) recover() error {
        entries, err := e.wal.ReadAll()
        if err != nil {
                return err
        }
        if len(entries) == 0 {
                return nil
        }

        fmt.Printf("[INFO] Replaying %d WAL entries\n", len(entries))
        for _, entry := range entries {
                col := e.GetCollection(entry.Collection)
                switch entry.Operation {
                case "INSERT":
                        doc := &Document{
                                ID:        entry.DocumentID,
                                Data:      entry.Data,
                                CreatedAt: entry.Timestamp,
                                UpdatedAt: entry.Timestamp,
                        }
                        doc.Checksum = computeChecksum(doc)
                        col.Documents[entry.DocumentID] = doc
                case "DELETE":
                        delete(col.Documents, entry.DocumentID)
                }
        }
        return nil
}

func (e *Engine) Close() error {
        if err := e.saveToDisk(); err != nil {
                return err
        }
        return e.wal.Close()
}
