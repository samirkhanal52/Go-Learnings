package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/samirkhanal52/go-cli-chat-app/models"
)

var dbConn *pgx.Conn

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
		log.Print("Connection Failed\nDetails:" + err.Error())
	}

	dbConn = dbConn1
}

func InsertMessage(messageModel models.Chats) {
	results, err := dbConn.Query(context.Background(), "INSERT INTO chats(user_id, chat_room_id, chat_message, created_at) VALUES($1, $2, $3, now()::timestamp(0)) RETURNING id;", messageModel.UserID, messageModel.ChatRoomName, messageModel.Message)
	if err != nil {
		log.Fatalf("Error while executing query\nDetails:%v", err)
	}
	defer results.Close()

	response := []models.ResponseModel{}

	// iterate through the results
	for results.Next() {
		values, err := results.Values()
		if err != nil {
			log.Fatalf("Error while iterating data\nDetails:%v", err)
		}

		// convert DB types to Go types
		id := values[0].(int32)
		status_code := "200"
		status_message := "Data Saved Successfully"

		response = append(response, models.ResponseModel{
			StatusID:      int(id),
			StatusCode:    status_code,
			StatusMessage: status_message,
		})
	}

	if response[0].StatusCode == "200" {
		fmt.Println(response[0].StatusMessage + "\n")
	} else {
		log.Println(response[0].StatusMessage + "\n")
	}
}
