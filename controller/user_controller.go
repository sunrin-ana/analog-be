package controller

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/pkg"
	"analog-be/service"
	"context"
	"github.com/NARUBROWN/spine/pkg/httperr"
	"github.com/NARUBROWN/spine/pkg/httpx"
	"github.com/NARUBROWN/spine/pkg/path"
	"github.com/NARUBROWN/spine/pkg/query"
	"github.com/NARUBROWN/spine/pkg/spine"
	"net/http"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Get gets a single user by their ID.
// @Summary      GetUserByID
// @Description  Get a single user by their ID.
// @Tags         User
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} dto.UserResponse
// @Failure      404 "Not Found"
// @Security     ApiKeyAuth
// @Router       /users/{id} [get]
func (c *UserController) Get(ctx context.Context, id path.Int) httpx.Response[dto.UserResponse] {
	user, err := c.userService.Get(ctx, &id.Value)
	if err != nil {
		// TODO: err 처리 로직 개선 필요(404 or 500), 다른 controller 도 동일
		return httpx.Response[dto.UserResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusNotFound, // user not found
			},
		}
	}

	res := dto.NewUserResponse(user)
	return httpx.Response[dto.UserResponse]{
		Body: res,
	}
}

// Create creates a new user.
// @Summary      CreateUser
// @Description  Create a new user.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body dto.UserCreateRequest true "User to create"
// @Success      201 {object} dto.UserResponse
// @Failure      400 "Bad Request"
// @Failure      500 "Internal Server Error"
// @Router       /users [post]
func (c *UserController) Create(ctx context.Context, req *dto.UserCreateRequest) httpx.Response[dto.UserResponse] {
	if req.Name == "" {
		return httpx.Response[dto.UserResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // name is required
			},
		}
	}

	user, err := c.userService.Create(ctx, req)
	if err != nil {
		return httpx.Response[dto.UserResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	res := dto.NewUserResponse(user)
	return httpx.Response[dto.UserResponse]{
		Body: res,
		Options: httpx.ResponseOptions{
			Status: http.StatusCreated,
		},
	}
}

// Update updates the current user's information.
// @Summary      UpdateUser
// @Description  Update the current user's information.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body dto.UserUpdateRequest true "User data to update"
// @Success      200 {object} dto.UserResponse
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /users [put]
func (c *UserController) Update(ctx context.Context, req *dto.UserUpdateRequest, spineCtx spine.Ctx) httpx.Response[dto.UserResponse] {
	v, ok := spineCtx.Get(string(pkg.UserIDKey))
	if !ok {
		return httpx.Response[dto.UserResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusUnauthorized,
			},
		}
	}

	id := v.(entity.ID)

	user, err := c.userService.Update(ctx, &id, req)
	if err != nil {
		return httpx.Response[dto.UserResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	res := dto.NewUserResponse(user)
	return httpx.Response[dto.UserResponse]{
		Body: res,
	}
}

// Delete deletes the current user.
// @Summary      DeleteUser
// @Description  Delete the current user.
// @Tags         User
// @Success      204 "No Content"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /users [delete]
func (c *UserController) Delete(ctx context.Context, spineCtx spine.Ctx) error {

	v, ok := spineCtx.Get(string(pkg.UserIDKey))
	if !ok {
		return &httperr.HTTPError{
			Status:  401,
			Message: "Authentication required",
			Cause:   nil,
		}
	}

	id := v.(entity.ID)

	err := c.userService.Delete(ctx, &id)
	if err != nil {
		return &httperr.HTTPError{
			Status:  500,
			Message: "Internal Server Error",
			Cause:   err,
		}
	}

	return nil
}

// Search searches for users by a query string.
// @Summary      SearchUsers
// @Description  Search for users by a query string.
// @Tags         User
// @Produce      json
// @Param        q query string true "Search query"
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} dto.PaginatedResult[dto.UserResponse]
// @Failure      400 "Bad Request"
// @Failure      404 "Not Found"
// @Router       /users/search/list [get]
func (c *UserController) Search(ctx context.Context, q query.Values, page query.Pagination) httpx.Response[dto.PaginatedResult[dto.UserResponse]] {
	searchQuery := q.Get("q")
	if searchQuery == "" {
		return httpx.Response[dto.PaginatedResult[dto.UserResponse]]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // search query is required
			},
		}
	}

	paginatedResult, err := c.userService.Search(ctx, searchQuery, page.Size, page.Page)
	if err != nil {
		return httpx.Response[dto.PaginatedResult[dto.UserResponse]]{
			Options: httpx.ResponseOptions{
				Status: http.StatusNotFound, // not found
			},
		}
	}

	userResponses := make([]dto.UserResponse, len(paginatedResult.Items))
	for i, user := range paginatedResult.Items {
		userResponses[i] = dto.NewUserResponse(user)
	}

	return httpx.Response[dto.PaginatedResult[dto.UserResponse]]{
		Body: dto.PaginatedResult[dto.UserResponse]{
			Items:  userResponses,
			Total:  paginatedResult.Total,
			Limit:  paginatedResult.Limit,
			Offset: paginatedResult.Offset,
		},
	}
}

// GetMe gets the currently authenticated user's information.
// @Summary      GetCurrentUser
// @Description  Get the currently authenticated user's information.
// @Tags         User
// @Produce      json
// @Success      200 {object} dto.UserResponse
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /users/me [get]
func (c *UserController) GetMe(ctx context.Context, spineCtx spine.Ctx) httpx.Response[dto.UserResponse] {
	v, ok := spineCtx.Get(string(pkg.UserIDKey))
	if !ok {
		return httpx.Response[dto.UserResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusUnauthorized,
			},
		}
	}

	id := v.(entity.ID)

	user, err := c.userService.Get(ctx, &id)
	if err != nil {
		return httpx.Response[dto.UserResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	res := dto.NewUserResponse(user)
	return httpx.Response[dto.UserResponse]{
		Body: res,
	}
}
