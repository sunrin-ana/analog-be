package controller

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/service"
	"context"
	"fmt"
	"strconv"

	"github.com/NARUBROWN/spine/pkg/query"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GetUser(ctx context.Context, q query.Values) (*dto.UserResponse, error) {
	idStr := q.Get("id")
	if idStr == "" {
		return nil, fmt.Errorf("user id is required")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	user, err := c.userService.Get(ctx, (*entity.ID)(&id))
	if err != nil {
		return nil, err
	}

	res := dto.NewUserResponse(user)
	return &res, nil
}

func (c *UserController) CreateUser(ctx context.Context, req dto.UserCreateRequest) (*dto.UserResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	user, err := c.userService.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	res := dto.NewUserResponse(user)
	return &res, nil
}

func (c *UserController) UpdateUser(ctx context.Context, q query.Values, req dto.UserUpdateRequest) (*dto.UserResponse, error) {
	idStr := q.Get("id")
	if idStr == "" {
		return nil, fmt.Errorf("user id is required")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	user, err := c.userService.Update(ctx, (*entity.ID)(&id), req)
	if err != nil {
		return nil, err
	}

	res := dto.NewUserResponse(user)
	return &res, nil
}

func (c *UserController) DeleteUser(ctx context.Context, q query.Values) (interface{}, error) {
	idStr := q.Get("id")
	if idStr == "" {
		return nil, fmt.Errorf("user id is required")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	err = c.userService.Delete(ctx, (*entity.ID)(&id))
	if err != nil {
		return nil, err
	}

	return map[string]string{"message": "user deleted successfully"}, nil
}

func (c *UserController) SearchUser(ctx context.Context, q query.Values) (*dto.UserListResponse, error) {
	searchQuery := q.Get("q")
	if searchQuery == "" {
		return nil, fmt.Errorf("search query is required")
	}

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

	paginatedResult, err := c.userService.Search(ctx, searchQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, len(paginatedResult.Items))
	for i, user := range paginatedResult.Items {
		userResponses[i] = dto.NewUserResponse(user)
	}

	return &dto.UserListResponse{
		Users:  userResponses,
		Total:  paginatedResult.Total,
		Limit:  paginatedResult.Limit,
		Offset: paginatedResult.Offset,
	}, nil
}
