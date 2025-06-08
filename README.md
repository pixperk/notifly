# Notifly: Multi-Channel Notification Microservice

Notifly is a robust, scalable microservice architecture designed for sending notifications through multiple channels like email and SMS. Built with Go, it implements modern software engineering practices including event-driven architecture, gRPC communication, and GraphQL API.

## Architecture Overview

Notifly is composed of the following microservices:

### User Service
- Handles user authentication and management
- Provides registration, login, and token validation
- Uses PostgreSQL for persistent storage
- Exposes gRPC API for internal communication

### Trigger Service
- Receives notification requests
- Validates authentication tokens
- Publishes notification events to NATS JetStream streams
- Implements authentication middleware for secure communication

### Notification Service
- Consumes notification events from NATS JetStream
- Implements durable, persistent subscriptions with explicit acknowledgements
- Dispatches notifications via appropriate channels (Email/SMS)
- Configurable retry policies based on notification type
- Intelligent error handling for different failure scenarios

### GraphQL API Gateway
- Provides a unified GraphQL API for client applications
- Handles authentication and request validation
- Proxies requests to appropriate microservices
- Implements GraphQL schema with mutations and queries

## Technology Stack

- **Language**: Go
- **API Protocols**: gRPC, GraphQL
- **Message Queue**: NATS with JetStream for persistent message storage
- **Database**: PostgreSQL
- **Authentication**: PASETO tokens
- **Notification Providers**:
  - Email: Brevo API
  - SMS: Twilio API
- **Containerization**: Docker & Docker Compose

## System Components

### Message Broker (NATS with JetStream)
NATS with JetStream serves as the central message broker, enabling loose coupling between services with persistent message storage:
- Notification events are published by the Trigger service to JetStream streams
- Events are consumed by Notification service using durable consumers
- Built-in persistent storage ensures message delivery even during service outages
- Automatic message replay and at-least-once delivery semantics
- Configurable retry policies with maximum delivery attempts

### Pull-Based Queue Implementation
The notification service implements a pull-based consumer pattern with JetStream:
- Uses explicit pull subscriptions rather than push-based delivery
- Consumer maintains state across restarts with durable name ("notif-consumer")
- Configurable batch processing with fetch size and wait time parameters
- DeliverAllPolicy ensures all messages in the stream are processed
- Messages are only removed from the queue after explicit acknowledgement
- Supports sophisticated error handling with different acknowledgement strategies:
  - Ack: For successfully processed messages
  - Nak: For temporary failures that should be retried
  - Term: For permanent failures that should not be retried
- Supports message inspection with metadata access for tracking delivery attempts
- Implements backoff strategy with configurable AckWait time
- Service restart automatically continues processing from last acknowledged message

### Message Processing
The notification service processes messages directly from JetStream:
- Explicit acknowledgement (Ack/Nak) for reliable message handling
- Configurable acknowledgement wait time and redelivery logic
- Automatic message filtering based on notification type
- Intelligent error handling with differentiated responses for transient vs. permanent failures

### Authentication
Token-based authentication using PASETO:
- Secure, stateless authentication using cookies
- gRPC interceptors for authorization
- Token validation middleware for GraphQL

## Getting Started

### Prerequisites
- Go 1.20 or higher
- Docker and Docker Compose
- Make (optional, for using Makefiles)

### Setup and Installation

1. Clone the repository:
   ```
   git clone https://github.com/pixperk/notifly.git
   cd notifly
   ```

2. Start required infrastructure:
   ```
   docker-compose up -d
   ```

3. Set up environment variables:
   - Copy the sample environment files in each service directory
   - Configure required API keys for Brevo and Twilio

4. Build and run services:
   ```
   cd user && make run
   cd trigger && make run
   cd notification && make run
   cd graphql && make run
   ```

## Service Configuration

Each service has its own configuration file (`app.env`) with the following settings:

### User Service
- Database connection parameters
- Token encryption key
- gRPC server port

### Trigger Service
- NATS connection details
- Token validation key
- gRPC server port

### Notification Service
- NATS connection details
- Email/SMS provider API keys
- Worker pool configuration

### GraphQL Service
- Internal service endpoints
- HTTP server configuration
- Authentication settings

## API Documentation

### GraphQL API

The GraphQL API provides the following operations:

#### Queries
- `healthCheck`: Simple health check endpoint

#### Mutations
- `signUp`: Register a new user
- `signIn`: Authenticate user and generate token
- `validateToken`: Validate an authentication token
- `triggerNotification`: Send a notification to a recipient

#### Notification Types
- `EMAIL`: Email notifications
- `SMS`: SMS notifications
- `PUSH`: Push notifications (planned for future)

## Architectural Patterns

### Event-Driven Architecture
- Services communicate asynchronously through events
- NATS JetStream provides persistent, reliable message delivery
- Supports at-least-once delivery semantics with message replay
- Enables service decoupling and independent scaling

### Circuit Breaker Pattern
- Gracefully handles service failures
- Implements configurable retry strategies through JetStream's delivery policies
- Prevents cascading failures with explicit message acknowledgement
- Differentiates between transient and permanent failures

### Repository Pattern
- Clean separation of data access layer
- Encapsulated database operations
- Facilitates testing with mock implementations

## Deployment

Notifly is designed for containerized deployment:

1. Build Docker images for each service:
   ```
   docker build -t notifly-user ./user
   docker build -t notifly-trigger ./trigger
   docker build -t notifly-notification ./notification
   docker build -t notifly-graphql ./graphql
   ```

2. Deploy using Docker Compose or Kubernetes:
   - Docker Compose for development/testing
   - Kubernetes recommended for production deployments

## Development

### Building from Source

Each service contains a Makefile with common commands:

```bash
make build    # Build the service binary
make run      # Run the service
make test     # Run tests
make proto    # Generate protobuf code (for gRPC services)
```

### Project Structure

```
├── common/              # Shared code and types
│   ├── auth/            # Authentication utilities
│   ├── client/          # gRPC client implementations
│   ├── proto/           # Protocol Buffers definitions
│   └── proto-gen/       # Generated gRPC code
├── graphql/             # GraphQL API gateway
│   ├── cmd/             # Service entry point
│   ├── generated/       # Generated GraphQL code
│   ├── models/          # GraphQL models
│   └── resolvers/       # GraphQL resolvers
├── notification/        # Notification dispatch service
│   ├── cmd/             # Service entry point
│   ├── dispatcher/      # Notification dispatchers
│   └── util/            # Utility functions
├── trigger/             # Notification trigger service
│   ├── cmd/             # Service entry point
│   └── util/            # Utility functions
├── user/                # User management service
│   ├── cmd/             # Service entry point
│   ├── db/              # Database related code
│   │   ├── migrations/  # SQL migrations
│   │   ├── queries/     # SQLC queries
│   │   └── sqlc/        # Generated database code
│   └── util/            # Utility functions
└── docker-compose.yml   # Local development setup
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
