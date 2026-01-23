package controller

import (
	"analog-be/dto"
	"analog-be/pkg"
	"analog-be/service"
	"context"
	"github.com/NARUBROWN/spine/pkg/path"
	"strconv"

	"github.com/NARUBROWN/spine/pkg/query"
)

type LogController struct {
	logService     *service.LogService
	commentService *service.CommentService
}

func NewLogController(logService *service.LogService, commentService *service.CommentService) *LogController {
	return &LogController{
		logService:     logService,
		commentService: commentService,
	}
}

func (c *LogController) GetList(ctx context.Context, q query.Values) (*dto.PaginatedResult[dto.LogResponse], error) {
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	paginatedResult, err := c.logService.GetList(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	logResponses := make([]dto.LogResponse, len(paginatedResult.Items))
	for i, log := range paginatedResult.Items {
		logResponses[i] = dto.NewLogResponse(log)
	}

	return &dto.PaginatedResult[dto.LogResponse]{
		Items:  logResponses,
		Total:  paginatedResult.Total,
		Limit:  paginatedResult.Limit,
		Offset: paginatedResult.Offset,
	}, nil
}

func (c *LogController) GetLog(ctx context.Context, id path.Int) (*dto.LogResponse, error) {
	log, err := c.logService.Get(ctx, &id.Value)
	if err != nil {
		return nil, err
	}

	res := dto.NewLogResponse(log)
	return &res, nil
}

func (c *LogController) SearchLogs(ctx context.Context, q query.Values) (*dto.PaginatedResult[dto.LogResponse], error) {
	searchQuery := q.Get("q")
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	paginatedResult, err := c.logService.Search(ctx, searchQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	logResponses := make([]dto.LogResponse, len(paginatedResult.Items))
	for i, log := range paginatedResult.Items {
		logResponses[i] = dto.NewLogResponse(log)
	}

	return &dto.PaginatedResult[dto.LogResponse]{
		Items:  logResponses,
		Total:  paginatedResult.Total,
		Limit:  paginatedResult.Limit,
		Offset: paginatedResult.Offset,
	}, nil
}

func (c *LogController) CreateLog(ctx context.Context, req dto.LogCreateRequest) (*dto.LogResponse, error) {
	if err := pkg.Validate(&req); err != nil {
		return nil, err
	}

	authorID, ok := pkg.GetUserID(ctx)
	if !ok {
		return nil, pkg.NewUnauthorizedError("Authentication required")
	}

	log, err := c.logService.Create(ctx, req, &authorID)
	if err != nil {
		return nil, err
	}

	res := dto.NewLogResponse(log)
	return &res, nil
}

func (c *LogController) UpdateLog(ctx context.Context, id path.Int, req dto.LogUpdateRequest) (*dto.LogResponse, error) {

	userID, ok := pkg.GetUserID(ctx)
	if !ok {
		return nil, pkg.NewUnauthorizedError("Authentication required")
	}

	log, err := c.logService.Get(ctx, &id.Value)
	if err != nil {
		return nil, err
	}

	hasPermission := false
	for _, author := range log.LoggedBy {
		if author.ID == userID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return nil, pkg.NewForbiddenError("You don't have permission to update this log")
	}

	updatedLog, err := c.logService.Update(ctx, &id.Value, req, &userID)
	if err != nil {
		return nil, err
	}

	res := dto.NewLogResponse(updatedLog)
	return &res, nil
}

func (c *LogController) DeleteLog(ctx context.Context, id path.Int) (interface{}, error) {

	userID, ok := pkg.GetUserID(ctx)
	if !ok {
		return nil, pkg.NewUnauthorizedError("Authentication required")
	}

	log, err := c.logService.Get(ctx, &id.Value)
	if err != nil {
		return nil, err
	}

	hasPermission := false
	for _, author := range log.LoggedBy {
		if author.ID == userID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return nil, pkg.NewForbiddenError("You don't have permission to delete this log")
	}

	err = c.logService.Delete(ctx, &id.Value)
	if err != nil {
		return nil, err
	}

	return map[string]string{"message": "log deleted successfully"}, nil
}

func (c *LogController) CreateComment(ctx context.Context, id path.Int, req dto.CommentCreateRequest) (*dto.CommentResponse, error) {

	if err := pkg.Validate(&req); err != nil {
		return nil, err
	}

	authorID, ok := pkg.GetUserID(ctx)
	if !ok {
		return nil, pkg.NewUnauthorizedError("Authentication required")
	}

	comment, err := c.commentService.Create(ctx, req, &id.Value, &authorID)
	if err != nil {
		return nil, err
	}

	res := dto.NewCommentResponse(comment)
	return &res, nil
}

func (c *LogController) UpdateComment(ctx context.Context, id path.Int, req dto.CommentUpdateRequest) (*dto.CommentResponse, error) {

	if err := pkg.Validate(&req); err != nil {
		return nil, err
	}

	comment, err := c.commentService.Update(ctx, &id.Value, req)
	if err != nil {
		return nil, err
	}

	res := dto.NewCommentResponse(comment)
	return &res, nil
}

func (c *LogController) DeleteComment(ctx context.Context, id path.Int) (interface{}, error) {

	err := c.commentService.Delete(ctx, &id.Value)
	if err != nil {
		return nil, err
	}

	return map[string]string{"message": "comment deleted successfully"}, nil
}

func (c *LogController) FindAllCommentByLogID(ctx context.Context, q query.Values, id path.Int) ([]dto.CommentResponse, error) {
	limit := 20
	if limitStr := q.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := q.Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	result, err := c.commentService.FindByLogID(ctx, &id.Value, limit, offset)
	if err != nil {
		return nil, err
	}

	commentResponses := make([]dto.CommentResponse, len(result.Items))
	for i, item := range result.Items {
		commentResponses[i] = dto.NewCommentResponse(item)
	}

	return commentResponses, nil
}
