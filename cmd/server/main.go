package main

import (
	"log"
	"net/http"
    "os"

	"team-app-backend/internal/database"
	"team-app-backend/internal/server"
)

func main() {
	// ۱. راه‌اندازی وابستگی‌ها (دیتابیس)
	dbPool, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer dbPool.Close() // بستن pool هنگام خروج

	// ۲. ساختن سرور (تزریق وابستگی)
	srv := server.NewServer(dbPool)

	// ۳. راه‌اندازی سرور
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting backend server on http://localhost:%s\n", port)
	
    // srv (سرور ما) حالا http.Handler است
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatal(err)
	}
}