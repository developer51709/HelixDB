# HelixDB

## Overview
HelixDB is a lightweight, local-first JSON database engine built in Go. It provides an HTTP/JSON API for document storage with write-ahead logging (WAL) for data integrity.

## Current State
- Core engine implemented with collections-based document storage
- HTTP API running on port 5000
- WAL-based crash recovery
- CORS enabled for cross-origin access

## Project Architecture

### Structure
```
cmd/helixdb/          - Main entry point, CLI parsing
internal/
  config/             - Configuration loading and schema
  server/             - HTTP server, routes, middleware
  storage/            - Storage engine, WAL
clients/              - Node.js and Python client libraries (stubs)
tests/                - Unit and integration tests (stubs)
api/                  - OpenAPI spec and example payloads
```

### Key Files
- `helixdb.config.json` - Server configuration (port, storage paths, security)
- `cmd/helixdb/main.go` - Application entry point
- `internal/storage/engine.go` - Core database engine with collections
- `internal/storage/wal.go` - Write-ahead log implementation
- `internal/server/routes.go` - HTTP API route handlers
- `internal/server/middleware.go` - CORS, logging, auth middleware

### API Endpoints
- `GET /` - Server info
- `GET /health` - Health check
- `GET /collections` - List all collections
- `POST /collections/:name` - Create document
- `GET /collections/:name` - List documents in collection
- `GET /collections/:name/:id` - Get document by ID
- `DELETE /collections/:name/:id` - Delete document
- `POST /collections/:name/query` - Query documents with filters

### Configuration
Server port defaults to 5000 (Replit compatible). Config is loaded from `helixdb.config.json`.

### Data Storage
- Data persisted to `./data/helix.db` (JSON format)
- WAL stored in `./data/wal/`

## Recent Changes
- 2026-02-20: Initial implementation of Go codebase from project skeleton
- Configured for Replit with port 5000
- Added WAL, storage engine, HTTP API, config loading
- Set up VM deployment target

## User Preferences
- None recorded yet
