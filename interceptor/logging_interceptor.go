package interceptor

import (
	"analog-be/pkg"
	"time"

	"github.com/NARUBROWN/spine/core"
	"go.uber.org/zap"
)

type LoggingInterceptor struct {
	logger *zap.Logger
}

func NewLoggingInterceptor() *LoggingInterceptor {
	return &LoggingInterceptor{
		logger: pkg.GetLogger(),
	}
}

func (i *LoggingInterceptor) PreHandle(ctx core.ExecutionContext, meta core.HandlerMeta) error {
	start := time.Now()
	ctx.Set("request_start", start)

	i.logger.Info("Request started",
		zap.String("method", ctx.Method()),
		zap.String("path", ctx.Path()),
		zap.String("controller", meta.ControllerType.Name()),
		zap.String("handler", meta.Method.Name),
		zap.Time("timestamp", start),
	)

	return nil
}

func (i *LoggingInterceptor) PostHandle(core.ExecutionContext, core.HandlerMeta) {}

func (i *LoggingInterceptor) AfterCompletion(ctx core.ExecutionContext, meta core.HandlerMeta, err error) {
	var start time.Time
	if startAny, ok := ctx.Get("request_start"); ok {
		if startTime, ok := startAny.(time.Time); ok {
			start = startTime
		} else {
			start = time.Now()
		}
	} else {
		start = time.Now()
	}

	duration := time.Since(start)

	if err != nil {
		i.logger.Error("Request failed",
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.Path()),
			zap.String("controller", meta.ControllerType.Name()),
			zap.String("handler", meta.Method.Name),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
	} else {
		i.logger.Info("Request completed",
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.Path()),
			zap.String("controller", meta.ControllerType.Name()),
			zap.String("handler", meta.Method.Name),
			zap.Duration("duration", duration),
		)
	}
}
