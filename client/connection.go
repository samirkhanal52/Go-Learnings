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
		log.Print("Database Connection Failed\nDetails:" + err.Error())
	}

	dbConn = dbConn1
}

//register user and insert into users
func registerUser(userModel models.Users) {
	results, err := dbConn.Query(context.Background(), "SELECT * FROM create_user($1, $2, $3)", userModel.UserName, userModel.Password, userModel.Email)
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
		status_code := values[1].(string)
		status_message := values[2].(string)

		response = append(response, models.ResponseModel{
			StatusID:      int(id),
			StatusCode:    status_code,
			StatusMessage: status_message,
		})
	}

	if response[0].StatusCode == "200" {
		fmt.Println(response[0].StatusMessage)
	} else {
		log.Fatalln(response[0].StatusMessage)
	}
}

//autheticate user and login
func loginUser(userModel models.Users) {
	results, err := dbConn.Query(context.Background(), "SELECT * FROM user_authorization($1, $2)", userModel.UserName, userModel.Password)
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
		status_code := values[1].(string)
		status_message := values[2].(string)

		response = append(response, models.ResponseModel{
			StatusID:      int(id),
			StatusCode:    status_code,
			StatusMessage: status_message,
		})
	}

	if response[0].StatusCode == "200" {
		fmt.Println(response[0].StatusMessage)
		models.UserID = response[0].StatusID
		models.UserName = userModel.UserName
	} else {
		log.Println(response[0].StatusMessage)
		Login()
	}
}
