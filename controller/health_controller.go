package controller

import (
	"context"
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

func (c *HealthController) Health(ctx context.Context) (*HealthResponse, error) {
	return &HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

func (c *HealthController) Ready(ctx context.Context) (*HealthResponse, error) {
	services := make(map[string]string)

	dbCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var result int
	err := c.db.NewSelect().ColumnExpr("1").Scan(dbCtx, &result)
	if err != nil {
		c.logger.Error("Database health check failed", zap.Error(err))
		services["database"] = "unhealthy"
		return &HealthResponse{
			Status:    "degraded",
			Timestamp: time.Now().Format(time.RFC3339),
			Services:  services,
		}, nil
	}
	services["database"] = "healthy"

	return &HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  services,
	}, nil
}

func (c *HealthController) Live(ctx context.Context) (*HealthResponse, error) {
	return &HealthResponse{
		Status:    "alive",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}
