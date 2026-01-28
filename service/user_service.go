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
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Get(ctx context.Context, id *entity.ID) (*entity.User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (s *UserService) Create(ctx context.Context, req *dto.UserCreateRequest) (*entity.User, error) {
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

	user, err := s.repository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, id *entity.ID, req *dto.UserUpdateRequest) (*entity.User, error) {
	user, err := s.repository.FindByID(ctx, id)
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

	user, err = s.repository.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id *entity.ID) error {
	_, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	err = s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *UserService) List(ctx context.Context, limit, offset int) (*dto.PaginatedResult[*entity.User], error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := s.repository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return &dto.PaginatedResult[*entity.User]{
		Items:  users,
		Total:  *total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *UserService) Search(ctx context.Context, query string, limit, offset int) (*dto.PaginatedResult[*entity.User], error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := s.repository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return &dto.PaginatedResult[*entity.User]{
		Items:  users,
		Total:  *total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
