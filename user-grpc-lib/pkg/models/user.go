// Package models contains all the domain models for the user service
package models

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

type PaginatedUsersModel struct {
	Users      []*UserModel
	Total      int64
	Page       int64
	PageSize   int64
	TotalPages int64
}
