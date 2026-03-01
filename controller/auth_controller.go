package controller

import (
	"analog-be/service"
	"context"
	"net/http"

	"github.com/NARUBROWN/spine/pkg/httpx"

	"github.com/NARUBROWN/spine/pkg/query"
)

type AuthController struct {
	oauthService *service.OAuthService
	userService  *service.UserService
}

func NewAuthController(anAccountOAuthService *service.OAuthService, userService *service.UserService) *AuthController {
	return &AuthController{
		oauthService: anAccountOAuthService,
		userService:  userService,
	}
}

// HandleAuthCallback handles the OAuth2 callback after a successful login.
// @Summary      HandleAuthCallback
// @Description  Handle the OAuth2 callback after a successful login.
// @Tags         Auth
// @Produce      json
// @Param        code query string true "Authorization code"
// @Param        state query string true "State"
// @Success      200 {object} dto.AuthResponse
// @Failure      400 "Bad Request"
// @Failure      500 "Internal Server Error"
// @Router       /auth/callback [get]
func (c *AuthController) HandleAuthCallback(ctx context.Context, q query.Values) httpx.Response[string] {
	code := q.Get("code")
	state := q.Get("state")

	if code == "" {
		return httpx.Response[string]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // code is required
			},
		}
	}
	if state == "" {
		return httpx.Response[string]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // state is required
			},
		}
	}

	result, err := c.oauthService.HandleCallback(ctx, code, state)
	if err != nil {
		return httpx.Response[string]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	return httpx.Response[string]{
		Options: httpx.ResponseOptions{
			Status:  http.StatusSeeOther,
			Cookies: result.Cookies,
			Headers: map[string]string{"Location": result.RedirectUri},
		},
	}
}

// RefreshToken refreshes the access token using a valid refresh token.
// @Summary      RefreshToken
// @Description  Refreshes the access token using a valid refresh token.
// @Tags         Auth
// @Success      200 "OK"
// @Failure      401 "Not Authorized"
// @Failure      500 "Internal Server Error"
// @Router       /auth/callback [get]
func (c *AuthController) RefreshToken(ctx context.Context, headers http.Header) {
	headers.Get("")
	// TODO: impl
}

func (c *AuthController) Logout(ctx context.Context) {
	// TODO: impl
}
