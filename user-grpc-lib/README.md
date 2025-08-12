# User gRPC Library

A comprehensive gRPC library for user management services, designed to integrate seamlessly with existing Go applications without duplicating business logic.

## Features

- Complete gRPC service definition for user management operations
- Client library with connection management and timeouts
- Server implementation using adapter pattern to reuse existing business logic
- Comprehensive test coverage
- Semantic versioning for library releases
- Clean separation between gRPC layer and business logic

## Installation

```bash
go get github.com/grpc/user-grpc-lib
```

## Quick Start

### Server Setup

```go
import (
    "google.golang.org/grpc"
    pb "github.com/grpc/user-grpc-lib/pkg/pb/user/v1"
    "github.com/grpc/user-grpc-lib/pkg/server"
)

// Create adapter that implements UserServiceInterface
type UserServiceAdapter struct {
    existingService *YourExistingUserService
}

// Implement the interface methods...
func (a *UserServiceAdapter) CreateUser(ctx context.Context, input server.UserCreateInput) (*server.UserModel, error) {
    // Convert and delegate to existing service
    // ... implementation
}

// Setup gRPC server
grpcServer := grpc.NewServer()
userServiceServer := server.NewUserServiceServer(userServiceAdapter)
pb.RegisterUserServiceServer(grpcServer, userServiceServer)
```

### Client Usage

```go
import (
    "github.com/grpc/user-grpc-lib/pkg/client"
    pb "github.com/grpc/user-grpc-lib/pkg/pb/user/v1"
)

// Create client
config := client.DefaultClientConfig("localhost:50051")
userClient, err := client.NewUserServiceClient(config)
if err != nil {
    log.Fatal(err)
}
defer userClient.Close()

// Create user
resp, err := userClient.CreateUser(ctx, &pb.CreateUserRequest{
    Email:     "user@example.com",
    Password:  "password123",
    FirstName: "John",
    LastName:  "Doe",
})
```

## Project Structure

```
├── proto/user/v1/              # Protocol buffer definitions
│   └── user_service.proto      # gRPC service and message definitions
├── pkg/
│   ├── pb/                     # Generated protobuf files
│   ├── server/                 # gRPC server implementation
│   │   ├── user_service_server.go
│   │   └── user_service_server_test.go
│   └── client/                 # gRPC client library
│       └── user_service_client.go
└── Makefile                    # Build and code generation tasks
```

## API Reference

### User Service

The `UserService` provides the following gRPC methods:

- `CreateUser` - Create a new user account
- `GetUserByEmail` - Retrieve user by email address
- `GetUserByID` - Retrieve user by unique ID
- `GetUsers` - List users with pagination
- `UpdateUser` - Update user profile information
- `DeleteUser` - Delete a user account
- `Login` - Authenticate user and return JWT token
- `UpdateUserRole` - Update user's role (admin operation)
- `UpdatePassword` - Change user's password

### Message Types

#### User
```protobuf
message User {
    string id = 1;
    string email = 2;
    string first_name = 3;
    string last_name = 4;
    Role role = 5;
    google.protobuf.Timestamp created_at = 6;
    google.protobuf.Timestamp updated_at = 7;
    google.protobuf.Timestamp deleted_at = 8;
    int32 rating = 9;
}
```

#### Role Enum
```protobuf
enum Role {
    ROLE_UNSPECIFIED = 0;
    ROLE_USER = 1;
    ROLE_MODERATOR = 2;
    ROLE_ADMIN = 3;
}
```

## Integration Pattern

This library uses the Adapter Pattern to integrate with existing services:

1. **Interface Definition**: Define `UserServiceInterface` that matches your existing service methods
2. **Adapter Implementation**: Create an adapter that converts between gRPC types and your domain types
3. **Model Conversion**: Use the built-in `ModelConverter` to handle type conversions
4. **Error Mapping**: Automatic conversion of domain errors to appropriate gRPC status codes

## Development

### Prerequisites

- Go 1.19+
- Protocol Buffers compiler (`protoc`)
- gRPC Go plugins

### Setup

```bash
# Install dependencies
make deps

# Generate protobuf files
make proto

# Run tests
make test

# Build library
make build
```

### Code Generation

```bash
# Generate protobuf and gRPC files
make proto

# Clean generated files
make clean
```

## Testing

The library includes comprehensive test coverage:

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific test
go test -v ./pkg/server -run TestUserServiceServer_CreateUser
```

## Versioning

This library follows semantic versioning. To create a new release:

```bash
# Tag new version
git tag v1.0.0
git push origin v1.0.0

# View current tags
make tag
```

## Error Handling

The library automatically maps common business logic errors to appropriate gRPC status codes:

- `not found` → `codes.NotFound`
- `already exists` → `codes.AlreadyExists`
- `invalid credentials` → `codes.Unauthenticated`
- `insufficient rights` → `codes.PermissionDenied`
- `validation` → `codes.InvalidArgument`
- Others → `codes.Internal`

## Configuration

### Client Configuration

```go
config := &client.ClientConfig{
    Address:    "localhost:50051",
    TLS:        false,
    Timeout:    30 * time.Second,
    MaxRetries: 3,
    KeepAlive: &keepalive.ClientParameters{
        Time:                10 * time.Second,
        Timeout:             5 * time.Second,
        PermitWithoutStream: true,
    },
}
```

### TLS Configuration

```go
config := client.DefaultClientConfig("localhost:50051")
config.TLS = true
config.TLSConfig = &tls.Config{
    ServerName: "your-server.com",
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This library is released under the MIT License.
# Updated module path
