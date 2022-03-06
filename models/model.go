package models

var (
	UserID   int
	UserName string
)

const ChannelName string = "chat"

type (
	Users struct {
		ID          int    `json:"id,omitempty"`
		UserName    string `json:"user_name,omitempty"`
		Email       string `json:"email,omitempty"`
		Password    string `json:"password,omitempty"`
		IsListening bool   `json:"is_listening"`
	}

	Chats struct {
		UserID       int    `json:"user_id,omitempty"`
		ChatRoomName string `json:"chat_room_id,omitempty"`
		Message      string `json:"message,omitempty"`
	}

	ResponseModel struct {
		StatusID      int    `json:"id,omitempty"`
		StatusCode    string `json:"status_code,omitempty"`
		StatusMessage string `json:"status_message"`
	}
)
