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
	ID                    bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Strategy              AuthStrategy  `json:"strategy" bson:"strategy"`
	GoogleID              string        `json:"-" bson:"googleId,omitempty"`
	Email                 string        `json:"email" bson:"email"`
	Name                  string        `json:"name" bson:"name"`
	Password              string        `json:"-" bson:"password,omitempty"`
	Picture               string        `json:"picture" bson:"picture,omitempty"`
	RefreshToken          string        `json:"-" bson:"refreshToken,omitempty"`
	RefreshTokenExpiresAt *time.Time    `json:"-" bson:"refreshTokenExpiresAt,omitempty"`
	LastLoginAt           *time.Time    `json:"lastLoginAt" bson:"lastLoginAt"`
	CreatedAt             time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt             time.Time     `json:"updatedAt" bson:"updatedAt"`
}

type UserDTO struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Picture     string     `json:"picture"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
	CreatedAt   time.Time  `json:"createdAt"`
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:          u.ID.Hex(),
		Email:       u.Email,
		Name:        u.Name,
		Picture:     u.Picture,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
	}
}
