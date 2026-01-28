package controller

import (
	"analog-be/dto"
	"analog-be/pkg"
	"analog-be/service"
	"context"
	"github.com/NARUBROWN/spine/pkg/httperr"
	"github.com/NARUBROWN/spine/pkg/httpx"
	"net/http"

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

func (c *AuthController) InitiateLogin(ctx context.Context, req *dto.LoginInitRequest) httpx.Response[dto.LoginInitResponse] {
	if req.RedirectUri == "" {
		return httpx.Response[dto.LoginInitResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // redirectUri is required
			},
		}
	}

	result, err := c.anAccountOAuthService.InitiateLogin(ctx, req.RedirectUri)
	if err != nil {
		return httpx.Response[dto.LoginInitResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}
	return httpx.Response[dto.LoginInitResponse]{
		Body: *result,
	}
}

func (c *AuthController) InitiateSignup(ctx context.Context, req *dto.SignupInitRequest) httpx.Response[dto.SignupInitResponse] {
	if req.RedirectUri == "" {
		return httpx.Response[dto.SignupInitResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // redirectUri is required
			},
		}
	}

	result, err := c.anAccountOAuthService.InitiateSignup(ctx, req.RedirectUri)
	if err != nil {
		return httpx.Response[dto.SignupInitResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	return httpx.Response[dto.SignupInitResponse]{
		Body: *result,
	}
}

func (c *AuthController) HandleLoginCallback(ctx context.Context, q query.Values) httpx.Response[dto.AuthResponse] {
	code := q.Get("code")
	state := q.Get("state")

	if code == "" {
		return httpx.Response[dto.AuthResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // code is required
			},
		}
	}
	if state == "" {
		return httpx.Response[dto.AuthResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // state is required
			},
		}
	}

	result, err := c.anAccountOAuthService.HandleCallback(ctx, code, state)
	if err != nil {
		return httpx.Response[dto.AuthResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	return httpx.Response[dto.AuthResponse]{
		Body: *result,
	}
}

func (c *AuthController) HandleSignupCallback(ctx context.Context, q query.Values) httpx.Response[dto.AuthResponse] {
	code := q.Get("code")
	state := q.Get("state")

	if code == "" {
		return httpx.Response[dto.AuthResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // code is required
			},
		}
	}
	if state == "" {
		return httpx.Response[dto.AuthResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // state is required
			},
		}
	}

	result, err := c.anAccountOAuthService.HandleCallback(ctx, code, state)
	if err != nil {
		return httpx.Response[dto.AuthResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	return httpx.Response[dto.AuthResponse]{
		Body: *result,
	}
}

func (c *AuthController) Logout(ctx context.Context, req *dto.LogoutRequest) error {
	sessionToken, ok := pkg.GetSessionToken(ctx)
	if !ok || sessionToken == "" {
		sessionToken = req.SessionToken
	}

	if sessionToken == "" {
		return httperr.BadRequest("sessionToken is required")
	}

	err := c.anAccountOAuthService.Logout(ctx, sessionToken)
	if err != nil {
		return &httperr.HTTPError{
			Status:  500,
			Message: "Internal Server Error",
			Cause:   err,
		}
	}

	return nil
}

func (c *AuthController) GetCurrentUser(ctx context.Context, q query.Values) httpx.Response[dto.UserDTO] {
	sessionToken, ok := pkg.GetSessionToken(ctx)
	if !ok || sessionToken == "" {
		return httpx.Response[dto.UserDTO]{
			Options: httpx.ResponseOptions{
				Status: http.StatusUnauthorized, // unauthorized
			},
		}
	}

	user, err := c.anAccountOAuthService.ValidateSession(ctx, sessionToken)
	if err != nil {
		return httpx.Response[dto.UserDTO]{
			Options: httpx.ResponseOptions{
				Status: http.StatusUnauthorized, // unauthorized
			},
		}
	}

	return httpx.Response[dto.UserDTO]{
		Body: dto.UserDTO{
			ID:           user.ID,
			Name:         user.Name,
			ProfileImage: user.ProfileImage,
			PartOf:       user.PartOf,
			Generation:   user.Generation,
			Connections:  user.Connections,
		},
	}
}

func (c *AuthController) RefreshToken(ctx context.Context, req *dto.TokenRefreshRequest) httpx.Response[dto.TokenResponse] {
	if req.RefreshToken == "" {
		return httpx.Response[dto.TokenResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // refresh_token is required
			},
		}
	}

	tokenResp, err := c.anAccountOAuthService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		return httpx.Response[dto.TokenResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // failed to refresh token
			},
		}
	}

	return httpx.Response[dto.TokenResponse]{
		Body: *tokenResp,
	}
}
