package service

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/repository"
	"context"
)

type TopicService struct {
	topicRepository *repository.TopicRepository
}

func NewTopicService(topicRepository *repository.TopicRepository) *TopicService {
	return &TopicService{topicRepository: topicRepository}
}

func (s *TopicService) Create(ctx context.Context, req dto.CommentCreateRequest) (*entity.Topic, error) {
	topic := &entity.Topic{
		Name: req.Content,
	}

	return s.topicRepository.Create(ctx, topic)
}
func (s *TopicService) FindAll(ctx context.Context, limit int, offset int) ([]*entity.Topic, error) {
	topics, err := s.topicRepository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return topics, nil
}
func (s *TopicService) Search(ctx context.Context, query string, limit int, offset int) ([]*entity.Topic, error) {
	topics, err := s.topicRepository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return topics, nil
}
func (s *TopicService) Delete(ctx context.Context, id *entity.ID) error {
	return s.topicRepository.Delete(ctx, id)
}
