package interceptor

import (
	"analog-be/pkg"
	"analog-be/repository"
	"strings"

	"github.com/NARUBROWN/spine/core"
	"go.uber.org/zap"
)

type AuthInterceptor struct {
	sessionRepo repository.SessionRepository
	logger      *zap.Logger
}

func NewAuthInterceptor(sessionRepo repository.SessionRepository, logger *zap.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

func (i *AuthInterceptor) PreHandle(ctx core.ExecutionContext, _ core.HandlerMeta) error {
	authHeader := ctx.Header("Authorization")
	if authHeader == "" {
		return pkg.NewUnauthorizedError("Missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return pkg.NewUnauthorizedError("Invalid authorization header format")
	}

	sessionToken := parts[1]

	session, err := i.sessionRepo.FindByToken(ctx.Context(), sessionToken)
	if err != nil {
		i.logger.Debug("Invalid session token", zap.Error(err))
		return pkg.NewUnauthorizedError("Invalid or expired session")
	}

	ctx.Set(string(pkg.UserIDKey), session.UserID)
	ctx.Set(string(pkg.SessionTokenKey), sessionToken)

	return nil
}

func (i *AuthInterceptor) PostHandle(core.ExecutionContext, core.HandlerMeta) {}

func (i *AuthInterceptor) AfterCompletion(core.ExecutionContext, core.HandlerMeta, error) {}

func (i *AuthInterceptor) OptionalPreHandle(ctx core.ExecutionContext, _ core.HandlerMeta) error {
	authHeader := ctx.Header("Authorization")
	if authHeader == "" {
		return nil
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil
	}

	sessionToken := parts[1]

	session, err := i.sessionRepo.FindByToken(ctx.Context(), sessionToken)
	if err != nil {
		return nil
	}

	ctx.Set(string(pkg.UserIDKey), session.UserID)
	ctx.Set(string(pkg.SessionTokenKey), sessionToken)

	return nil
}
