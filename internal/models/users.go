package models

import (
	"time"

	"github.com/RakaMurdiarta/online-shop-system/pkg/common/constants"
	"gorm.io/gorm"
)

type User struct {
	ID         uint
	Email      string `gorm:"size:100;uniqueIndex;not null"`
	Password   *string
	FullName   string
	Role       constants.UserRole `gorm:"type:user_role_enum;default:'buyer'"`
	Provider   constants.Provider `gorm:"type:provider_enum;default:'local'"`
	ProviderID *string
	AvatarURL  *string
	Session    *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Carts []Cart
}
