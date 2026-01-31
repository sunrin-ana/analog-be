package dto

import "analog-be/entity"

type TopicCreateRequest struct {
	Name string `json:"name"`
}

type TopicResponse struct {
	ID    entity.ID `json:"id"`
	Name  string    `json:"name"`
	Count int64     `json:"count"`
}

func NewTopicResponse(t *entity.Topic) TopicResponse {
	return TopicResponse{
		ID:   t.ID,
		Name: t.Name,
	}
}
