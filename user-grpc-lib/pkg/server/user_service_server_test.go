package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/Kiraberos/grpc/user-grpc-lib/proto/user/v1"
)

// MockUserService implements UserServiceInterface for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, input UserCreateInput) (*UserModel, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserModel), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*UserModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserModel), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*UserModel, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserModel), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id string, input UserUpdateInput) (*UserModel, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserModel), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string, actorID string, actorRole Role) error {
	args := m.Called(ctx, id, actorID, actorRole)
	return args.Error(0)
}

func (m *MockUserService) ListUsers(ctx context.Context, page, pageSize int64) (*PaginatedUsersModel, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PaginatedUsersModel), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) UpdateUserRole(ctx context.Context, id string, role Role, actorRole Role) error {
	args := m.Called(ctx, id, role, actorRole)
	return args.Error(0)
}

func (m *MockUserService) UpdatePassword(ctx context.Context, id string, input UserPasswordUpdateInput) error {
	args := m.Called(ctx, id, input)
	return args.Error(0)
}

func TestUserServiceServer_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		request        *pb.CreateUserRequest
		mockSetup      func(*MockUserService)
		expectedError  codes.Code
		expectedResult bool
	}{
		{
			name: "successful user creation",
			request: &pb.CreateUserRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, UserCreateInput{
					Email:     "test@example.com",
					Password:  "password123",
					FirstName: "John",
					LastName:  "Doe",
				}).Return(&UserModel{
					ID:        "123",
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Role:      RoleUser,
					CreatedAt: timestamppb.Now(),
					UpdatedAt: timestamppb.Now(),
				}, nil)
			},
			expectedResult: true,
		},
		{
			name: "missing email",
			request: &pb.CreateUserRequest{
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			mockSetup:     func(m *MockUserService) {},
			expectedError: codes.InvalidArgument,
		},
		{
			name: "user already exists",
			request: &pb.CreateUserRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errors.New("user already exists"))
			},
			expectedError: codes.AlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{}
			tt.mockSetup(mockService)

			server := NewUserServiceServer(mockService)

			resp, err := server.CreateUser(context.Background(), tt.request)

			if tt.expectedError != codes.OK {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedError, st.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)
				assert.Equal(t, tt.request.Email, resp.User.Email)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserServiceServer_GetUserByID(t *testing.T) {
	tests := []struct {
		name           string
		request        *pb.GetUserByIDRequest
		mockSetup      func(*MockUserService)
		expectedError  codes.Code
		expectedResult bool
	}{
		{
			name: "successful user retrieval",
			request: &pb.GetUserByIDRequest{
				Id: "123",
			},
			mockSetup: func(m *MockUserService) {
				m.On("GetUserByID", mock.Anything, "123").Return(&UserModel{
					ID:        "123",
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Role:      RoleUser,
				}, nil)
			},
			expectedResult: true,
		},
		{
			name: "missing id",
			request: &pb.GetUserByIDRequest{
				Id: "",
			},
			mockSetup:     func(m *MockUserService) {},
			expectedError: codes.InvalidArgument,
		},
		{
			name: "user not found",
			request: &pb.GetUserByIDRequest{
				Id: "999",
			},
			mockSetup: func(m *MockUserService) {
				m.On("GetUserByID", mock.Anything, "999").Return(nil, errors.New("user not found"))
			},
			expectedError: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{}
			tt.mockSetup(mockService)

			server := NewUserServiceServer(mockService)

			resp, err := server.GetUserByID(context.Background(), tt.request)

			if tt.expectedError != codes.OK {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedError, st.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)
				assert.Equal(t, tt.request.Id, resp.User.Id)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserServiceServer_Login(t *testing.T) {
	tests := []struct {
		name           string
		request        *pb.LoginRequest
		mockSetup      func(*MockUserService)
		expectedError  codes.Code
		expectedResult bool
	}{
		{
			name: "successful login",
			request: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockUserService) {
				m.On("Login", mock.Anything, "test@example.com", "password123").Return("jwt-token", nil)
				m.On("GetUserByEmail", mock.Anything, "test@example.com").Return(&UserModel{
					ID:        "123",
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Role:      RoleUser,
				}, nil)
			},
			expectedResult: true,
		},
		{
			name: "invalid credentials",
			request: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockUserService) {
				m.On("Login", mock.Anything, "test@example.com", "wrongpassword").Return("", errors.New("invalid credentials"))
			},
			expectedError: codes.Unauthenticated,
		},
		{
			name: "missing email",
			request: &pb.LoginRequest{
				Password: "password123",
			},
			mockSetup:     func(m *MockUserService) {},
			expectedError: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{}
			tt.mockSetup(mockService)

			server := NewUserServiceServer(mockService)

			resp, err := server.Login(context.Background(), tt.request)

			if tt.expectedError != codes.OK {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedError, st.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.Token)
				assert.NotNil(t, resp.User)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestModelConverter_ConvertRoleToProto(t *testing.T) {
	converter := NewModelConverter()

	tests := []struct {
		domainRole Role
		protoRole  pb.Role
	}{
		{RoleUser, pb.Role_ROLE_USER},
		{RoleModerator, pb.Role_ROLE_MODERATOR},
		{RoleAdmin, pb.Role_ROLE_ADMIN},
		{Role("invalid"), pb.Role_ROLE_UNSPECIFIED},
	}

	for _, tt := range tests {
		result := converter.ConvertRoleToProto(tt.domainRole)
		assert.Equal(t, tt.protoRole, result)
	}
}

func TestModelConverter_ConvertRoleFromProto(t *testing.T) {
	converter := NewModelConverter()

	tests := []struct {
		protoRole  pb.Role
		domainRole Role
	}{
		{pb.Role_ROLE_USER, RoleUser},
		{pb.Role_ROLE_MODERATOR, RoleModerator},
		{pb.Role_ROLE_ADMIN, RoleAdmin},
		{pb.Role_ROLE_UNSPECIFIED, RoleUser}, // Default fallback
	}

	for _, tt := range tests {
		result := converter.ConvertRoleFromProto(tt.protoRole)
		assert.Equal(t, tt.domainRole, result)
	}
}
