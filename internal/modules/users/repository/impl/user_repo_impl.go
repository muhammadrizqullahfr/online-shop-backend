package impl

import (
	"context"
	"errors"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/repository"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/constants"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	*database.TransactionManagerImpl
}

func NewUserRepository(db *database.TransactionManagerImpl) repository.UserRepository {
	return &userRepositoryImpl{
		TransactionManagerImpl: db,
	}
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, user *models.User) error {
	return r.GetTx(ctx).Create(user).Error
}

func (r *userRepositoryImpl) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.GetTx(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.GetTx(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (r *userRepositoryImpl) FindByProvider(ctx context.Context, provider constants.Provider, providerID string) (*models.User, error) {
	var user models.User
	err := r.GetTx(ctx).Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, user *models.User) error {
	return r.GetTx(ctx).Model(user).Updates(user).Error
}

func (r *userRepositoryImpl) IsSeller(ctx context.Context, userID uint) (bool, error) {
	var result struct {
		ID uint
	}

	err := r.GetTx(ctx).
		Model(&models.User{}).
		Select("id").
		Where("id = ? AND role = ?", userID, "seller").
		Take(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *userRepositoryImpl) FindBySession(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	err := r.GetTx(ctx).Where("session = ?", token).First(&user).Error
	return &user, err
}

func (r *userRepositoryImpl) DeleteUser(ctx context.Context, id uint) error {
	return r.GetTx(ctx).Delete(&models.User{}, id).Error
}

func (r *userRepositoryImpl) ListUsers(ctx context.Context, limit, offset int, search string, role constants.UserRole) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.GetTx(ctx).Model(&models.User{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("email ILIKE ? OR full_name ILIKE ?", like, like)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, total, err
}
