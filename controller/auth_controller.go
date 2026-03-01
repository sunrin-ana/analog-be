package controller

import (
	"analog-be/dto"
	"analog-be/pkg"
	"analog-be/service"
	"context"
	"net/http"

	"github.com/NARUBROWN/spine/pkg/httperr"
	"github.com/NARUBROWN/spine/pkg/httpx"

	"github.com/NARUBROWN/spine/pkg/query"
)

type AuthController struct {
	anAccountOAuthService service.AnAccountService
	userService           service.UserService
}

func NewAuthController(anAccountOAuthService service.AnAccountService, userService service.UserService) *AuthController {
	return &AuthController{
		anAccountOAuthService: anAccountOAuthService,
		userService:           userService,
	}
}

// InitiateLogin initiates the OAuth2 login flow.
// @Summary      InitiateLogin
// @Description  Initiate the OAuth2 login flow.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginInitRequest true "Login initiation request"
// @Success      200 {object} dto.LoginInitResponse
// @Failure      400 "Bad Request"
// @Failure      500 "Internal Server Error"
// @Router       /auth/login/init [post]
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

// InitiateSignup initiates the OAuth2 signup flow.
// @Summary      InitiateSignup
// @Description  Initiate the OAuth2 signup flow.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.SignupInitRequest true "Signup initiation request"
// @Success      200 {object} dto.SignupInitResponse
// @Failure      400 "Bad Request"
// @Failure      500 "Internal Server Error"
// @Router       /auth/signup/init [post]
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

// HandleLoginCallback handles the OAuth2 callback after a successful login.
// @Summary      HandleLoginCallback
// @Description  Handle the OAuth2 callback after a successful login.
// @Tags         Auth
// @Produce      json
// @Param        code query string true "Authorization code"
// @Param        state query string true "State"
// @Success      200 {object} dto.AuthResponse
// @Failure      400 "Bad Request"
// @Failure      500 "Internal Server Error"
// @Router       /auth/login/callback [get]
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

// HandleSignupCallback handles the OAuth2 callback after a successful signup.
// @Summary      HandleSignupCallback
// @Description  Handle the OAuth2 callback after a successful signup.
// @Tags         Auth
// @Produce      json
// @Param        code query string true "Authorization code"
// @Param        state query string true "State"
// @Success      200 {object} dto.AuthResponse
// @Failure      400 "Bad Request"
// @Failure      500 "Internal Server Error"
// @Router       /auth/signup/callback [get]
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

// Logout logs out the current user.
// @Summary      Logout
// @Description  Log out the current user.
// @Tags         Auth
// @Accept       json
// @Success      204 "No Content"
// @Failure      400 "Bad Request"
// @Failure      500 "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /auth/logout [post]
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

// GetCurrentUser gets the currently authenticated user's information.
// @Summary      GetCurrentUser
// @Description  Get the currently authenticated user's information.
// @Tags         Auth
// @Produce      json
// @Success      200 {object} dto.UserDTO
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /auth/me [get]
func (c *AuthController) GetCurrentUser(ctx context.Context) httpx.Response[dto.UserDTO] {
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

// RefreshToken refreshes the access token using a refresh token.
// @Summary      RefreshToken
// @Description  Refresh the access token using a refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.TokenRefreshRequest true "Token refresh request"
// @Success      200 {object} dto.TokenResponse
// @Failure      400 "Bad Request"
// @Failure      500 "Internal Server Error"
// @Router       /auth/token/refresh [post]
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
