package main

import (
	"log"
	"net/http"

	"github.com/thecemakin/hr-project/internal/platform/config"
	"github.com/thecemakin/hr-project/internal/platform/db"
	server "github.com/thecemakin/hr-project/internal/platform/http"
)

func main() {
	log.Println("HR Backend API starting...")

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Setup database connection (ignoring currently to just test boot)
	database, err := db.Connect(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v\nContinuing without db to allow boot test.", err)
	} else {
		log.Println("Database connection established:", database.Name())
	}

	// 3. Setup HTTP server and routing
	srv := server.NewServer()

	log.Printf("Listening and serving HTTP on :%s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, srv.Router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
