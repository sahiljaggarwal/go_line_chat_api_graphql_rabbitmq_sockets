package graphqlschema

import (
	"errors"
	"line/src/configs/db"
	"line/src/configs/env"
	"line/src/services"

	// "line/src/services"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/graphql-go/graphql"
	gql "github.com/graphql-go/graphql"
)

/*********** Types ***********/

// User Type
var userType = gql.NewObject(gql.ObjectConfig{
	Name: "User",
	Fields: gql.Fields{
		"id": &gql.Field{
			Type: gql.String,
		},
		"name": &gql.Field{
			Type: gql.String,
		},
		"email": &gql.Field{
			Type: gql.String,
		},
		"online_status": &gql.Field{
			Type: gql.Boolean,
		},
		"profile_image": &gql.Field{
			Type: gql.String,
		},
		"token": &gql.Field{
			Type: gql.String,
		},
	},
})

// SignUp Type
var signUpInputType = gql.NewInputObject(gql.InputObjectConfig{
	Name: "SignUpInput",
	Fields: gql.InputObjectConfigFieldMap{
		"name":     &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
		"email":    &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
		"password": &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
	},
})

// SignIn Type
var signInInputType = gql.NewInputObject(gql.InputObjectConfig{
	Name: "SignInInput",
	Fields: gql.InputObjectConfigFieldMap{
		"email":    &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
		"password": &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
		// "id": &gql.InputObjectFieldConfig{Type: gql.String},
		// "token": &gql.InputObjectFieldConfig{Type: gql.String},

	},
})

// Find All Users Type
var findAllUsersInputType = gql.NewInputObject(gql.InputObjectConfig{
	Name: "FindAllUsersInput",
	Fields: gql.InputObjectConfigFieldMap{
		"searchQuery": &gql.InputObjectFieldConfig{
			Type: gql.String,
		},
		"limit": &gql.InputObjectFieldConfig{
			Type: gql.Int,
		},
		"offset": &gql.InputObjectFieldConfig{
			Type: gql.Int,
		},
	},
})

// Pagination Type
var paginatedUsersType = gql.NewObject(gql.ObjectConfig{
	Name: "PaginatedUsers",
	Fields: gql.Fields{
		"totalCount": &gql.Field{
			Type: gql.Int,
		},
		"users": &gql.Field{
			Type: gql.NewList(userType), // List of users
		},
	},
})

// Message Type
var messageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Message",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String), // Ensure ID is non-null
		},
		"sender_id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String), // Ensure sender_id is non-null
		},
		"conversation_id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String), // Ensure conversation_id is non-null
		},
		"text": &graphql.Field{
			Type: graphql.String, // text can be nullable
		},
		"created_at": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String), // Ensure created_at is non-null
		},
	},
})

/************ Query Resolvers ************/
var rootQuery = gql.NewObject(gql.ObjectConfig{
	Name: "Query",
	Fields: gql.Fields{
		// testing ping pong
		"ping": &gql.Field{
			Type: gql.String,
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				return "pong", nil
			},
		},
		// find all users
		"findAllUsers": &gql.Field{
			Type: paginatedUsersType,
			Args: gql.FieldConfigArgument{
				"input": &gql.ArgumentConfig{
					Type: findAllUsersInputType,
				},
			},
			Resolve: func(p gql.ResolveParams) (interface{}, error) {

				fiberCtx, ok := p.Context.Value("fiberCtx").(*fiber.Ctx)
				if !ok || fiberCtx == nil {
					return nil, errors.New("failed to get fiber context")
				}

				token := fiberCtx.Locals("Authorization")

				if token != nil {
					tokenString, ok := token.(string)
					if !ok {
						return nil, errors.New("invalid token format")
					}

					_, err := ValidateToken(tokenString)
					if err != nil {
						return nil, err
					}
				} else {
					log.Println("No authorization token provided, proceeding without token")
				}

				searchQuery := ""
				limit := 10
				offset := 0

				if input, ok := p.Args["input"].(map[string]interface{}); ok {
					if val, ok := input["searchQuery"].(string); ok {
						searchQuery = val
					}
					if val, ok := input["limit"].(int); ok {
						limit = val
					}
					if val, ok := input["offset"].(int); ok {
						offset = val
					}
				} else {
					log.Print("No input provided, using default values")
				}

				log.Printf("Searching users with query: %s, limit: %d, offset: %d", searchQuery, limit, offset)

				userService := services.UserService{DB: db.DB}

				if userService.DB == nil {
					return nil, errors.New("database connection is nil")
				}

				usersData, err := userService.FindAllUsers(searchQuery, limit, offset)
				if err != nil {
					log.Print("Error fetching users", err)
					return nil, err
				}

				return usersData, nil
			},
		},

		// find conversations
		"findConversation": &gql.Field{
			Type: gql.NewList(gql.NewObject(gql.ObjectConfig{
				Name: "Conversation",
				Fields: gql.Fields{
					"receiver_id": &gql.Field{Type: gql.String},
					"id":          &gql.Field{Type: gql.String},
				},
			})),
			Args: gql.FieldConfigArgument{
				"friendId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.String)},
			},
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				fiberCtx, ok := p.Context.Value("fiberCtx").(*fiber.Ctx)
				if !ok || fiberCtx == nil {
					return nil, errors.New("failed to get fiber context")
				}

				token := fiberCtx.Locals("Authorization")
				if token == nil {
					return nil, errors.New("authorization token is required")
				}

				tokenString, ok := token.(string)
				if !ok {
					return nil, errors.New("invalid token format")
				}

				claims, err := ValidateToken(tokenString)
				if err != nil {
					return nil, err
				}

				userId, ok := claims["id"].(string)
				if !ok {
					return nil, errors.New("user ID not found in token claims")
				}

				friendId := p.Args["friendId"].(string)

				conversationService := services.ConversationService{DB: db.DB}
				conversations, err := conversationService.FindConversationBySenderAndReceiverId(userId, friendId)
				if err != nil {
					log.Print("Error finding conversations", err)
					return nil, err
				}

				return conversations, nil
			},
		},

		// find conversation messages
		"findConversationMessages": &gql.Field{
			Type: gql.NewList(messageType),
			Args: gql.FieldConfigArgument{
				"conversationId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.String)},
			},
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				conversationId := p.Args["conversationId"].(string)

				messageService := services.MessageService{DB: db.DB}
				messagesData, err := messageService.FindConversationMessages(conversationId)
				if err != nil {
					log.Print("Error fetching messages", err)
					return nil, err
				}

				return messagesData["data"], nil
			},
		},
	},
})

/************ Mutation Resolvers ************/
var mutation = gql.NewObject(gql.ObjectConfig{
	Name: "Mutation",
	Fields: gql.Fields{

		// Signup Resolver
		"signUp": &gql.Field{
			Type: userType,
			Args: gql.FieldConfigArgument{
				"input": &gql.ArgumentConfig{
					Type: gql.NewNonNull(signUpInputType),
				},
			},
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				input := p.Args["input"].(map[string]interface{})
				name := input["name"].(string)
				email := input["email"].(string)
				password := input["password"].(string)

				userService := services.UserService{DB: db.DB}
				userData, err := userService.CreateUser(name, email, password)
				if err != nil {
					log.Print("Error creating user", err)
					return nil, err
				}

				return userData["user"], nil
			},
		},

		// Signin Resolver
		"signIn": &gql.Field{
			Type: userType,
			Args: gql.FieldConfigArgument{
				"input": &gql.ArgumentConfig{
					Type: gql.NewNonNull(signInInputType),
				},
			},
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				input := p.Args["input"].(map[string]interface{})
				email := input["email"].(string)
				password := input["password"].(string)

				userService := services.UserService{DB: db.DB}
				userData, err := userService.SigninUser(email, password)
				if err != nil {
					return nil, err
				}
				log.Print("service data ", userData)
				return userData, nil
			},
		},

		// Create Conversation Resolver
		"createConversation": &gql.Field{
			Type: gql.NewObject(gql.ObjectConfig{
				Name: "CreateConversationResponse",
				Fields: gql.Fields{
					"conversation": &gql.Field{
						Type: gql.NewObject(gql.ObjectConfig{
							Name: "Conversation",
							Fields: gql.Fields{
								"sender_id":   &gql.Field{Type: gql.String},
								"receiver_id": &gql.Field{Type: gql.String},
							},
						}),
					},
					"message": &gql.Field{Type: gql.String},
				},
			}),
			Args: gql.FieldConfigArgument{
				// "senderId":   &gql.ArgumentConfig{Type: gql.NewNonNull(gql.String)},
				"receiverId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.String)},
			},
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				fiberCtx, ok := p.Context.Value("fiberCtx").(*fiber.Ctx)
				if !ok || fiberCtx == nil {
					return nil, errors.New("failed to get fiber context")
				}

				// Validate token
				token := fiberCtx.Locals("Authorization")
				var claims jwt.MapClaims
				if token == nil {
					return nil, errors.New("authorization token is required")
				}

				tokenString, ok := token.(string)
				if !ok {
					return nil, errors.New("invalid token format")
				}

				claims, err := ValidateToken(tokenString)
				if err != nil {
					return nil, err
				}

				senderId := claims["id"].(string)
				receiverId := p.Args["receiverId"].(string)

				conversationService := services.ConversationService{DB: db.DB}
				conversationData, err := conversationService.CreateConversation(senderId, receiverId)
				if err != nil {
					log.Print("Error creating conversation", err)
					return nil, err
				}

				return conversationData, nil
			},
		},

		// Create Message Resolver
		"createMessage": &gql.Field{
			Type: gql.NewObject(gql.ObjectConfig{
				Name: "CreateMessageResponse",
				Fields: gql.Fields{
					"message": &gql.Field{Type: gql.String},
					"data": &gql.Field{Type: gql.NewObject(gql.ObjectConfig{
						Name: "Message",
						Fields: gql.Fields{
							"id":              &gql.Field{Type: gql.String},
							"sender_id":       &gql.Field{Type: gql.String},
							"receiver_id":     &gql.Field{Type: gql.String},
							"conversation_id": &gql.Field{Type: gql.String},
							"text":            &gql.Field{Type: gql.String},
						},
					})},
				},
			}),
			Args: gql.FieldConfigArgument{
				"input": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.NewInputObject(gql.InputObjectConfig{
						Name: "CreateMessageInput",
						Fields: gql.InputObjectConfigFieldMap{
							"receiver_id":     &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
							"conversation_id": &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
							"text":            &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
						},
					})),
				},
			},
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				input := p.Args["input"].(map[string]interface{})

				// log.Printf("Received inputs: %+v", input)

				// senderId := input["senderId"].(string)
				conversationId := input["conversation_id"].(string)
				text := input["text"].(string)
				// log.Printf("Received tex id inputs: %+v %v", text, conversationId)

				fiberCtx, ok := p.Context.Value("fiberCtx").(*fiber.Ctx)
				if !ok || fiberCtx == nil {
					return nil, errors.New("failed to get fiber context")
				}

				// Validate token
				token := fiberCtx.Locals("Authorization")
				// log.Print("token ", token)
				var claims jwt.MapClaims
				if token == nil {
					return nil, errors.New("authorization token is required")
				}

				tokenString, ok := token.(string)
				if !ok {
					return nil, errors.New("invalid token format")
				}

				claims, err := ValidateToken(tokenString)
				if err != nil {
					return nil, err
				}
				senderId := claims["id"].(string)
				receiverId := input["receiver_id"].(string)
				messageService := services.MessageService{DB: db.DB}
				response, err := messageService.CreateMessage(senderId, receiverId, conversationId, text)
				if err != nil {
					log.Print("Error creating message:", err)
					return nil, err
				}

				return response, nil
			},
		},
	},
})

// Schema configuration
var Schema, _ = gql.NewSchema(gql.SchemaConfig{
	Query:    rootQuery,
	Mutation: mutation,
})

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, errors.New("authorization token is required")
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, errors.New("invalid authorization format")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(env.SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
