package controller

import (
	"analog-be/dto"
	"analog-be/service"
	"context"
	"github.com/NARUBROWN/spine/pkg/query"
	"strconv"
)

type TopicController struct {
	topicService *service.TopicService
}

func NewTopicController(topicService *service.TopicService) *TopicController {
	return &TopicController{topicService: topicService}
}

func (c *TopicController) GetList(ctx context.Context, q query.Values) ([]dto.TopicResponse, error) {
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

	topics, err := c.topicService.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	topicResponses := make([]dto.TopicResponse, len(topics))
	for i, topic := range topics {
		topicResponses[i] = dto.NewTopicResponse(topic)
	}

	return topicResponses, nil
}
