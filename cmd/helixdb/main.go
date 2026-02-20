package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/helixdb/helixdb/internal/config"
	"github.com/helixdb/helixdb/internal/server"
	"github.com/helixdb/helixdb/internal/storage"
)

func main() {
	command := "serve"
	configPath := ""

	args := os.Args[1:]
	for i, arg := range args {
		switch arg {
		case "serve", "backup", "recover":
			command = arg
		case "--config", "-c":
			if i+1 < len(args) {
				configPath = args[i+1]
			}
		}
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("[ERROR] Failed to load config: %v", err)
	}

	switch command {
	case "serve":
		runServe(cfg)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Usage: helixdb [serve|backup|recover] [--config path]")
		os.Exit(1)
	}
}

func runServe(cfg config.Config) {
	engine, err := storage.NewEngine(cfg.Storage.DataFile, cfg.Storage.WALDirectory)
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize storage engine: %v", err)
	}

	srv := server.New(engine, cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("[INFO] Shutting down HelixDB...")
		if err := srv.Shutdown(); err != nil {
			log.Printf("[ERROR] Shutdown error: %v", err)
		}
		os.Exit(0)
	}()

	if err := srv.Start(); err != nil {
		log.Fatalf("[ERROR] Server failed: %v", err)
	}
}
