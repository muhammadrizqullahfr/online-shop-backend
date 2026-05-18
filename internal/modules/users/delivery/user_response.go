package delivery

import (
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/constants"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
)

type UserResponse struct {
	ID        uint               `json:"id"`
	Email     string             `json:"email"`
	FullName  string             `json:"name"`
	Role      constants.UserRole `json:"role"`
	Provider  constants.Provider `json:"provider"`
	AvatarURL string             `json:"avatar_url"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type UserListResponse struct {
	Items  []UserResponse          `json:"items"`
	Paging response.PagingResponse `json:"paging"`
}

func ToUserResponse(u *models.User) *UserResponse {
	if u == nil {
		return nil
	}
	avatar := ""
	if u.AvatarURL != nil {
		avatar = *u.AvatarURL
	}
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Role:      u.Role,
		Provider:  u.Provider,
		AvatarURL: avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
