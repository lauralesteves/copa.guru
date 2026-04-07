package models

type GoogleLoginRequestDTO struct {
	Code        string `json:"code" validate:"required"`
	RedirectURI string `json:"redirectUri" validate:"required"`
}

type LoginResponseDTO struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	User         *UserDTO `json:"user"`
}

type RefreshRequestDTO struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshResponseDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
