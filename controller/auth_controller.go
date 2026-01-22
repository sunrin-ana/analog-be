package controller

import (
	"analog-be/dto"
	"analog-be/pkg"
	"analog-be/service"
	"context"
	"fmt"

	"github.com/NARUBROWN/spine/pkg/query"
)

type AuthController struct {
	anAccountOAuthService *service.AnAccountService
	userService           *service.UserService
}

func NewAuthController(anAccountOAuthService *service.AnAccountService, userService *service.UserService) *AuthController {
	return &AuthController{
		anAccountOAuthService: anAccountOAuthService,
		userService:           userService,
	}
}

func (c *AuthController) InitiateLogin(ctx context.Context, req dto.LoginInitRequest) (*dto.LoginInitResponse, error) {
	if req.RedirectUri == "" {
		return nil, fmt.Errorf("redirectUri is required")
	}

	return c.anAccountOAuthService.InitiateLogin(ctx, req.RedirectUri)
}

func (c *AuthController) InitiateSignup(ctx context.Context, req dto.SignupInitRequest) (*dto.SignupInitResponse, error) {
	if req.RedirectUri == "" {
		return nil, fmt.Errorf("redirectUri is required")
	}

	return c.anAccountOAuthService.InitiateSignup(ctx, req.RedirectUri)
}

func (c *AuthController) HandleLoginCallback(ctx context.Context, q query.Values) (*dto.AuthResponse, error) {
	code := q.Get("code")
	state := q.Get("state")

	if code == "" {
		return nil, fmt.Errorf("code is required")
	}
	if state == "" {
		return nil, fmt.Errorf("state is required")
	}

	return c.anAccountOAuthService.HandleCallback(ctx, code, state)
}

func (c *AuthController) HandleSignupCallback(ctx context.Context, q query.Values) (*dto.AuthResponse, error) {
	code := q.Get("code")
	state := q.Get("state")

	if code == "" {
		return nil, fmt.Errorf("code is required")
	}
	if state == "" {
		return nil, fmt.Errorf("state is required")
	}

	return c.anAccountOAuthService.HandleCallback(ctx, code, state)
}

func (c *AuthController) Logout(ctx context.Context, req dto.LogoutRequest) (interface{}, error) {
	sessionToken, ok := pkg.GetSessionToken(ctx)
	if !ok || sessionToken == "" {
		sessionToken = req.SessionToken
	}

	if sessionToken == "" {
		return nil, fmt.Errorf("sessionToken is required")
	}

	err := c.anAccountOAuthService.Logout(ctx, sessionToken)
	if err != nil {
		return nil, err
	}

	return map[string]string{"message": "logged out successfully"}, nil
}

func (c *AuthController) GetCurrentUser(ctx context.Context, q query.Values) (*dto.UserDTO, error) {
	sessionToken, ok := pkg.GetSessionToken(ctx)
	if !ok || sessionToken == "" {
		return nil, fmt.Errorf("unauthorized")
	}

	user, err := c.anAccountOAuthService.ValidateSession(ctx, sessionToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	return &dto.UserDTO{
		ID:           user.ID,
		Name:         user.Name,
		ProfileImage: user.ProfileImage,
		PartOf:       user.PartOf,
		Generation:   user.Generation,
		Connections:  user.Connections,
	}, nil
}

func (c *AuthController) RefreshToken(ctx context.Context, req dto.TokenRefreshRequest) (*dto.TokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, fmt.Errorf("refresh_token is required")
	}

	tokenResp, err := c.anAccountOAuthService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return tokenResp, nil
}
