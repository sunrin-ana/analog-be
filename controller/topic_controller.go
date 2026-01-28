package controller

import (
	"analog-be/dto"
	"analog-be/service"
	"context"
	"github.com/NARUBROWN/spine/pkg/httpx"
	"github.com/NARUBROWN/spine/pkg/query"
	"net/http"
)

type TopicController struct {
	topicService *service.TopicService
}

func NewTopicController(topicService *service.TopicService) *TopicController {
	return &TopicController{topicService: topicService}
}

func (c *TopicController) GetList(ctx context.Context, page query.Pagination) httpx.Response[[]dto.TopicResponse] {

	topics, err := c.topicService.FindAll(ctx, page.Size, page.Page)
	if err != nil {
		return httpx.Response[[]dto.TopicResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	topicResponses := make([]dto.TopicResponse, len(topics))
	for i, topic := range topics {
		topicResponses[i] = dto.NewTopicResponse(topic)
	}

	return httpx.Response[[]dto.TopicResponse]{
		Body: topicResponses,
	}
}
