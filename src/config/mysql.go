package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DBConnect() {
	var err error

	DB_URL := os.Getenv("MYSQL_URL")

	DB, err = sql.Open("mysql", DB_URL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}

	// Check if the connection is alive
	if err = DB.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %s", err)
	}

	log.Println("Successfully connected to the database.")
}

func DBDisconnect() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Error disconnecting from the database: %s", err)
	}
	log.Println("Database connection closed.")
}
