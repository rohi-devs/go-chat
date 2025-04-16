package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func ConnectDB() {
	var err error
	connStr := "postgresql://neondb_owner:npg_54xTLPVuygmq@ep-wild-lab-a47ef400-pooler.us-east-1.aws.neon.tech/neondb?sslmode=require"
	DB, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	log.Println("Connected to PostgreSQL")
}

func InitTables() {
	_, err := DB.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);`)
	if err != nil {
		log.Fatal("Failed to initialize tables:", err)
	}
}
