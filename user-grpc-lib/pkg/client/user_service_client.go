package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	pb "github.com/Kiraberos/grpc/user-grpc-lib/proto/user/v1"
)

// Config ClientConfig holds configuration for the gRPC client
type Config struct {
	Address    string
	TLS        bool
	TLSConfig  *tls.Config
	Timeout    time.Duration
	MaxRetries int
	KeepAlive  *keepalive.ClientParameters
}

// DefaultClientConfig returns a default client configuration
func DefaultClientConfig(address string) *Config {
	return &Config{
		Address:    address,
		TLS:        false,
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		KeepAlive: &keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		},
	}
}

// UserServiceClient wraps the generated gRPC client with additional functionality
type UserServiceClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
	config *Config
}

// NewUserServiceClient creates a new user service client
func NewUserServiceClient(config *Config) (*UserServiceClient, error) {
	if config == nil {
		return nil, fmt.Errorf("client config is required")
	}

	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(*config.KeepAlive),
	}

	// Configure TLS or insecure connection
	if config.TLS {
		creds := credentials.NewTLS(config.TLSConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(config.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := pb.NewUserServiceClient(conn)

	return &UserServiceClient{
		client: client,
		conn:   conn,
		config: config,
	}, nil
}

// Close closes the client connection
func (c *UserServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreateUser creates a new user
func (c *UserServiceClient) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.client.CreateUser(ctx, req)
}

// GetUserByEmail retrieves a user by email
func (c *UserServiceClient) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.client.GetUserByEmail(ctx, req)
}

// GetUserByID retrieves a user by ID
func (c *UserServiceClient) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.client.GetUserByID(ctx, req)
}

// GetUsers retrieves a list of users with pagination
func (c *UserServiceClient) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.client.GetUsers(ctx, req)
}

// UpdateUser updates an existing user
func (c *UserServiceClient) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.client.UpdateUser(ctx, req)
}

// DeleteUser deletes a user
func (c *UserServiceClient) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	return c.client.DeleteUser(ctx, req)
}

// withTimeout adds a timeout to the context if one isn't already set
func (c *UserServiceClient) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); hasDeadline {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, c.config.Timeout)
}

// HealthCheck performs a simple health check by calling GetUsers with minimal parameters
func (c *UserServiceClient) HealthCheck(ctx context.Context) error {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	_, err := c.client.GetUsers(ctx, &pb.GetUsersRequest{
		Page:     1,
		PageSize: 1,
	})
	return err
}
