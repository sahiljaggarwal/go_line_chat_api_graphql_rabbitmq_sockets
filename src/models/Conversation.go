package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();" json:"id"`
	SenderId uuid.UUID `gorm:"type:uuid" json:"sender_id"`
	ReceiverId uuid.UUID `gorm:"type:uuid" json:"receiver_id"`
	gorm.Model
}

func (conversation *Conversation) BeforeCreate(tx *gorm.DB)(err error){
	conversation.ID = uuid.New()
	return
}