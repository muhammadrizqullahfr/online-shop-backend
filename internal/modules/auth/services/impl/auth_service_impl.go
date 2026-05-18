package impl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/services"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/repository"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/constants"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"
	"golang.org/x/crypto/bcrypt"
)

type authServiceImpl struct {
	userRepo repository.UserRepository
	conf     *config.Config
}

func NewAuthService(userRepo repository.UserRepository, conf *config.Config,
) services.AuthService {

	return &authServiceImpl{
		userRepo: userRepo,
		conf:     conf,
	}

}

func (s *authServiceImpl) LoginLocal(ctx context.Context, req *delivery.LoginRequest) (*delivery.LoginResponse, error) {
	var resp delivery.LoginResponse

	err := s.userRepo.WithTransaction(ctx, func(tx context.Context) error {
		user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.New("invalid credential")
		}

		if user.Provider != "local" || user.Password == nil {
			return errors.New(fmt.Errorf("please login with your %v account", user.Provider).Error())
		}

		err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password))
		if err != nil {
			return errors.New("invalid credential")
		}

		token, err := shared.GenerateToken(uint(user.ID), user.Role, s.conf.JwtSecretKey, 24)
		if err != nil {
			return errors.New("failed to generate token")
		}

		tokenRefresh, err := shared.GenerateToken(uint(user.ID), user.Role, s.conf.JwtSecretKey, 168)

		if err != nil {
			return errors.New("failed to generate refresh token")
		}

		user.Session = &tokenRefresh
		err = s.userRepo.UpdateUser(ctx, user)
		if err != nil {
			return err
		}

		resp.User.ID = fmt.Sprintf("%d", user.ID)
		resp.User.Role = string(user.Role)
		resp.User.Email = user.Email
		resp.User.Name = user.FullName
		resp.Token = token

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// RegisterLocal implements [services.AuthService].
func (s *authServiceImpl) RegisterLocal(ctx context.Context, req *delivery.SingUpRequest) (*delivery.RegisterResponse, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("user Already Registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("password Hash Error")
	}
	passStr := string(hashedPassword)

	userDomain := req.ToDomain(&passStr, constants.RoleCustomer, constants.ProviderLocal)

	if err := s.userRepo.CreateUser(ctx, userDomain); err != nil {
		return nil, err
	}

	user := delivery.UserRegister{
		ID:        fmt.Sprintf("%d", userDomain.ID),
		Email:     userDomain.Email,
		Name:      userDomain.FullName,
		Role:      string(userDomain.Role),
		CreatedAt: userDomain.CreatedAt.Format(time.RFC3339),
	}

	return &delivery.RegisterResponse{
		User:  user,
		Token: "",
	}, nil
}

func (s *authServiceImpl) HandleOAuthCallback(ctx context.Context, data *delivery.OAuthUserRequest) (*models.User, error) {
	var userResponse *models.User

	pID := data.ProviderID
	avatar := data.AvatarURL

	err := s.userRepo.WithTransaction(ctx, func(tx context.Context) error {
		user, err := s.userRepo.FindByProvider(tx, data.Provider, pID)
		if err == nil && user != nil {
			userResponse = user
			return nil
		}

		existingUser, err := s.userRepo.GetUserByEmail(tx, data.Email)
		if err == nil && existingUser != nil {
			existingUser.Provider = data.Provider
			existingUser.ProviderID = &pID

			if existingUser.AvatarURL == nil || *existingUser.AvatarURL == "" {
				existingUser.AvatarURL = &avatar
			}

			if err := s.userRepo.UpdateUser(tx, existingUser); err != nil {
				return err
			}
			userResponse = existingUser
			return nil
		}

		newUser := &models.User{
			Email:      data.Email,
			FullName:   data.FullName,
			Provider:   data.Provider,
			ProviderID: &pID,
			AvatarURL:  &avatar,
			Role:       constants.RoleCustomer,
		}

		if err := s.userRepo.CreateUser(tx, newUser); err != nil {
			return err
		}

		userResponse = newUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	return userResponse, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*delivery.RefreshTokenResponse, error) {

	claims, err := shared.ValidateToken(refreshToken, s.conf.JwtSecretKey)

	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	user, err := s.userRepo.FindBySession(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("session not found or already logged out")
	}

	if user.Session == nil {
		return nil, errors.New("invalid session")
	}

	newAccessToken, err := shared.GenerateToken(claims.UserID, claims.Role, s.conf.JwtSecretKey, 24)
	if err != nil {
		return nil, err
	}

	return &delivery.RefreshTokenResponse{
		NewAccessToken: newAccessToken,
	}, nil
}

func (s *authServiceImpl) Me(ctx context.Context, userID uint) (*delivery.MeResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &delivery.MeResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.FullName,
		Role:  user.Role,
	}, nil
}
