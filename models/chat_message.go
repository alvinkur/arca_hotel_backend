package models

import "time"

type ChatMessage struct {
	ID         uint      `gorm:"primaryKey;column:id_chat_message"       json:"id_chat_message"`
	ChatID     uint      `gorm:"column:id_chat;not null;index"           json:"id_chat"`
	SenderType string    `gorm:"not null"                                  json:"sender_type"`
	Message    string    `gorm:"not null"                                  json:"message"`
	Date       time.Time `gorm:"autoCreateTime"                            json:"date"`

	Chat Chat `gorm:"foreignKey:ChatID;references:ID" json:"chat,omitempty"`
}

func (ChatMessage) TableName() string {
	return "chat_message"
}
