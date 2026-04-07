package models

import "fmt"

func (r *GoogleLoginRequest) Validate() error {
	if r.Code == "" {
		return fmt.Errorf("code is required")
	}
	if r.RedirectURI == "" {
		return fmt.Errorf("redirectUri is required")
	}
	return nil
}

func (r *RefreshRequest) Validate() error {
	if r.RefreshToken == "" {
		return fmt.Errorf("refreshToken is required")
	}
	return nil
}
