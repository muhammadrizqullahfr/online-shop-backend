package impl

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/repository"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/constants"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"golang.org/x/crypto/bcrypt"
)

type userServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(
	userRepo repository.UserRepository,
) services.UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req *delivery.CreateUserRequest) (*delivery.UserResponse, error) {
	var created *models.User

	err := s.userRepo.WithTransaction(ctx, func(tx context.Context) error {
		existing, err := s.userRepo.GetUserByEmail(tx, req.Email)
		if err != nil {
			return err
		}
		if existing != nil {
			return errors.New("email already registered")
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}
		pass := string(hashed)

		role := req.Role
		if role == "" {
			role = constants.UserRole("buyer")
		}

		user := &models.User{
			Email:    req.Email,
			Password: &pass,
			FullName: req.FullName,
			Role:     role,
			Provider: constants.ProviderLocal,
		}

		if err := s.userRepo.CreateUser(tx, user); err != nil {
			return err
		}
		created = user
		return nil
	})
	if err != nil {
		return nil, err
	}

	return delivery.ToUserResponse(created), nil
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, id uint) (*delivery.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return delivery.ToUserResponse(user), nil
}

func (s *userServiceImpl) ListUsers(ctx context.Context, page, limit int, search, role string) (*delivery.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	users, total, err := s.userRepo.ListUsers(ctx, limit, offset, search, constants.UserRole(role))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	items := make([]delivery.UserResponse, 0, len(users))
	for i := range users {
		items = append(items, *delivery.ToUserResponse(&users[i]))
	}

	totalPage := 0
	if total > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &delivery.UserListResponse{
		Items: items,
		Paging: response.PagingResponse{
			CurrentPage: page,
			TotalPage:   totalPage,
			TotalItems:  total,
		},
	}, nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id uint) error {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	return s.userRepo.DeleteUser(ctx, id)
}
