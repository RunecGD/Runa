package model

type Message struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	SenderID uint   `json:"sender_id"`
	Content  string `json:"content"`
}
