package service

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/repository"
	"context"
	"fmt"
	"time"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := s.repository.FindByID(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &dto.UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		ProfileImage: user.ProfileImage,
		JoinedAt:     user.JoinedAt.Format(time.RFC3339),
		PartOf:       user.PartOf,
		Generation:   user.Generation,
		Connections:  user.Connections,
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req dto.UserCreateRequest) (*dto.UserResponse, error) {
	now := time.Now()

	user := &entity.User{
		ID:           req.ID,
		Name:         req.Name,
		ProfileImage: req.ProfileImage,
		JoinedAt:     now,
		PartOf:       req.PartOf,
		Generation:   req.Generation,
		Connections:  req.Connections,
	}

	if user.Connections == nil {
		user.Connections = []string{}
	}

	err := s.repository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &dto.UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		ProfileImage: user.ProfileImage,
		JoinedAt:     user.JoinedAt.Format(time.RFC3339),
		PartOf:       user.PartOf,
		Generation:   user.Generation,
		Connections:  user.Connections,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, req dto.UserUpdateRequest) (*dto.UserResponse, error) {
	user, err := s.repository.FindByID(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.ProfileImage != nil {
		user.ProfileImage = *req.ProfileImage
	}
	if req.PartOf != nil {
		user.PartOf = *req.PartOf
	}
	if req.Generation != nil {
		user.Generation = *req.Generation
	}
	if req.Connections != nil {
		user.Connections = *req.Connections
	}

	err = s.repository.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &dto.UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		ProfileImage: user.ProfileImage,
		JoinedAt:     user.JoinedAt.Format(time.RFC3339),
		PartOf:       user.PartOf,
		Generation:   user.Generation,
		Connections:  user.Connections,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	_, err := s.repository.FindByID(ctx, int(id))
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	err = s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) (*dto.UserListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.repository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	total, err := s.repository.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	userResponses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &dto.UserResponse{
			ID:           user.ID,
			Name:         user.Name,
			ProfileImage: user.ProfileImage,
			JoinedAt:     user.JoinedAt.Format(time.RFC3339),
			PartOf:       user.PartOf,
			Generation:   user.Generation,
			Connections:  user.Connections,
		}
	}

	return &dto.UserListResponse{
		Users:  userResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *UserService) SearchUsers(ctx context.Context, query string, limit, offset int) (*dto.UserListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.repository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	total := len(users)

	userResponses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &dto.UserResponse{
			ID:           user.ID,
			Name:         user.Name,
			ProfileImage: user.ProfileImage,
			JoinedAt:     user.JoinedAt.Format(time.RFC3339),
			PartOf:       user.PartOf,
			Generation:   user.Generation,
			Connections:  user.Connections,
		}
	}

	return &dto.UserListResponse{
		Users:  userResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
