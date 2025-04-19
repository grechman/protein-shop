package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func InitDB() {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}
	Conn, err = pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func CloseDB() {
	Conn.Close(context.Background())
}
