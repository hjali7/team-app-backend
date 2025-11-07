package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the team-app-backend! API is running.")
	})

	// 2. تعریف روت "/health" - این برای Kubernetes حیاتی است
	// تا بتواند سلامت اپلیکیشن ما را بررسی کند
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// یک پاسخ JSON ساده که نشان می‌دهد سرویس فعال است
		fmt.Fprintln(w, `{"status": "ok", "service": "team-app-backend"}`)
	})

	port := ":8080"
	log.Printf("Starting backend server on http://localhost%s\n", port)

	// 3. راه‌اندازی سرور HTTP
	// اگر سرور نتواند اجرا شود، برنامه با خطا متوقف می‌شود
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}