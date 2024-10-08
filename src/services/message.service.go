package services

import (
	"encoding/json"
	"line/src/configs/rabbitmq"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type MessageService struct {
	DB *gorm.DB
}

type Message struct {
	// Id             uuid.UUID `json:"id"`
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();" json:"id"`
	SenderId       uuid.UUID `json:"sender_id"`
	ReceiverId     uuid.UUID `json:"receiver_id"`
	ConversationId uuid.UUID `json:"conversation_id"`
	Text           string    `json:"text"`
}

func (ms *MessageService) CreateMessage(senderId, receiverId, conversationId, text string) (map[string]interface{}, error) {
	message := Message{
		SenderId:       uuid.MustParse(senderId),
		ReceiverId:     uuid.MustParse(receiverId),
		ConversationId: uuid.MustParse(conversationId),
		Text:           text,
	}

	err := ms.DB.Create(&message).Error
	if err != nil {
		return nil, err
	}
	go PublishMessage(message)

	return map[string]interface{}{
		"message": "Message Created",
		"data":    message,
	}, nil
}

func (ms *MessageService) FindConversationMessages(conversationId string) (map[string]interface{}, error) {
	messages := []Message{}

	err := ms.DB.Where("conversation_id = ?", conversationId).Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "messages",
		"data":    messages,
	}, nil
}

func PublishMessage(msg Message) error {
	ch, err := rabbitmq.GetConnection().Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		"chat", // queue name
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBody,
		},
	)
	return err
}
