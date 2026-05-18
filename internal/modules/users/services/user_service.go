package services

import (
	"context"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/delivery"
)

type UserService interface {
	CreateUser(ctx context.Context, req *delivery.CreateUserRequest) (*delivery.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*delivery.UserResponse, error)
	ListUsers(ctx context.Context, page, limit int, search, role string) (*delivery.UserListResponse, error)
	DeleteUser(ctx context.Context, id uint) error
}
