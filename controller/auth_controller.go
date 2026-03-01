package controller

import (
	"analog-be/pkg"
	"analog-be/service"
	"context"
	"net/http"

	"github.com/NARUBROWN/spine/pkg/httpx"
	"github.com/NARUBROWN/spine/pkg/spine"
	"github.com/labstack/gommon/log"

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
		log.Error(err)
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
func (c *AuthController) RefreshToken(ctx context.Context, spineCtx spine.Ctx) (*httpx.Response[string], error) {
	v, ok := spineCtx.Get("Cookie")

	if !ok {
		return nil, pkg.NewUnauthorizedError("Missing cookies")
	}

	cookies, ok := v.(map[string]*httpx.Cookie)

	if !ok || cookies == nil || cookies["alog_tkn"] == nil {
		return nil, pkg.NewUnauthorizedError("Missing cookies")
	}

	refreshToken := cookies["refresh_tkn"]

	result, err := c.oauthService.RefreshToken(ctx, refreshToken.Value)

	if err != nil {
		return nil, err
	}

	return &httpx.Response[string]{
		Options: httpx.ResponseOptions{
			Status:  http.StatusOK,
			Cookies: result[:],
		},
	}, nil
}

func (c *AuthController) Logout(ctx context.Context, spineCtx spine.Ctx) httpx.Response[string] {
	v, ok := spineCtx.Get("Cookie")

	if ok {
		cookies, ok := v.(map[string]*httpx.Cookie)

		if ok && cookies == nil && cookies["refresh_tkn"] != nil {
			_ = c.oauthService.Logout(ctx, cookies["refresh_tkn"].Value)
		}
	}

	return httpx.Response[string]{
		Options: httpx.ResponseOptions{
			Status: http.StatusOK,
			Cookies: []httpx.Cookie{
				{
					Name:   "alog_tkn",
					Value:  "",
					Path:   "/",
					MaxAge: 0,
				},
				{
					Name:   "refresh_tkn",
					Value:  "",
					Path:   "/",
					MaxAge: 0,
				},
			},
		},
	}
}
