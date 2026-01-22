package interceptor

import (
	"analog-be/pkg"
	"database/sql"
	"errors"

	"github.com/NARUBROWN/spine/core"
	"go.uber.org/zap"
)

type ErrorInterceptor struct {
	logger *zap.Logger
}

func NewErrorInterceptor() *ErrorInterceptor {
	return &ErrorInterceptor{
		logger: pkg.GetLogger(),
	}
}

func (i *ErrorInterceptor) PreHandle(core.ExecutionContext, core.HandlerMeta) error {
	return nil
}

func (i *ErrorInterceptor) PostHandle(core.ExecutionContext, core.HandlerMeta) {}

func (i *ErrorInterceptor) AfterCompletion(ctx core.ExecutionContext, meta core.HandlerMeta, err error) {
	if err == nil {
		return
	}

	var appErr *pkg.AppError
	if errors.As(err, &appErr) {
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		return
	}

	i.logger.Error("Unexpected error",
		zap.Error(err),
		zap.String("method", ctx.Method()),
		zap.String("path", ctx.Path()),
		zap.String("controller", meta.ControllerType.Name()),
		zap.String("handler", meta.Method.Name),
	)
}
