package service

import (
	"analog-be/entity"
	"analog-be/repository"
	"context"
)

type TopicService interface {
	Create(ctx context.Context, topic *entity.Topic) (*entity.Topic, error)
	FindAll(ctx context.Context, limit int, offset int) ([]*entity.Topic, error)
	Search(ctx context.Context, query string, limit int, offset int) ([]*entity.Topic, error)
	Delete(ctx context.Context, id *entity.ID) error
}

type TopicServiceImpl struct {
	topicRepository repository.TopicRepository
}

func NewTopicService(topicRepository repository.TopicRepository) TopicService {
	return &TopicServiceImpl{topicRepository: topicRepository}
}

func (s *TopicServiceImpl) Create(ctx context.Context, topic *entity.Topic) (*entity.Topic, error) {
	return s.topicRepository.Create(ctx, topic)
}
func (s *TopicServiceImpl) FindAll(ctx context.Context, limit int, offset int) ([]*entity.Topic, error) {
	topics, err := s.topicRepository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return topics, nil
}
func (s *TopicServiceImpl) Search(ctx context.Context, query string, limit int, offset int) ([]*entity.Topic, error) {
	topics, err := s.topicRepository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return topics, nil
}
func (s *TopicServiceImpl) Delete(ctx context.Context, id *entity.ID) error {
	return s.topicRepository.Delete(ctx, id)
}
