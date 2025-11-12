package main

import (
	"log"
	"net/http"
    "os"

	"team-app-backend/internal/database"
	"team-app-backend/internal/server"
)

func main() {
	dbPool, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer dbPool.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	srv := server.NewServer(dbPool, jwtSecret)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting backend server on http://localhost:%s\n", port)
	
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatal(err)
	}
}