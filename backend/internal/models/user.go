package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthStrategy string

const (
	StrategyGoogle AuthStrategy = "google"
	StrategyEmail  AuthStrategy = "email"
)

type User struct {
	ID bson.ObjectID `json:"id" bson:"_id,omitempty"`

	Name    string `json:"name" bson:"name"`
	Email   string `json:"email" bson:"email"`
	Picture string `json:"picture" bson:"picture,omitempty"`

	Auth     Auth         `json:"auth" bson:"auth"`
	Strategy AuthStrategy `json:"strategy" bson:"strategy"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type Auth struct {
	GoogleID string `json:"-" bson:"googleId,omitempty"`
	Password string `json:"-" bson:"password,omitempty"`

	AccessToken           string     `json:"-" bson:"accessToken,omitempty"`
	RefreshToken          string     `json:"-" bson:"refreshToken,omitempty"`
	RefreshTokenExpiresAt *time.Time `json:"-" bson:"refreshTokenExpiresAt,omitempty"`
	LastLoginAt           *time.Time `json:"lastLoginAt" bson:"lastLoginAt"`
}

func (u *User) ToDTO() *UserDTO {
	dto := &UserDTO{
		ID:        u.ID.Hex(),
		Email:     u.Email,
		Name:      u.Name,
		Picture:   u.Picture,
		CreatedAt: u.CreatedAt,
	}

	if u.Auth.LastLoginAt != nil {
		dto.LastLoginAt = u.Auth.LastLoginAt
	}

	return dto
}

func (u *User) ToLoginResponseDTO() *LoginResponseDTO {
	return &LoginResponseDTO{
		AccessToken:  u.Auth.AccessToken,
		RefreshToken: u.Auth.RefreshToken,
		User:         u.ToDTO(),
	}
}

func (a *Auth) ToRefreshResponseDTO() *RefreshResponseDTO {
	return &RefreshResponseDTO{
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
	}
}
