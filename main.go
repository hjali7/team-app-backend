package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	if _, err := InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthCheckHandler)

	port := ":8080"
	log.Printf("Starting backend server on http://localhost%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the team-app-backenghjgkyhjghjfgghjfghjfghjfd! API is running v2.")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	serverStatus := "ok"
	dbStatus := "ok"
	dbError := ""

	pool := GetDBPool()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // ۲ ثانیه مهلت
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		dbStatus = "error"
		dbError = err.Error()
		log.Printf("Health check failed: Database ping error: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	response := map[string]string{
		"status":      serverStatus,
		"database":    dbStatus,
		"db_error":    dbError,
		"service":     "team-app-backend",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}