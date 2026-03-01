package controller

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/service"
	"context"
	"log"
	"net/http"

	"github.com/NARUBROWN/spine/pkg/httpx"
	"github.com/NARUBROWN/spine/pkg/query"
)

type TopicController struct {
	topicService service.TopicService
}

func NewTopicController(topicService service.TopicService) *TopicController {
	return &TopicController{topicService: topicService}
}

// GetList gets a paginated list of topics.
// @Summary      GetListOfTopics
// @Description  Get a paginated list of topics.
// @Tags         Topic
// @Produce      json
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {array} dto.TopicResponse
// @Failure      404 "Not Found"
// @Router       /topic [get]
func (c *TopicController) GetList(ctx context.Context, page query.Pagination) httpx.Response[[]dto.TopicResponse] {

	topics, err := c.topicService.FindAll(ctx, page.Size, page.Page)
	if err != nil {
		return httpx.Response[[]dto.TopicResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusNotFound, // not found
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

// Create creates a new topic.
// @Summary      CreateTopic
// @Description  Create a new topic.
// @Tags         Topic
// @Accept       json
// @Produce      json
// @Param        topic body dto.TopicCreateRequest true "Topic to create"
// @Success      201 {object} dto.TopicResponse
// @Failure      400 "Bad Request"
// @Failure      500 "Internal Server Error"
// @Router       /topic [post]
func (c *TopicController) Create(ctx context.Context, req *dto.TopicCreateRequest) httpx.Response[dto.TopicResponse] {
	if req.Name == "" {
		return httpx.Response[dto.TopicResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusBadRequest, // name is required
			},
		}
	}

	topic, err := c.topicService.Create(ctx, &entity.Topic{
		Name: req.Name,
	})
	if err != nil {
		log.Println(err.Error())
		return httpx.Response[dto.TopicResponse]{
			Options: httpx.ResponseOptions{
				Status: http.StatusInternalServerError, // internal server error
			},
		}
	}

	return httpx.Response[dto.TopicResponse]{
		Body: dto.NewTopicResponse(topic),
		Options: httpx.ResponseOptions{
			Status: http.StatusCreated,
		},
	}
}
