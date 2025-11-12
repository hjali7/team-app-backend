package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() (*pgxpool.Pool, error) {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	if host == "" {
		host = "localhost" 
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
    if dbname == "" {
        dbname = "team_app_db"
    }
    if password == "" {
        // password = "yourlocalpassword" 
    }


	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Database connection pool created successfully.")
	return pool, nil
}