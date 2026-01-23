package service

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/repository"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AnAccountService struct {
	stateRepo   *repository.OAuthStateRepository
	sessionRepo *repository.SessionRepository
	userRepo    *repository.UserRepository
	httpClient  *http.Client
}

func NewAnAccountOAuthService(
	stateRepo *repository.OAuthStateRepository,
	sessionRepo *repository.SessionRepository,
	userRepo *repository.UserRepository,
) *AnAccountService {
	return &AnAccountService{
		stateRepo:   stateRepo,
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *AnAccountService) InitiateLogin(ctx context.Context, redirectUri string) (*dto.LoginInitResponse, error) {
	return s.initiateOAuth(ctx, redirectUri, false)
}

func (s *AnAccountService) InitiateSignup(ctx context.Context, redirectUri string) (*dto.SignupInitResponse, error) {
	resp, err := s.initiateOAuth(ctx, redirectUri, true)
	if err != nil {
		return nil, err
	}
	return &dto.SignupInitResponse{
		AuthorizationUrl: resp.AuthorizationUrl,
		State:            resp.State,
	}, nil
}

func (s *AnAccountService) initiateOAuth(ctx context.Context, redirectUri string, isSignup bool) (*dto.LoginInitResponse, error) {
	state := uuid.New().String()

	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}

	codeChallenge := generateCodeChallenge(codeVerifier)

	oauthState := &entity.OAuthState{
		State:        state,
		CodeVerifier: codeVerifier,
		RedirectUri:  redirectUri,
		IsSignup:     isSignup,
		ExpiresAt:    time.Now().UTC().Add(10 * time.Minute),
		CreatedAt:    time.Now().UTC(),
	}

	err = s.stateRepo.Create(ctx, oauthState)
	if err != nil {
		return nil, fmt.Errorf("failed to save OAuth state: %w", err)
	}

	authUrl := s.buildAuthorizationURL(state, codeChallenge, redirectUri)

	return &dto.LoginInitResponse{
		AuthorizationUrl: authUrl,
		State:            state,
	}, nil
}

func (s *AnAccountService) HandleCallback(ctx context.Context, code, state string) (*dto.AuthResponse, error) {
	oauthState, err := s.stateRepo.FindByState(ctx, state)
	if err != nil {
		return nil, fmt.Errorf("invalid state: %w", err)
	}

	if time.Now().UTC().After(oauthState.ExpiresAt) {
		return nil, fmt.Errorf("state expired")
	}

	defer s.stateRepo.Delete(ctx, state)

	tokenResp, err := s.exchangeCodeForToken(code, oauthState.CodeVerifier, oauthState.RedirectUri)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	userInfo, err := s.getUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	user, err := s.findOrCreateUser(ctx, userInfo, oauthState.IsSignup)
	if err != nil {
		return nil, err
	}

	session, err := s.createSession(ctx, &user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &dto.AuthResponse{
		SessionToken: session.SessionToken,
		User: &dto.UserDTO{
			ID:           user.ID,
			Name:         user.Name,
			ProfileImage: user.ProfileImage,
			PartOf:       user.PartOf,
			Generation:   user.Generation,
			Connections:  user.Connections,
		},
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func (s *AnAccountService) Logout(ctx context.Context, sessionToken string) error {
	return s.sessionRepo.Delete(ctx, sessionToken)
}

func (s *AnAccountService) ValidateSession(ctx context.Context, sessionToken string) (*entity.User, error) {
	session, err := s.sessionRepo.FindByToken(ctx, sessionToken)
	if err != nil {
		return nil, fmt.Errorf("invalid session")
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		s.sessionRepo.Delete(ctx, sessionToken)
		return nil, fmt.Errorf("session expired")
	}

	if session.User == nil {
		return nil, fmt.Errorf("user not found")
	}

	return session.User, nil
}

func (s *AnAccountService) RefreshAccessToken(refreshToken string) (*dto.TokenResponse, error) {
	baseURL := getAnAccountBaseURL()
	clientID := getAnAccountClientID()
	clientSecret := getAnAccountClientSecret()

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", baseURL+"/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("refresh token failed: %s", string(body))
	}

	var tokenResp dto.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

func (s *AnAccountService) buildAuthorizationURL(state, codeChallenge, redirectUri string) string {
	baseURL := getAnAccountBaseURL()
	clientID := getAnAccountClientID()

	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectUri)
	params.Set("scope", "openid profile email")
	params.Set("state", state)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")

	return fmt.Sprintf("%s/oauth2/authorize?%s", baseURL, params.Encode())
}

func (s *AnAccountService) exchangeCodeForToken(code, codeVerifier, redirectUri string) (*dto.TokenResponse, error) {
	baseURL := getAnAccountBaseURL()
	clientID := getAnAccountClientID()
	clientSecret := getAnAccountClientSecret()

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectUri)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code_verifier", codeVerifier)

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

func (s *AnAccountService) getUserInfo(accessToken string) (*dto.UserInfoResponse, error) {
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

func (s *AnAccountService) findOrCreateUser(ctx context.Context, userInfo *dto.UserInfoResponse, isSignup bool) (*entity.User, error) {
	var userID entity.ID
	fmt.Sscanf(userInfo.Sub, "%d", &userID)

	user, err := s.userRepo.FindByID(ctx, &userID)
	if err == nil {
		if isSignup {
			return nil, fmt.Errorf("user already exists, please login instead")
		}
		return user, nil
	}

	if !isSignup {
		return nil, fmt.Errorf("user not found, please sign up first")
	}

	now := time.Now()
	newUser := &entity.User{
		ID:           userID,
		Name:         userInfo.Name,
		ProfileImage: userInfo.Picture,
		JoinedAt:     now,
		PartOf:       "", // TODO: impl
		Generation:   0,  // TODO: impl
		Connections:  []string{},
	}

	newUser, err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

func (s *AnAccountService) createSession(ctx context.Context, userID *entity.ID) (*entity.Session, error) {
	sessionToken, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	session := &entity.Session{
		SessionToken: sessionToken,
		UserID:       *userID,
		ExpiresAt:    time.Now().UTC().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt:    time.Now().UTC(),
	}

	err = s.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func generateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
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
