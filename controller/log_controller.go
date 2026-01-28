package controller

import (
	"analog-be/dto"
	"analog-be/pkg"
	"analog-be/service"
	"context"
	"github.com/NARUBROWN/spine/pkg/httperr"
	"github.com/NARUBROWN/spine/pkg/httpx"
	"github.com/NARUBROWN/spine/pkg/path"
	"github.com/NARUBROWN/spine/pkg/query"
	"net/http"
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

func (c *LogController) GetListOfLog(ctx context.Context, q query.Values, page query.Pagination) httpx.Response[dto.PaginatedResult[dto.LogResponse]] {
	paginatedResult, err := c.logService.GetList(ctx, page.Size, page.Page)
	if err != nil {
		return httpx.Response[dto.PaginatedResult[dto.LogResponse]]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
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

func (c *LogController) GetLog(ctx context.Context, id path.Int) httpx.Response[dto.LogResponse] {
	log, err := c.logService.Get(ctx, &id.Value)
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

func (c *LogController) SearchLogs(ctx context.Context, q query.Values, page query.Pagination) httpx.Response[dto.PaginatedResult[dto.LogResponse]] {
	searchQuery := q.Get("q")

	paginatedResult, err := c.logService.Search(ctx, searchQuery, page.Size, page.Page)
	if err != nil {
		return httpx.Response[dto.PaginatedResult[dto.LogResponse]]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
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

func (c *LogController) CreateComment(ctx context.Context, id path.Int, commentId path.Int, req *dto.CommentCreateRequest) httpx.Response[dto.CommentResponse] {

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
				Status: http.StatusNotFound, // comment not found
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

	comment, err = c.commentService.Create(ctx, req, &commentId.Value, &authorID)
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

func (c *LogController) DeleteComment(ctx context.Context, id path.Int) error {

	err := c.commentService.Delete(ctx, &id.Value)
	if err != nil {
		return &httperr.HTTPError{
			Status:  500,
			Message: "Internal Server Error",
			Cause:   err,
		}
	}

	return nil
}

func (c *LogController) FindAllCommentByLogID(ctx context.Context, page query.Pagination, id path.Int) httpx.Response[dto.PaginatedResult[dto.CommentResponse]] {

	result, err := c.commentService.FindByLogID(ctx, &id.Value, page.Size, page.Page)
	if err != nil {
		return httpx.Response[dto.PaginatedResult[dto.CommentResponse]]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
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
