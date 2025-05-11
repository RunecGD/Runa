package model

type Message struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}
