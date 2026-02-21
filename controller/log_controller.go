package controller

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/pkg"
	"analog-be/service"
	"context"
	"net/http"

	"github.com/NARUBROWN/spine/pkg/httperr"
	"github.com/NARUBROWN/spine/pkg/httpx"
	"github.com/NARUBROWN/spine/pkg/path"
	"github.com/NARUBROWN/spine/pkg/query"
	"github.com/NARUBROWN/spine/pkg/spine"
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

// GetListOfLog gets a paginated list of logs.
// @Summary      GetListOfLog
// @Description  Get a paginated list of logs.
// @Tags         Log
// @Produce      json
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} dto.PaginatedResult[dto.LogResponse]
// @Failure		 404 "Not Found"
// @Router       /logs [get]
func (c *LogController) GetListOfLog(ctx context.Context, page query.Pagination) httpx.Response[dto.PaginatedResult[dto.LogResponse]] {
	paginatedResult, err := c.logService.GetList(ctx, page.Size, page.Page)
	if err != nil {
		return httpx.Response[dto.PaginatedResult[dto.LogResponse]]{
			Options: httpx.ResponseOptions{
				Status: http.StatusNotFound, // not found
			},
		}
	}

	logResponses := make([]dto.LogResponse, len(paginatedResult.Items))
	for i, log := range paginatedResult.Items {
		logResponses[i] = dto.NewLogResponse(log)
	}

	return httpx.Response[dto.PaginatedResult[dto.LogResponse]]{
		Body: dto.PaginatedResult[dto.LogResponse]{
			Items:  logResponses,
			Total:  paginatedResult.Total,
			Limit:  paginatedResult.Limit,
			Offset: paginatedResult.Offset,
		},
	}
}

// GetListOfTopicLog gets a paginated list of logs for a specific topic.
// @Summary      GetListOfTopicLog
// @Description  Get a paginated list of logs for a specific topic.
// @Tags         Log
// @Produce      json
// @Param        topicId path int true "Topic ID"
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} dto.PaginatedResult[dto.LogResponse]
// @Failure      404 "Not Found"
// @Router       /logs/topic/list/{topicId} [get]
func (c *LogController) GetListOfTopicLog(ctx context.Context, topicID path.Int, page query.Pagination) httpx.Response[dto.PaginatedResult[dto.LogResponse]] {
	paginatedResult, err := c.logService.GetListByTopicID(ctx, &topicID.Value, page.Size, page.Page)
	if err != nil {
		return httpx.Response[dto.PaginatedResult[dto.LogResponse]]{
			Options: httpx.ResponseOptions{
				Status: http.StatusNotFound, // not found
			},
		}
	}

	logResponses := make([]dto.LogResponse, len(paginatedResult.Items))
	for i, log := range paginatedResult.Items {
		logResponses[i] = dto.NewLogResponse(log)
	}

	return httpx.Response[dto.PaginatedResult[dto.LogResponse]]{
		Body: dto.PaginatedResult[dto.LogResponse]{
			Items:  logResponses,
			Total:  paginatedResult.Total,
			Limit:  paginatedResult.Limit,
			Offset: paginatedResult.Offset,
		},
	}
}

// GetListOfGenerationLog gets a paginated list of logs for a specific generation.
// @Summary      GetListOfGenerationLog
// @Description  Get a paginated list of logs for a specific generation.
// @Tags         Log
// @Produce      json
// @Param        generation path int true "Generation"
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} dto.PaginatedResult[dto.LogResponse]
// @Failure      404 "Not Found"
// @Router       /logs/generation/list/{generation} [get]
func (c *LogController) GetListOfGenerationLog(ctx context.Context, generation path.Int, page query.Pagination) httpx.Response[dto.PaginatedResult[dto.LogResponse]] {
	paginatedResult, err := c.logService.GetListByGeneration(ctx, uint16(generation.Value), page.Size, page.Page)
	if err != nil {
		return httpx.Response[dto.PaginatedResult[dto.LogResponse]]{
			Options: httpx.ResponseOptions{
				Status: http.StatusNotFound, // not found
			},
		}
	}

	logResponses := make([]dto.LogResponse, len(paginatedResult.Items))
	for i, log := range paginatedResult.Items {
		logResponses[i] = dto.NewLogResponse(log)
	}

	return httpx.Response[dto.PaginatedResult[dto.LogResponse]]{
		Body: dto.PaginatedResult[dto.LogResponse]{
			Items:  logResponses,
			Total:  paginatedResult.Total,
			Limit:  paginatedResult.Limit,
			Offset: paginatedResult.Offset,
		},
	}
}

// GetLog gets a single log by its ID.
// @Summary      GetLog
// @Description  Get a single log by its ID.
// @Tags         Log
// @Produce      json
// @Param        id path int true "Log ID"
// @Success      200 {object} dto.LogResponse
// @Failure      404 "Not Found"
// @Router       /logs/{id} [get]
func (c *LogController) GetLog(ctx context.Context, id path.Int) httpx.Response[dto.LogResponse] {
	log, err := c.logService.Get(ctx, &id.Value)
	if err != nil {
		return httpx.Response[dto.LogResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusNotFound, // not found
			},
		}
	}

	res := dto.NewLogResponse(log)
	return httpx.Response[dto.LogResponse]{
		Body: res,
	}
}

// SearchLogs searches for logs by a query string.
// @Summary      SearchLogs
// @Description  Search for logs by a query string.
// @Tags         Log
// @Produce      json
// @Param        q query string true "Search query"
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} dto.PaginatedResult[dto.LogResponse]
// @Failure      404 "Not Found"
// @Router       /logs/search/list [get]
func (c *LogController) SearchLogs(ctx context.Context, q query.Values, page query.Pagination) httpx.Response[dto.PaginatedResult[dto.LogResponse]] {
	searchQuery := q.Get("q")

	paginatedResult, err := c.logService.Search(ctx, searchQuery, page.Size, page.Page)
	if err != nil {
		return httpx.Response[dto.PaginatedResult[dto.LogResponse]]{
			Options: httpx.ResponseOptions{
				Status: http.StatusNotFound, // not found
			},
		}
	}

	logResponses := make([]dto.LogResponse, len(paginatedResult.Items))
	for i, log := range paginatedResult.Items {
		logResponses[i] = dto.NewLogResponse(log)
	}

	return httpx.Response[dto.PaginatedResult[dto.LogResponse]]{
		Body: dto.PaginatedResult[dto.LogResponse]{
			Items:  logResponses,
			Total:  paginatedResult.Total,
			Limit:  paginatedResult.Limit,
			Offset: paginatedResult.Offset,
		},
	}
}

// CreateLog creates a new log.
// @Summary      CreateLog
// @Description  Create a new log.
// @Tags         Log
// @Accept       json
// @Produce      json
// @Param        log body dto.LogCreateRequest true "Log to create"
// @Success      200 {object} dto.LogResponse
// @Failure      400 "Bad Request"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /logs [post]
func (c *LogController) CreateLog(ctx context.Context, req *dto.LogCreateRequest) httpx.Response[dto.LogResponse] {
	if err := pkg.Validate(req); err != nil {
		return httpx.Response[dto.LogResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // validation error
			},
		}
	}

	authorID, ok := pkg.GetUserID(ctx)
	if !ok {
		return httpx.Response[dto.LogResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusUnauthorized, // authentication required
			},
		}
	}

	log, err := c.logService.Create(ctx, req, &authorID)
	if err != nil {
		return httpx.Response[dto.LogResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	res := dto.NewLogResponse(log)
	return httpx.Response[dto.LogResponse]{
		Body: res,
	}
}

// UpdateLog updates an existing log.
// @Summary      UpdateLog
// @Description  Update an existing log.
// @Tags         Log
// @Accept       json
// @Produce      json
// @Param        id path int true "Log ID"
// @Param        log body dto.LogUpdateRequest true "Log data to update"
// @Success      200 {object} dto.LogResponse
// @Failure      401 "Unauthorized"
// @Failure      403 "Forbidden"
// @Failure      500 "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /logs/{id} [put]
func (c *LogController) UpdateLog(ctx context.Context, id path.Int, req *dto.LogUpdateRequest) httpx.Response[dto.LogResponse] {

	userID, ok := pkg.GetUserID(ctx)
	if !ok {
		return httpx.Response[dto.LogResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusUnauthorized, // authentication required
			},
		}
	}

	log, err := c.logService.Get(ctx, &id.Value)
	if err != nil {
		return httpx.Response[dto.LogResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	hasPermission := false
	for _, author := range log.LoggedBy {
		if author.ID == userID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return httpx.Response[dto.LogResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusForbidden, // You don't have permission to update this log
			},
		}
	}

	updatedLog, err := c.logService.Update(ctx, &id.Value, req, &userID)
	if err != nil {
		return httpx.Response[dto.LogResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	res := dto.NewLogResponse(updatedLog)
	return httpx.Response[dto.LogResponse]{
		Body: res,
	}
}

// DeleteLog deletes a log by its ID.
// @Summary      DeleteLog
// @Description  Delete a log by its ID.
// @Tags         Log
// @Param        id path int true "Log ID"
// @Success      204 "No Content"
// @Failure      401 "Unauthorized"
// @Failure      403 "Forbidden"
// @Failure      500 "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /logs/{id} [delete]
func (c *LogController) DeleteLog(ctx context.Context, id path.Int) error {

	userID, ok := pkg.GetUserID(ctx)
	if !ok {
		return httperr.Unauthorized("Authentication required")
	}

	log, err := c.logService.Get(ctx, &id.Value)
	if err != nil {
		return &httperr.HTTPError{
			Status:  500,
			Message: "Internal Server Error",
			Cause:   err,
		}
	}

	hasPermission := false
	for _, author := range log.LoggedBy {
		if author.ID == userID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return &httperr.HTTPError{
			Status:  403,
			Message: "You don't have permission to delete this log",
			Cause:   nil,
		}
	}

	err = c.logService.Delete(ctx, &id.Value)
	if err != nil {
		return &httperr.HTTPError{
			Status:  500,
			Message: "Internal Server Error",
			Cause:   err,
		}
	}

	return nil
}

// CreateComment creates a new comment on a specific log.
// @Summary      CreateComment
// @Description  Create a new comment on a specific log.
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        id path int true "Log ID"
// @Param        comment body dto.CommentCreateRequest true "Comment to create"
// @Success      200 {object} dto.CommentResponse
// @Failure      400 "Bad Request"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /logs/{id}/comments [post]
func (c *LogController) CreateComment(ctx context.Context, id path.Int, req *dto.CommentCreateRequest, spineCtx spine.Ctx) httpx.Response[dto.CommentResponse] {

	if err := pkg.Validate(req); err != nil {
		return httpx.Response[dto.CommentResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // validation error
			},
		}
	}

	v, ok := spineCtx.Get(string(pkg.UserIDKey))
	if !ok {
		return httpx.Response[dto.CommentResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusUnauthorized,
			},
		}
	}

	authorID := v.(entity.ID)

	comment, err := c.commentService.Create(ctx, req, &id.Value, &authorID)
	if err != nil {
		return httpx.Response[dto.CommentResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	res := dto.NewCommentResponse(comment)
	return httpx.Response[dto.CommentResponse]{
		Body: res,
	}
}

// UpdateComment updates an existing comment.
// @Summary      UpdateComment
// @Description  Update an existing comment.
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        id path int true "Log ID"
// @Param        commentId path int true "Comment ID"
// @Param        comment body dto.CommentUpdateRequest true "Comment data to update"
// @Success      200 {object} dto.CommentResponse
// @Failure      400 "Bad Request"
// @Failure      401 "Unauthorized"
// @Failure      403 "Forbidden"
// @Failure      500 "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /logs/{id}/comments/{commentId} [put]
func (c *LogController) UpdateComment(ctx context.Context, id path.Int, commentId path.Int, req *dto.CommentUpdateRequest) httpx.Response[dto.CommentResponse] {

	if err := pkg.Validate(req); err != nil {
		return httpx.Response[dto.CommentResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // validation error
			},
		}
	}

	authorID, ok := pkg.GetUserID(ctx)
	if !ok {
		return httpx.Response[dto.CommentResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusUnauthorized, // authentication required
			},
		}
	}

	comment, err := c.commentService.GetById(ctx, &commentId.Value)
	if err != nil {
		return httpx.Response[dto.CommentResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	if comment.AuthorID != authorID {
		return httpx.Response[dto.CommentResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusForbidden, // forbidden
			},
		}
	}

	if comment.LogID != id.Value {
		return httpx.Response[dto.CommentResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // invalid log id
			},
		}
	}

	comment, err = c.commentService.Update(ctx, &id.Value, req)
	if err != nil {
		return httpx.Response[dto.CommentResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	res := dto.NewCommentResponse(comment)
	return httpx.Response[dto.CommentResponse]{
		Body: res,
	}
}

// DeleteComment deletes a comment by its ID.
// @Summary      DeleteComment
// @Description  Delete a comment by its ID.
// @Tags         Comment
// @Param        id path int true "Log ID"
// @Param        commentId path int true "Comment ID"
// @Success      204 "No Content"
// @Failure      401 "Unauthorized"
// @Failure      403 "Forbidden"
// @Failure      404 "Not Found"
// @Failure      500 "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /logs/{id}/comments/{commentId} [delete]
func (c *LogController) DeleteComment(ctx context.Context, id path.Int, commentId path.Int) error {
	authorID, ok := pkg.GetUserID(ctx)
	if !ok {
		return httperr.Unauthorized("Authentication required")
	}

	comment, err := c.commentService.GetById(ctx, &commentId.Value)
	if err != nil {
		return httperr.NotFound("Comment not found")
	}

	if comment.AuthorID != authorID {
		return &httperr.HTTPError{
			Status:  403,
			Message: "Forbidden",
			Cause:   nil,
		}
	}

	if comment.LogID != id.Value {
		return httperr.BadRequest("Invalid Log ID")
	}

	err = c.commentService.Delete(ctx, &id.Value)
	if err != nil {
		return &httperr.HTTPError{
			Status:  500,
			Message: "Internal Server Error",
			Cause:   err,
		}
	}

	return nil
}

// FindAllCommentByLogID gets a paginated list of all comments for a specific log.
// @Summary      FindAllCommentByLogID
// @Description  Get a paginated list of all comments for a specific log.
// @Tags         Comment
// @Produce      json
// @Param        id path int true "Log ID"
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} dto.PaginatedResult[dto.CommentResponse]
// @Failure      404 "Not Found"
// @Router       /logs/{id}/comments [get]
func (c *LogController) FindAllCommentByLogID(ctx context.Context, page query.Pagination, id path.Int) httpx.Response[dto.PaginatedResult[dto.CommentResponse]] {

	result, err := c.commentService.FindByLogID(ctx, &id.Value, page.Size, page.Page)
	if err != nil {
		return httpx.Response[dto.PaginatedResult[dto.CommentResponse]]{
			Options: httpx.ResponseOptions{
				Status: http.StatusNotFound, // not found
			},
		}
	}

	commentResponses := make([]dto.CommentResponse, len(result.Items))
	for i, item := range result.Items {
		commentResponses[i] = dto.NewCommentResponse(item)
	}

	return httpx.Response[dto.PaginatedResult[dto.CommentResponse]]{
		Body: dto.PaginatedResult[dto.CommentResponse]{
			Items:  commentResponses,
			Total:  result.Total,
			Limit:  result.Limit,
			Offset: result.Offset,
		},
	}
}
