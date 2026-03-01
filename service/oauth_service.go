package service

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/repository"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/NARUBROWN/spine/pkg/httpx"
)

type OAuthService struct {
	userRepo     *repository.UserRepository
	tokenService *TokenService
	httpClient   *http.Client
}

func NewAnAccountOAuthService(
	userRepo *repository.UserRepository,
) *OAuthService {
	return &OAuthService{
		userRepo: userRepo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *OAuthService) HandleCallback(ctx context.Context, code, state string) (*dto.AuthResponse, error) {
	tokenResp, err := s.exchangeCodeForToken(code, os.Getenv("BASE_URL")+"/api/auth/callback")
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	userInfo, err := s.getUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	user, err := s.findOrCreateUser(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenService.Sign(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}

	tknCookie := httpx.Cookie{
		Name:     "alog_tkn",
		Value:    *token,
		Path:     "/",
		Domain:   os.Getenv("BASE_URL"),
		MaxAge:   60*60*4 - 60, // Token의 Max Age는 4시간
		HttpOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}

	refreshCookie := httpx.Cookie{
		Name:     "refresh_tkn",
		Value:    refreshToken.Token,
		Path:     "/",
		Domain:   os.Getenv("BASE_URL"),
		Expires:  &refreshToken.ExpiresAt, // Token의 Max Age는 20일
		HttpOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}

	redirectUri, _ := url.QueryUnescape(state)

	return &dto.AuthResponse{
		Cookies:     []httpx.Cookie{tknCookie, refreshCookie},
		RedirectUri: redirectUri,
	}, nil
}

func (s *OAuthService) RefreshToken(ctx context.Context, refreshToken string) (*[2]httpx.Cookie, error) {
	result, err := s.tokenService.RefreshToken(ctx, refreshToken)

	if err != nil {
		return nil, err
	}

	cookies := s.bakeAuthCookies(result.Token, result.RefreshToken)

	return &cookies, nil
}

func (s *OAuthService) Logout(ctx context.Context, refreshToken string) error {
	return nil
}

func (s *OAuthService) bakeAuthCookies(token *string, refreshToken *entity.RefreshToken) [2]httpx.Cookie {
	tknCookie := httpx.Cookie{
		Name:     "alog_tkn",
		Value:    *token,
		Path:     "/",
		Domain:   os.Getenv("BASE_URL"),
		MaxAge:   60*60*4 - 60, // Token의 Max Age는 4시간
		HttpOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}

	refreshCookie := httpx.Cookie{
		Name:     "refresh_tkn",
		Value:    refreshToken.Token,
		Path:     "/",
		Domain:   os.Getenv("BASE_URL"),
		Expires:  &refreshToken.ExpiresAt, // Token의 Max Age는 20일
		HttpOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}

	return [2]httpx.Cookie{tknCookie, refreshCookie}
}

func (s *OAuthService) exchangeCodeForToken(code, redirectUri string) (*dto.TokenResponse, error) {
	baseURL := getAnAccountBaseURL()
	clientID := getAnAccountClientID()
	clientSecret := getAnAccountClientSecret()

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectUri)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", baseURL+"/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request failed: %s", string(body))
	}

	var tokenResp dto.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (s *OAuthService) getUserInfo(accessToken string) (*dto.UserInfoResponse, error) {
	baseURL := getAnAccountBaseURL()

	req, err := http.NewRequest("GET", baseURL+"/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userinfo request failed: %s", string(body))
	}

	var userInfo dto.UserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (s *OAuthService) findOrCreateUser(ctx context.Context, userInfo *dto.UserInfoResponse) (*entity.User, error) {
	var userID entity.ID
	fmt.Sscanf(userInfo.Sub, "%d", &userID)

	user, err := s.userRepo.FindByID(ctx, &userID)
	if err == nil {
		return user, nil
	}

	now := time.Now()
	newUser := &entity.User{
		ID:           userID,
		Name:         userInfo.Name,
		ProfileImage: userInfo.Picture,
		JoinedAt:     now,
		PartOf:       "", // ---> | 계정 승인 X
		Generation:   0,  // ---> | 상태로 지정
		Connections:  []string{},
	}

	newUser, err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

func getAnAccountBaseURL() string {
	baseURL := os.Getenv("AN_ACCOUNT_BASE_URL")
	if baseURL == "" {
		baseURL = "https://accounts.ana.st"
	}
	return baseURL
}

func getAnAccountClientID() string {
	clientID := os.Getenv("AN_ACCOUNT_CLIENT_ID")

	return clientID
}

func getAnAccountClientSecret() string {
	clientSecret := os.Getenv("AN_ACCOUNT_CLIENT_SECRET")
	if clientSecret == "" {
		panic("AN_ACCOUNT_CLIENT_SECRET must be set")
	}
	return clientSecret
}
