package dto

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
	ID          string   `json:"id" validate:"required,max=100"`
	Title       string   `json:"title" validate:"required,min=1,max=200"`
	Topics      []string `json:"topics" validate:"max=20,dive,max=50"`
	Generations []uint16 `json:"generations" validate:"max=50"`
	Content     string   `json:"content" validate:"required,min=1,max=50000"`
	LoggedBy    []int64  `json:"loggedBy" validate:"max=100"`
}

type LogUpdateRequest struct {
	Title       *string   `json:"title"`
	Topics      *[]string `json:"topics"`
	Generations *[]uint16 `json:"generations"`
	Content     *string   `json:"content"`
	LoggedBy    *[]int64  `json:"loggedBy"`
}

type LogResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Topics      []string          `json:"topics"`
	Generations []uint16          `json:"generations"`
	Content     string            `json:"content"`
	CreatedAt   string            `json:"createdAt"`
	LoggedBy    []int64           `json:"loggedBy"`
	Comments    []CommentResponse `json:"comments,omitempty"`
}

type LogListResponse struct {
	Logs   []*LogResponse `json:"logs"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

type CommentCreateRequest struct {
	Author  string `json:"author" validate:"required,min=1,max=100"`
	Content string `json:"content" validate:"required,min=1,max=5000"`
}

type CommentUpdateRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}

type CommentResponse struct {
	ID        string `json:"id"`
	LogID     string `json:"logId"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}
