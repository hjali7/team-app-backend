package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

// از یک متغیر سراسری برای نگهداری "استخر اتصال" (Connection Pool) استفاده می‌کنیم
var (
	dbPool *pgxpool.Pool
	once   sync.Once // برای اطمینان از اینکه اتصال فقط یک بار برقرار می‌شود
)

// InitDB مسئول خواندن متغیرهای محیطی و ایجاد استخر اتصال است
func InitDB() (*pgxpool.Pool, error) {
	// این تابع فقط یک بار اجرا خواهد شد، حتی اگر چندین بار فراخوانی شود
	once.Do(func() {
		// خواندن پیکربندی از متغیرهای محیطی
		// این متغیرها هم در docker-compose.yml محلی و هم در K8s (ArgoCD) تنظیم شده‌اند
		host := os.Getenv("POSTGRES_HOST")
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbname := os.Getenv("POSTGRES_DB")
		port := os.Getenv("POSTGRES_PORT") // ممکن است در K8s تنظیم نشده باشد، پیش‌فرض می‌گذاریم

		if port == "" {
			port = "5432" // پورت پیش‌فرض Postgres
		}

		// ساختن رشته‌ی اتصال (Connection String)
		connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

		// ایجاد استخر اتصال جدید
		pool, err := pgxpool.New(context.Background(), connString)
		if err != nil {
			log.Fatalf("Unable to create connection pool: %v\n", err)
		}

		// پینگ کردن دیتابیس برای اطمینان از صحت اتصال
		if err := pool.Ping(context.Background()); err != nil {
			log.Fatalf("Unable to ping database: %v\n", err)
		}

		log.Println("Database connection pool created successfully.")
		dbPool = pool
	})

	return dbPool, nil
}

// GetDBPool به بقیه اپلیکیشن اجازه می‌دهد تا به استخر اتصال دسترسی داشته باشند
func GetDBPool() *pgxpool.Pool {
	if dbPool == nil {
		// اگر هنوز InitDB فراخوانی نشده، آن را فراخوانی می‌کنیم
		// این یک اقدام احتیاطی است
		InitDB()
	}
	return dbPool
}