package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

var DbConnect *pgx.Conn

func init() {
	loadEnv()
	newDbInstance()
}

func loadEnv() {
	if _, ok := os.LookupEnv("DB_URL"); !ok {
		os.Setenv("DB_URL", "postgres://postgres:logic@localhost:5432/cli_chat_app")
		log.Print("DB env setup completed..")
	}
}

func newDbInstance() {
	dbURL := os.Getenv("DB_URL")

	dbConn1, err := pgx.Connect(context.Background(), dbURL)

	if err != nil {
		log.Print("Database Connection Failed\nDetails:" + err.Error())
	}

	DbConnect = dbConn1
}
