package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/samirkhanal52/go-cli-chat-app/client/connection"
	"github.com/samirkhanal52/go-cli-chat-app/models"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

var client *socketio_client.Client

func main() {
	Login()

	fmt.Println("WELCOME " + models.UserName)

	initSocket()

	fmt.Print("Write Message:")

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "quit" {
			return
		}

		if err := client.Emit("message", models.UserName+":"+command); err != nil {
			log.Printf("Emit message error:%v\n", err)
			continue
		}
	}
}

//Initialize socket
func initSocket() {
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["uid"] = "1"
	opts.Query["cid"] = "conf_123"

	//Get HOST URI
	if _, ok := os.LookupEnv("HOST"); !ok {
		os.Setenv("HOST", "http://127.0.0.1:4444")
	}

	clientTemp, err := socketio_client.NewClient(os.Getenv("HOST"), opts)

	if err != nil {
		log.Printf("New Client Error:%v\n", err)
		os.Exit(0)
		return
	}

	client = clientTemp

	client.On("error", func() {
		log.Printf("Error\n")
		os.Exit(0)
	})

	client.On("connection", func() {
		log.Printf("on connect\n")
	})

	client.On("message", func(msg string) {
		log.Printf("%v\n", msg)
	})

	client.On("disconnection", func() {
		log.Printf("Disconnected\n")
	})

	go func() {
		authStr := "{\"uid\":\"" + opts.Query["uid"] + "\",\"cid\":\"" + opts.Query["cid"] + "\"}"
		for {
			err := client.Emit("authenticate", authStr)
			if err != nil {
				log.Printf("Emit auth error:%v\n", err)
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

//User Login
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

	response := connection.LoginUser(user)

	if response[0].StatusCode == "200" {
		fmt.Println(response[0].StatusMessage)
		models.UserID = response[0].StatusID
		models.UserName = user.UserName
	} else {
		log.Println(response[0].StatusMessage)
		Login()
	}
}

//User Registration
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

	response := connection.RegisterUser(user)

	if response[0].StatusCode == "200" {
		fmt.Println(response[0].StatusMessage)
	} else {
		log.Fatalln(response[0].StatusMessage)
	}
}
