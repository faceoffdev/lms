package database

import (
	"database/sql"
	"fmt"
	"lms/internal/config"
	"log"
)

func Connect(c config.PostgresConfig) *sql.DB {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable application_name=catizen_bot_go TimeZone=Europe/Moscow",
		c.Host, c.User, c.Pass, c.Dbname, c.Port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open Postgres: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}
	log.Printf("Connected to Postgres: %s", dsn)

	log.Print("PostgreSQL database is connected")

	return db
}
