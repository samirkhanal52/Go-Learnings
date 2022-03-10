package connection

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	config "github.com/samirkhanal52/go-cli-chat-app"
	"github.com/samirkhanal52/go-cli-chat-app/models"
)

//register user and insert into users
func RegisterUser(userModel models.Users) []models.ResponseModel {
	query := "SELECT * FROM create_user($1, $2, $3)"
	results, err := config.DbConnect.Query(context.Background(), query, userModel.UserName, userModel.Password, userModel.Email)
	if err != nil {
		log.Fatalf("Error while executing query\nDetails:%v", err)
	}
	defer results.Close()

	return response(results)
}

//autheticate user and login
func LoginUser(userModel models.Users) []models.ResponseModel {
	query := "SELECT * FROM user_authorization($1, $2)"
	results, err := config.DbConnect.Query(context.Background(), query, userModel.UserName, userModel.Password)
	if err != nil {
		log.Fatalf("Error while executing query\nDetails:%v", err)
	}
	defer results.Close()

	return response(results)
}

func response(results pgx.Rows) []models.ResponseModel {
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

	return response
}
