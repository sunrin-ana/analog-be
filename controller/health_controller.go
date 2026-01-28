package controller

import (
	"context"
	"github.com/NARUBROWN/spine/pkg/httpx"
	"time"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type HealthController struct {
	db     *bun.DB
	logger *zap.Logger
}

func NewHealthController(db *bun.DB, logger *zap.Logger) *HealthController {
	return &HealthController{
		db:     db,
		logger: logger,
	}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func (c *HealthController) Health(ctx context.Context) httpx.Response[HealthResponse] {
	return httpx.Response[HealthResponse]{
		Body: HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}
}

func (c *HealthController) Ready(ctx context.Context) httpx.Response[HealthResponse] {
	services := make(map[string]string)

	dbCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var result int
	err := c.db.NewSelect().ColumnExpr("1").Scan(dbCtx, &result)
	if err != nil {
		c.logger.Error("Database health check failed", zap.Error(err))
		services["database"] = "unhealthy"

		return httpx.Response[HealthResponse]{
			Body: HealthResponse{
				Status:    "degraded",
				Timestamp: time.Now().Format(time.RFC3339),
				Services:  services,
			},
		}
	}
	services["database"] = "healthy"

	return httpx.Response[HealthResponse]{
		Body: HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now().Format(time.RFC3339),
			Services:  services,
		},
	}
}

func (c *HealthController) Live(ctx context.Context) httpx.Response[HealthResponse] {
	return httpx.Response[HealthResponse]{
		Body: HealthResponse{
			Status:    "alive",
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}
}
