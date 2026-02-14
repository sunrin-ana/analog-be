package dto

import (
	"analog-be/entity"
	"time"
)

type LogListRequest struct {
	Limit  int `json:"limit" form:"limit"`
	Offset int `json:"offset" form:"offset"`
}

type LogSearchRequest struct {
	Query  string `json:"q" form:"q" binding:"required"`
	Limit  int    `json:"limit" form:"limit"`
	Offset int    `json:"offset" form:"offset"`
}

type LogCreateRequest struct {
	Title       string      `json:"title" validate:"required,min=1,max=200"`
	TopicIDs    []entity.ID `json:"topicIDs" validate:"max=20,dive,max=50"`
	Generations []uint16    `json:"generations" validate:"max=50"`
	Content     string      `json:"content" validate:"required,min=1,max=50000"`
	CoAuthorIDs []entity.ID `json:"coAuthorIDs" validate:"max=100"`
}

type LogUpdateRequest struct {
	Title       *string      `json:"title"`
	TopicIDs    *[]entity.ID `json:"topicIDs"`
	Generations *[]uint16    `json:"generations"`
	Content     *string      `json:"content"`
	CoAuthorIDs *[]entity.ID `json:"coAuthorIDs"`
}

type LogResponse struct {
	ID          entity.ID       `json:"id"`
	Title       string          `json:"title"`
	Topics      []TopicResponse `json:"topics"`
	Generations []uint16        `json:"generations"`
	Content     string          `json:"content"`
	CreatedAt   time.Time       `json:"createdAt"`
	LoggedBy    []UserResponse  `json:"loggedBy"`
}

type CommentCreateRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}

type CommentUpdateRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}

type CommentResponse struct {
	ID        entity.ID    `json:"id"`
	LogID     entity.ID    `json:"logId"`
	Author    UserResponse `json:"author"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"createdAt"`
}

type TopicResponse struct {
	ID    entity.ID `json:"id"`
	Name  string    `json:"name"`
	Count int64     `json:"count"`
}

func NewLogResponse(l *entity.Log) LogResponse {
	var topics []TopicResponse
	if l.Topics != nil {
		topics = make([]TopicResponse, len(l.Topics))
		for i, topic := range l.Topics {
			topics[i] = NewTopicResponse(topic)
		}
	}

	var loggedBy []UserResponse
	if l.LoggedBy != nil {
		loggedBy = make([]UserResponse, len(l.LoggedBy))
		for i, user := range l.LoggedBy {
			loggedBy[i] = NewUserResponse(user)
		}
	}

	return LogResponse{
		ID:          l.ID,
		Title:       l.Title,
		Topics:      topics,
		Generations: l.Generations,
		Content:     l.PreRendered,
		CreatedAt:   l.CreatedAt,
		LoggedBy:    loggedBy,
	}
}

func NewCommentResponse(c *entity.Comment) CommentResponse {
	var author UserResponse
	if c.Author != nil {
		author = NewUserResponse(c.Author)
	}

	return CommentResponse{
		ID:        c.ID,
		LogID:     c.LogID,
		Author:    author,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
	}
}

func NewTopicResponse(t *entity.Topic) TopicResponse {
	return TopicResponse{
		ID:   t.ID,
		Name: t.Name,
	}
}
