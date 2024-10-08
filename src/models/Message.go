package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();" json:"id"`
	SenderId       uuid.UUID `gorm:"type:uuid" json:"sender_id"`
	ReceiverId     uuid.UUID `gorm:"type:uuid" json:"receiver_id"`
	ConversationId uuid.UUID `gorm:"type:uuid" json:"conversation_id"`
	Text           string    `json:"text"`
	gorm.Model
}

// func (message *Message) BeforeCreate(tx *gorm.DB) (err error) {
// 	message.ID = uuid.New()
// 	fmt.Printf("Generated ID: %s\n", message.ID)
// 	return
// }
