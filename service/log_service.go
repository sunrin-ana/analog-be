package service

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type LogService struct {
	logRepository     *repository.LogRepository
	commentRepository *repository.CommentRepository
}

func NewLogService(logRepository *repository.LogRepository, commentRepository *repository.CommentRepository) *LogService {
	return &LogService{
		logRepository:     logRepository,
		commentRepository: commentRepository,
	}
}

func (s *LogService) GetLog(ctx context.Context, id string) (*dto.LogResponse, error) {
	log, err := s.logRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	comments := make([]dto.CommentResponse, 0)
	if log.Comments != nil {
		for _, comment := range log.Comments {
			comments = append(comments, dto.CommentResponse{
				ID:        comment.ID,
				LogID:     comment.LogID,
				Author:    comment.Author,
				Content:   comment.Content,
				CreatedAt: comment.CreatedAt.Format(time.RFC3339),
			})
		}
	}

	return &dto.LogResponse{
		ID:          log.ID,
		Title:       log.Title,
		Topics:      log.Topics,
		Generations: log.Generations,
		Content:     log.Content,
		CreatedAt:   log.CreatedAt.Format(time.RFC3339),
		LoggedBy:    log.LoggedBy,
		Comments:    comments,
	}, nil
}

func (s *LogService) GetListOfLog(ctx context.Context, limit int, offset int) (*dto.LogListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	logs, err := s.logRepository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.logRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	logResponses := make([]*dto.LogResponse, len(logs))
	for i, log := range logs {
		logResponses[i] = &dto.LogResponse{
			ID:          log.ID,
			Title:       log.Title,
			Topics:      log.Topics,
			Generations: log.Generations,
			Content:     log.Content,
			CreatedAt:   log.CreatedAt.Format(time.RFC3339),
			LoggedBy:    log.LoggedBy,
		}
	}

	return &dto.LogListResponse{
		Logs:   logResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *LogService) SearchLogs(ctx context.Context, query string, limit int, offset int) (*dto.LogListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	logs, err := s.logRepository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	logResponses := make([]*dto.LogResponse, len(logs))
	for i, log := range logs {
		logResponses[i] = &dto.LogResponse{
			ID:          log.ID,
			Title:       log.Title,
			Topics:      log.Topics,
			Generations: log.Generations,
			Content:     log.Content,
			CreatedAt:   log.CreatedAt.Format(time.RFC3339),
			LoggedBy:    log.LoggedBy,
		}
	}

	return &dto.LogListResponse{
		Logs:   logResponses,
		Total:  len(logs),
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *LogService) CreateLog(ctx context.Context, req dto.LogCreateRequest) (*dto.LogResponse, error) {
	now := time.Now().UTC()

	log := &entity.Log{
		ID:          req.ID,
		Title:       req.Title,
		Topics:      req.Topics,
		Generations: req.Generations,
		Content:     req.Content,
		CreatedAt:   now,
		LoggedBy:    req.LoggedBy,
	}

	if log.Topics == nil {
		log.Topics = []string{}
	}
	if log.Generations == nil {
		log.Generations = []uint16{}
	}
	if log.LoggedBy == nil {
		log.LoggedBy = []int64{}
	}

	err := s.logRepository.Create(ctx, log)
	if err != nil {
		return nil, err
	}

	return &dto.LogResponse{
		ID:          log.ID,
		Title:       log.Title,
		Topics:      log.Topics,
		Generations: log.Generations,
		Content:     log.Content,
		CreatedAt:   log.CreatedAt.Format(time.RFC3339),
		LoggedBy:    log.LoggedBy,
		Comments:    []dto.CommentResponse{},
	}, nil
}

func (s *LogService) UpdateLog(ctx context.Context, id string, req dto.LogUpdateRequest) (*dto.LogResponse, error) {
	log, err := s.logRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		log.Title = *req.Title
	}
	if req.Topics != nil {
		log.Topics = *req.Topics
	}
	if req.Generations != nil {
		log.Generations = *req.Generations
	}
	if req.Content != nil {
		log.Content = *req.Content
	}
	if req.LoggedBy != nil {
		log.LoggedBy = *req.LoggedBy
	}

	err = s.logRepository.Update(ctx, log)
	if err != nil {
		return nil, err
	}

	return &dto.LogResponse{
		ID:          log.ID,
		Title:       log.Title,
		Topics:      log.Topics,
		Generations: log.Generations,
		Content:     log.Content,
		CreatedAt:   log.CreatedAt.Format(time.RFC3339),
		LoggedBy:    log.LoggedBy,
	}, nil
}

func (s *LogService) DeleteLog(ctx context.Context, id string) error {
	err := s.commentRepository.DeleteByLogID(ctx, id)
	if err != nil {
		return err
	}

	return s.logRepository.Delete(ctx, id)
}

func (s *LogService) CreateComment(ctx context.Context, logID string, req dto.CommentCreateRequest) (*dto.CommentResponse, error) {
	_, err := s.logRepository.FindByID(ctx, logID)
	if err != nil {
		return nil, err
	}

	comment := &entity.Comment{
		ID:        uuid.New().String(),
		LogID:     logID,
		Author:    req.Author,
		Content:   req.Content,
		CreatedAt: time.Now().UTC(),
	}

	err = s.commentRepository.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &dto.CommentResponse{
		ID:        comment.ID,
		LogID:     comment.LogID,
		Author:    comment.Author,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *LogService) UpdateComment(ctx context.Context, commentID string, req dto.CommentUpdateRequest) (*dto.CommentResponse, error) {
	comment, err := s.commentRepository.FindByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	comment.Content = req.Content

	err = s.commentRepository.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &dto.CommentResponse{
		ID:        comment.ID,
		LogID:     comment.LogID,
		Author:    comment.Author,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *LogService) DeleteComment(ctx context.Context, commentID string) error {
	return s.commentRepository.Delete(ctx, commentID)
}
