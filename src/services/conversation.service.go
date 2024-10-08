package services

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationService struct {
	DB *gorm.DB
}

type Conversation struct {
	Id         uuid.UUID `json:"id"`
	SenderId   string    `json:"sender_id"`
	ReceiverId string    `json:"receiver_id"`
}

func (cs *ConversationService) CreateConversation(senderId, receiverId string) (map[string]interface{}, error) {
	conversation := Conversation{
		SenderId:   senderId,
		ReceiverId: receiverId,
	}

	err := cs.DB.Create(&conversation).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"conversation": conversation,
		"message":      "Conversation Create Successfully",
	}, nil
}

func (cs *ConversationService) FindConversationBySenderAndReceiverId(userId, friendId string) ([]Conversation, error) {
	var conversations []Conversation
	err := cs.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userId, friendId, friendId, userId).First(&conversations).Error
	if err != nil {
		return nil, err
	}
	log.Printf("Conversations found: %+v", conversations)

	return conversations, nil
}
