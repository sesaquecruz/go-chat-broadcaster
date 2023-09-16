package model

type Message struct {
	Id         string `json:"id"`
	RoomId     string `json:"room_id"`
	SenderId   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Text       string `json:"text"`
	CreatedAt  string `json:"created_at"`
}
