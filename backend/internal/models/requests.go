package models

type GoogleLoginRequest struct {
	Code        string `json:"code"`
	RedirectURI string `json:"redirectUri"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}
