package controller

import (
	"analog-be/dto"
	"analog-be/pkg"
	"analog-be/service"
	"context"
	"fmt"
	"strconv"

	"github.com/NARUBROWN/spine/pkg/query"
)

type LogController struct {
	svc *service.LogService
}

func NewLogController(svc *service.LogService) *LogController {
	return &LogController{
		svc: svc,
	}
}

func (c *LogController) GetListOfLog(ctx context.Context, q query.Values) (*dto.LogListResponse, error) {
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return c.svc.GetListOfLog(ctx, limit, offset)
}

func (c *LogController) GetLog(ctx context.Context, q query.Values) (*dto.LogResponse, error) {
	id := q.Get("id")
	if id == "" {
		return nil, fmt.Errorf("log id is required")
	}

	return c.svc.GetLog(ctx, id)
}

func (c *LogController) SearchLogs(ctx context.Context, q query.Values) (*dto.LogListResponse, error) {
	searchQuery := q.Get("q")
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return c.svc.SearchLogs(ctx, searchQuery, limit, offset)
}

func (c *LogController) CreateLog(ctx context.Context, req dto.LogCreateRequest) (*dto.LogResponse, error) {
	if err := pkg.Validate(&req); err != nil {
		return nil, err
	}

	return c.svc.CreateLog(ctx, req)
}

func (c *LogController) UpdateLog(ctx context.Context, q query.Values, req dto.LogUpdateRequest) (*dto.LogResponse, error) {
	id := q.Get("id")
	if id == "" {
		return nil, fmt.Errorf("log id is required")
	}

	userID, ok := pkg.GetUserID(ctx)
	if !ok {
		return nil, pkg.NewUnauthorizedError("Authentication required")
	}

	log, err := c.svc.GetLog(ctx, id)
	if err != nil {
		return nil, err
	}

	hasPermission := false
	for _, authorID := range log.LoggedBy {
		if authorID == userID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return nil, pkg.NewForbiddenError("You don't have permission to update this log")
	}

	return c.svc.UpdateLog(ctx, id, req)
}

func (c *LogController) DeleteLog(ctx context.Context, q query.Values) (interface{}, error) {
	id := q.Get("id")
	if id == "" {
		return nil, fmt.Errorf("log id is required")
	}

	userID, ok := pkg.GetUserID(ctx)
	if !ok {
		return nil, pkg.NewUnauthorizedError("Authentication required")
	}

	log, err := c.svc.GetLog(ctx, id)
	if err != nil {
		return nil, err
	}

	hasPermission := false
	for _, authorID := range log.LoggedBy {
		if authorID == userID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return nil, pkg.NewForbiddenError("You don't have permission to delete this log")
	}

	err = c.svc.DeleteLog(ctx, id)
	if err != nil {
		return nil, err
	}

	return map[string]string{"message": "log deleted successfully"}, nil
}

func (c *LogController) CreateComment(ctx context.Context, q query.Values, req dto.CommentCreateRequest) (*dto.CommentResponse, error) {
	logID := q.Get("id")
	if logID == "" {
		return nil, fmt.Errorf("log id is required")
	}

	if err := pkg.Validate(&req); err != nil {
		return nil, err
	}

	return c.svc.CreateComment(ctx, logID, req)
}

func (c *LogController) UpdateComment(ctx context.Context, q query.Values, req dto.CommentUpdateRequest) (*dto.CommentResponse, error) {
	commentID := q.Get("commentId")
	if commentID == "" {
		return nil, fmt.Errorf("comment id is required")
	}

	if err := pkg.Validate(&req); err != nil {
		return nil, err
	}

	return c.svc.UpdateComment(ctx, commentID, req)
}

func (c *LogController) DeleteComment(ctx context.Context, q query.Values) (interface{}, error) {
	commentID := q.Get("commentId")
	if commentID == "" {
		return nil, fmt.Errorf("comment id is required")
	}

	err := c.svc.DeleteComment(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return map[string]string{"message": "comment deleted successfully"}, nil
}
