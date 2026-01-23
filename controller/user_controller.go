package controller

import (
	"analog-be/dto"
	"analog-be/service"
	"context"
	"fmt"
	"github.com/NARUBROWN/spine/pkg/path"
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

func (c *UserController) Get(ctx context.Context, id path.Int) (dto.UserResponse, error) {

	user, err := c.userService.Get(ctx, &id.Value)
	if err != nil {
		return dto.UserResponse{}, err
	}

	res := dto.NewUserResponse(user)
	return res, nil
}

func (c *UserController) Create(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error) {
	if req.Name == "" {
		return dto.UserResponse{}, fmt.Errorf("name is required")
	}

	user, err := c.userService.Create(ctx, req)
	if err != nil {
		return dto.UserResponse{}, err
	}

	res := dto.NewUserResponse(user)
	return res, nil
}

func (c *UserController) Update(ctx context.Context, id path.Int, req dto.UserUpdateRequest) (dto.UserResponse, error) {

	user, err := c.userService.Update(ctx, &id.Value, req)
	if err != nil {
		return dto.UserResponse{}, err
	}

	res := dto.NewUserResponse(user)
	return res, nil
}

func (c *UserController) Delete(ctx context.Context, id path.Int) (interface{}, error) {

	err := c.userService.Delete(ctx, &id.Value)
	if err != nil {
		return nil, err
	}

	return map[string]string{"message": "user deleted successfully"}, nil
}

func (c *UserController) Search(ctx context.Context, q query.Values) (dto.PaginatedResult[dto.UserResponse], error) {
	searchQuery := q.Get("q")
	if searchQuery == "" {
		return dto.PaginatedResult[dto.UserResponse]{}, fmt.Errorf("search query is required")
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
		return dto.PaginatedResult[dto.UserResponse]{}, err
	}

	userResponses := make([]dto.UserResponse, len(paginatedResult.Items))
	for i, user := range paginatedResult.Items {
		userResponses[i] = dto.NewUserResponse(user)
	}

	return dto.PaginatedResult[dto.UserResponse]{
		Items:  userResponses,
		Total:  paginatedResult.Total,
		Limit:  paginatedResult.Limit,
		Offset: paginatedResult.Offset,
	}, nil
}
