package handlers

import (
	"encoding/json"
	"errors"
	"line/src/configs/rabbitmq"
	"line/src/models"

	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var clients = make(map[uuid.UUID]*websocket.Conn)

func HandleWebSocket(c *websocket.Conn) {
	defer c.Close()

	token := c.Query("token")
	claims, err := ValidateToken(token)
	if err != nil {
		log.Println("Token validation failed:", err)
		c.WriteMessage(websocket.CloseMessage, []byte("Unauthorized"))
		return
	}

	userID, ok := claims["id"].(string)
	log.Println("userId ", userID)
	if !ok {
		log.Println("User ID not found in claims")
		c.WriteMessage(websocket.CloseMessage, []byte("Unauthorized"))
		return
	}

	userUUID := uuid.MustParse(userID)
	clients[userUUID] = c
	defer delete(clients, userUUID)
	consumeMessagesForUser(userUUID)
}

func consumeMessagesForUser(userUUID uuid.UUID) {
	ch, err := rabbitmq.GetConnection().Channel()
	if err != nil {
		log.Println("Error getting RabbitMQ Channel:", err)
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"chat",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to register a consume:", err)
		return
	}

	for msg := range msgs {
		var message models.Message
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		if message.ReceiverId == userUUID {
			client, exists := clients[userUUID]
			if exists {
				err := client.WriteJSON(message)
				if err != nil {
					log.Println("Error sending message to client:", err)
				} else {
					log.Printf("Message sent to user: %s", userUUID)
				}
			} else {
				log.Printf("User %s not connected", userUUID)
			}
		}
	}
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, errors.New("authorization token is required")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("my-secret-key"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
