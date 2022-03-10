package connection

import (
	"context"
	"log"

	config "github.com/samirkhanal52/go-cli-chat-app"
	"github.com/samirkhanal52/go-cli-chat-app/models"
)

func InsertMessage(messageModel models.Chats) []models.ResponseModel {
	query := "INSERT INTO chats(user_id, chat_room_id, chat_message, created_at) VALUES($1, $2, $3, now()::timestamp(0)) RETURNING id;"
	results, err := config.DbConnect.Query(context.Background(), query, messageModel.UserID, messageModel.ChatRoomName, messageModel.Message)
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

	return response
}
