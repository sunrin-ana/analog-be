package dto

import (
	"analog-be/entity"

	"github.com/NARUBROWN/spine/pkg/httpx"
)

type LoginInitRequest struct {
	RedirectUri string `json:"redirectUri" binding:"required"` // 인증 진행 이후 redirect 될 URI
}

type SignupInitRequest struct {
	RedirectUri string `json:"redirectUri" binding:"required"`
}

type SignupInitResponse struct {
	AuthorizationUrl string `json:"authorizationUrl"`
	State            string `json:"state"`
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type OAuthCallbackRequest struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

type AuthResponse struct {
	Cookies     []httpx.Cookie
	RedirectUri string
}

type UserDTO struct {
	ID           entity.ID `json:"id"`
	Name         string    `json:"name"`
	ProfileImage string    `json:"profileImage"`
	PartOf       string    `json:"partOf"`
	Generation   uint16    `json:"generation"`
	Connections  []string  `json:"connections"`
}

type LogoutRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

type UserInfoResponse struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	Picture           string `json:"picture,omitempty"`
}
