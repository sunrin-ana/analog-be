package service

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/repository"
	"context"
)

type CommentService struct {
	commentRepository *repository.CommentRepository
	logRepository     *repository.LogRepository
}

func NewCommentService(commentRepository *repository.CommentRepository, logRepository *repository.LogRepository) *CommentService {
	return &CommentService{commentRepository: commentRepository, logRepository: logRepository}
}

func (s *CommentService) Create(ctx context.Context, req dto.CommentCreateRequest, logID *entity.ID, authorID *entity.ID) (*entity.Comment, error) {
	_, err := s.logRepository.FindByID(ctx, logID)
	if err != nil {
		return nil, err
	}

	comment := &entity.Comment{
		LogID:    *logID,
		AuthorID: *authorID,
		Content:  req.Content,
	}

	comment, err = s.commentRepository.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) Update(ctx context.Context, commentID *entity.ID, req dto.CommentUpdateRequest) (*entity.Comment, error) {
	comment, err := s.commentRepository.FindByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	comment.Content = req.Content

	err = s.commentRepository.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) Delete(ctx context.Context, commentID *entity.ID) error {
	return s.commentRepository.Delete(ctx, commentID)
}

func (s *CommentService) FindByLogID(ctx context.Context, logID *entity.ID, limit, offset int) (*dto.PaginatedResult[*entity.Comment], error) {
	_, err := s.logRepository.FindByID(ctx, logID)
	if err != nil {
		return nil, err
	}

	comment, total, err := s.commentRepository.FindByLogID(ctx, logID)
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedResult[*entity.Comment]{
		Items:  comment,
		Total:  *total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
