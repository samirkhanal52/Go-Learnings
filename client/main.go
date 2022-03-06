package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/samirkhanal52/go-cli-chat-app/middleware"
	"github.com/samirkhanal52/go-cli-chat-app/models"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

func main() {

	middleware.Login()

	fmt.Println("WELCOME " + models.UserName)

	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["uid"] = "1"
	opts.Query["cid"] = "conf_123"
	uri := "http://127.0.0.1:4444"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Printf("New Client Error:%v\n", err)
		return
	}

	client.On("error", func() {
		log.Printf("Error\n")
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
