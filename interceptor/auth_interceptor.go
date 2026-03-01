package interceptor

import (
	"analog-be/entity"
	"analog-be/pkg"
	"analog-be/service"
	"fmt"

	"github.com/NARUBROWN/spine/core"
	"github.com/NARUBROWN/spine/pkg/httperr"
	"github.com/NARUBROWN/spine/pkg/httpx"
	"go.uber.org/zap"
)

type AuthInterceptor struct {
	tokenService service.TokenService
	logger       *zap.Logger
}

func NewAuthInterceptor(tokenService service.TokenService, logger *zap.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		tokenService: tokenService,
		logger:       logger,
	}
}

func (i *AuthInterceptor) PreHandle(ctx core.ExecutionContext, _ core.HandlerMeta) error {
	v, ok := ctx.Get("Cookie")

	if !ok {
		return pkg.NewUnauthorizedError("Missing cookies")
	}

	cookies, ok := v.(map[string]*httpx.Cookie)

	if !ok || cookies == nil || cookies["alog_tkn"] == nil {
		return pkg.NewUnauthorizedError("Missing cookies")
	}

	token := cookies["alog_tkn"]

	claims, err := i.tokenService.Verify(token.Value)
	if err != nil {
		i.logger.Debug("Invalid token", zap.Error(err))
		return pkg.NewUnauthorizedError("Invalid or expired session")
	}

	var id entity.ID
	_, err = fmt.Sscanf(claims.ID, "%d", &id)
	if err != nil {
		return httperr.Unauthorized("Invalid token")
	}

	ctx.Set(string(pkg.UserClaims), claims)
	ctx.Set(string(pkg.UserID), id)

	return nil
}

func (i *AuthInterceptor) PostHandle(core.ExecutionContext, core.HandlerMeta) {}

func (i *AuthInterceptor) AfterCompletion(core.ExecutionContext, core.HandlerMeta, error) {}
