package middleware

import (
	"bufio"
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
		os.Setenv("PORT", ":4444")
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

func Register() {
	fmt.Print("Welcome\nRegister\nAlready User?(Y/N):")

	reader := bufio.NewScanner(os.Stdin)

	if err := reader.Err(); err != nil {
		log.Print("Error..")
	}

	for reader.Scan() {
		userEntry := reader.Text()

		if userEntry == "Y" {
			return
		} else if userEntry == "N" {
			break
		} else {
			fmt.Print("Invalid Command(Y/N):")
		}
	}

	user := models.Users{}
	//Username
	fmt.Print("User Name:")
	for reader.Scan() {
		userEntry := reader.Text()

		if userEntry == "quit" {
			log.Print("Quiting..")
			return
		} else if userEntry != "" {
			user.UserName = userEntry
			break
		} else if userEntry == "" {
			fmt.Print("Please Enter User Name:")
		}
	}

	//Password
	fmt.Print("Password:")
	for reader.Scan() {
		userEntry := reader.Text()

		if userEntry == "quit" {
			log.Print("Quiting..")
			return
		} else if userEntry != "" {
			user.Password = userEntry
			break
		} else if userEntry == "" {
			fmt.Print("Please Enter Password:")
		}
	}

	//Email
	fmt.Print("Email:")
	for reader.Scan() {
		userEntry := reader.Text()

		if userEntry == "quit" {
			log.Print("Quiting..")
			break
		} else if userEntry != "" {
			user.Email = userEntry
			break
		} else if userEntry == "" {
			fmt.Print("Please Enter Email:")
		}
	}

	registerUser(user)
}

func Login() {
	Register()

	fmt.Println("Login")

	reader := bufio.NewScanner(os.Stdin)

	if err := reader.Err(); err != nil {
		log.Print("Error..")
	}

	user := models.Users{}
	//Username
	fmt.Print("User Name:")
	for reader.Scan() {
		userEntry := reader.Text()

		if userEntry == "quit" {
			log.Print("Quiting..")
			return
		} else if userEntry != "" {
			user.UserName = userEntry
			break
		} else if userEntry == "" {
			fmt.Print("Please Enter User Name:")
		}
	}

	//Password
	fmt.Print("Password:")
	for reader.Scan() {
		userEntry := reader.Text()

		if userEntry == "quit" {
			log.Print("Quiting..")
			return
		} else if userEntry != "" {
			user.Password = userEntry
			break
		} else if userEntry == "" {
			fmt.Print("Please Enter Password:")
		}
	}

	loginUser(user)
}

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
