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
	oauthService service.OAuthService
	userService  service.UserService
}

func NewAuthController(anAccountOAuthService service.OAuthService, userService service.UserService) *AuthController {
	return &AuthController{
		oauthService: anAccountOAuthService,
		userService:  userService,
	}
}

// HandleAuthCallback handles the OAuth2 callback after a successful login.
// @Summary      Handle OAuth2 Callback
// @Description  Handles the OAuth2 callback, sets auth cookies, and redirects to the original state.
// @Tags         Auth
// @Produce      json
// @Param        code   query     string  true  "Authorization code"
// @Param        state  query     string  true  "State (original redirect URI)"
// @Success      303    {string}  string  "Redirecting to original state"
// @Header       303    {string}  Location "Redirect URI"
// @Failure      400    {object}  pkg.AppError "Missing code or state"
// @Failure      500    {object}  pkg.AppError "Internal Server Error"
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
// @Summary      Refresh Access Token
// @Description  Refreshes the access token using the refresh token stored in the `refresh_tkn` cookie.
// @Tags         Auth
// @Produce      json
// @Success      200  {string}  string       "Successfully refreshed tokens"
// @Header       200  {string}  Set-Cookie   "alog_tkn and refresh_tkn cookies"
// @Failure      401  {object}  pkg.AppError "Unauthorized (missing or invalid refresh token)"
// @Failure      500  {object}  pkg.AppError "Internal Server Error"
// @Router       /auth/token [put]
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

// Logout logs out the user by clearing auth cookies.
// @Summary      Logout
// @Description  Logs out the user by clearing the `alog_tkn` and `refresh_tkn` cookies.
// @Tags         Auth
// @Produce      json
// @Success      200  {string}  string      "Successfully logged out"
// @Header       200  {string}  Set-Cookie  "Expired alog_tkn and refresh_tkn cookies"
// @Router       /auth/token [delete]
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
