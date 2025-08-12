package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/grpc/user-grpc-lib/proto/user/v1"
)

// UserServiceInterface represents the existing HTTP service interface
type UserServiceInterface interface {
	CreateUser(ctx context.Context, input UserCreateInput) (*UserModel, error)
	GetUserByID(ctx context.Context, id string) (*UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	UpdateUser(ctx context.Context, id string, input UserUpdateInput) (*UserModel, error)
	DeleteUser(ctx context.Context, id string, actorID string, actorRole Role) error
	ListUsers(ctx context.Context, page, pageSize int64) (*PaginatedUsersModel, error)
	Login(ctx context.Context, email, password string) (string, error)
	UpdateUserRole(ctx context.Context, id string, role Role, actorRole Role) error
	UpdatePassword(ctx context.Context, id string, input UserPasswordUpdateInput) error
}

// Role represents user roles
type Role string

const (
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
	RoleAdmin     Role = "admin"
)

type UserCreateInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type UserUpdateInput struct {
	FirstName *string
	LastName  *string
}

type UserPasswordUpdateInput struct {
	CurrentPassword string
	NewPassword     string
}

// UserModel represents the domain user model
type UserModel struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Role      Role
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
	DeletedAt *timestamppb.Timestamp
	Rating    int32
}

// PaginatedUsersModel represents paginated user results
type PaginatedUsersModel struct {
	Users      []*UserModel
	Total      int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

// ModelConverter handles conversion between domain models and gRPC messages
type ModelConverter struct{}

func NewModelConverter() *ModelConverter {
	return &ModelConverter{}
}

// ConvertUserToProto converts domain user model to protobuf message
func (c *ModelConverter) ConvertUserToProto(user *UserModel) *pb.User {
	if user == nil {
		return nil
	}

	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      c.ConvertRoleToProto(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
		Rating:    user.Rating,
	}
}

// ConvertRoleToProto converts domain role to protobuf role
func (c *ModelConverter) ConvertRoleToProto(role Role) pb.Role {
	switch role {
	case RoleUser:
		return pb.Role_ROLE_USER
	case RoleModerator:
		return pb.Role_ROLE_MODERATOR
	case RoleAdmin:
		return pb.Role_ROLE_ADMIN
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}

// ConvertRoleFromProto converts protobuf role to domain role
func (c *ModelConverter) ConvertRoleFromProto(role pb.Role) Role {
	switch role {
	case pb.Role_ROLE_USER:
		return RoleUser
	case pb.Role_ROLE_MODERATOR:
		return RoleModerator
	case pb.Role_ROLE_ADMIN:
		return RoleAdmin
	default:
		return RoleUser
	}
}

// UserServiceServer implements the gRPC UserService server
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService UserServiceInterface
	converter   *ModelConverter
}

// NewUserServiceServer creates a new gRPC user service server
func NewUserServiceServer(userService UserServiceInterface) *UserServiceServer {
	return &UserServiceServer{
		userService: userService,
		converter:   NewModelConverter(),
	}
}

// CreateUser implements the CreateUser gRPC method
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	input := UserCreateInput{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := s.userService.CreateUser(ctx, input)
	if err != nil {
		return nil, s.convertError(err)
	}

	return &pb.CreateUserResponse{
		User: s.converter.ConvertUserToProto(user),
	}, nil
}

// GetUserByEmail implements the GetUserByEmail gRPC method
func (s *UserServiceServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	user, err := s.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, s.convertError(err)
	}

	return &pb.GetUserByEmailResponse{
		User: s.converter.ConvertUserToProto(user),
	}, nil
}

// GetUserByID implements the GetUserByID gRPC method
func (s *UserServiceServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	user, err := s.userService.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, s.convertError(err)
	}

	return &pb.GetUserByIDResponse{
		User: s.converter.ConvertUserToProto(user),
	}, nil
}

// GetUsers implements the GetUsers gRPC method
func (s *UserServiceServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	result, err := s.userService.ListUsers(ctx, page, pageSize)
	if err != nil {
		return nil, s.convertError(err)
	}

	users := make([]*pb.User, len(result.Users))
	for i, user := range result.Users {
		users[i] = s.converter.ConvertUserToProto(user)
	}

	return &pb.GetUsersResponse{
		Users:      users,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// UpdateUser implements the UpdateUser gRPC method
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	input := UserUpdateInput{}
	if req.FirstName != nil {
		input.FirstName = req.FirstName
	}
	if req.LastName != nil {
		input.LastName = req.LastName
	}

	user, err := s.userService.UpdateUser(ctx, req.Id, input)
	if err != nil {
		return nil, s.convertError(err)
	}

	return &pb.UpdateUserResponse{
		User: s.converter.ConvertUserToProto(user),
	}, nil
}

// DeleteUser implements the DeleteUser gRPC method
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id == "" || req.ActorId == "" {
		return nil, status.Error(codes.InvalidArgument, "id and actor_id are required")
	}

	actorRole := s.converter.ConvertRoleFromProto(req.ActorRole)

	err := s.userService.DeleteUser(ctx, req.Id, req.ActorId, actorRole)
	if err != nil {
		return nil, s.convertError(err)
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

// convertError converts domain errors to appropriate gRPC status codes
func (s *UserServiceServer) convertError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Map common business logic errors to gRPC status codes
	switch {
	case contains(errMsg, "not found"):
		return status.Error(codes.NotFound, errMsg)
	case contains(errMsg, "already exists"):
		return status.Error(codes.AlreadyExists, errMsg)
	case contains(errMsg, "invalid credentials"):
		return status.Error(codes.Unauthenticated, errMsg)
	case contains(errMsg, "insufficient rights"):
		return status.Error(codes.PermissionDenied, errMsg)
	case contains(errMsg, "validation"):
		return status.Error(codes.InvalidArgument, errMsg)
	default:
		return status.Error(codes.Internal, fmt.Sprintf("internal error: %v", err))
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsInside(s, substr))))
}

func containsInside(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
