# Line Chat API

### Project Overview
Line is a real-time chat API designed to facilitate instant messaging between users. Built on modern technologies, this API supports various functionalities, including user authentication, conversation management, and message retrieval, ensuring a seamless chat experience.

### Tech Stack
- **Go**: The programming language used for building the API.
- **Fiber**: A fast web framework for Go, enabling efficient routing and middleware support.
- **GraphQL**: A query language for APIs that allows clients to request only the data they need.
- **RabbitMQ**: A message broker that facilitates communication between services in a decoupled manner.
- **WebSockets**: For real-time communication between clients and the server.
- **PostgreSQL**: A robust relational database used to store user and message data.

### Features
- **User Authentication**: Sign up and sign in functionalities with token-based authentication.
- **Real-time Messaging**: Support for real-time chat through WebSocket connections.
- **User Management**: Ability to find all users with search capabilities.
- **Conversation Management**: Retrieve conversations and associated messages.
- **Token Validation**: All API requests require a valid JWT token for access control.
- **Pagination**: Efficiently manage large datasets through pagination in various queries

### Configuration Environment
Set the following environment variables in your configuration file:

```plaintext
PORT = 3000
HOST = http://localhost

DB_HOST = localhost
DB_USER = postgres
DB_PASSWORD = your_password
DB_NAME = line
DB_PORT = 5432

SECRET_KEY = my-secret-key

RABBITMQ_URL = amqp://guest:guest@localhost:5672/
```

Here’s the updated README section with a brief description of each route:

---

## GraphQL API List

The GraphQL API requires a token for all operations except for sign-up and sign-in resolvers. Below is a list of available routes (queries and mutations):

### Queries
- **ping**: A simple ping-pong test to check server connectivity.
- **findAllUsers**: Retrieves a paginated list of users.
- **findConversation**: Finds conversations between specified users.
- **findConversationMessages**: Retrieves messages from a specified conversation.
- **findMessage**: Fetches a specific message by its ID.

### Mutations
- **signUp**: Registers a new user (no token required).
- **signIn**: Authenticates a user and returns a token (no token required).
- **createConversation**: Initiates a new conversation between users.
- **createMessage**: Sends a new message in a specified conversation.

---

Let me know if you need any further modifications!

### WebSocket Connection
To connect to the WebSocket for real-time messaging, use the following URL structure, including the token:

```plaintext
ws://localhost:3000/ws?token=YOUR_JWT_TOKEN
```

### GraphQL Endpoint
The GraphQL API can be accessed at the following endpoint:

```plaintext
http://localhost:3000/graphql
```
