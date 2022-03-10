package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/samirkhanal52/go-cli-chat-app/models"
	"github.com/samirkhanal52/go-cli-chat-app/server/connection"
)

func main() {
	socketServer := socketio.NewServer(nil)

	socketServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("Connected ID:", s.ID())
		s.Join(models.ChannelName)
		socketServer.BroadcastToRoom("/", models.ChannelName, "message", models.UserName+" joined the chat.")

		socketServer.OnEvent("/", "message", func(s socketio.Conn, msg string) {
			s.SetContext(msg)
			fmt.Println("message:", msg)
			socketServer.BroadcastToRoom("/", models.ChannelName, "message", msg)

			chatMessage := models.Chats{
				UserID:       models.UserID,
				ChatRoomName: models.ChannelName,
				Message:      msg,
			}

			response := connection.InsertMessage(chatMessage)

			if response[0].StatusCode == "200" {
				fmt.Println(response[0].StatusMessage + "\n")
			} else {
				log.Println(response[0].StatusMessage + "\n")
			}

			socketServer.OnError("/", func(s socketio.Conn, e error) {
				fmt.Println("Error:", e)
			})

			socketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
				fmt.Println("Closed", reason)
			})
		})

		return nil
	})

	//Get Host URL
	if _, ok := os.LookupEnv("PORT"); !ok {
		os.Setenv("PORT", ":4444")
		log.Print("Host connected..")
	}

	srvMux := http.NewServeMux()
	srvMux.Handle("/socket.io/", socketServer)

	go func() {
		if err := socketServer.Serve(); err != nil {
			log.Fatalf("Socket Listen error: %s\n", err)
		}
	}()

	defer socketServer.Close()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         os.Getenv("PORT"),
		Handler:      srvMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		fmt.Println("Listening on port...", os.Getenv("PORT"))

		if err := http.ListenAndServe(os.Getenv("PORT"), srvMux); err != nil {
			log.Fatal(err)
		}

	}()

	<-stopChan
	log.Println("Server Shutting Down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	fmt.Println("Server Shut Down...")
}
