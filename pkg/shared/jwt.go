package shared

import (
	"errors"
	"time"

	"github.com/RakaMurdiarta/online-shop-system/pkg/common/constants"
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	UserID uint               `json:"user_id"`
	Role   constants.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role constants.UserRole, secret string, expireInHours int) (string, error) {

	claims := &JwtCustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireInHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ecommerce-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, secret string) (*JwtCustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signin token")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token or expire")
}
