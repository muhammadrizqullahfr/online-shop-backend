package delivery

import "github.com/RakaMurdiarta/online-shop-system/pkg/common/constants"

type CreateUserRequest struct {
	Email    string             `json:"email" validate:"required,email,max=100"`
	Password string             `json:"password" validate:"required,min=6,max=72"`
	FullName string             `json:"full_name" validate:"required,max=255"`
	Role     constants.UserRole `json:"role" validate:"omitempty,oneof=customer admin"`
}
