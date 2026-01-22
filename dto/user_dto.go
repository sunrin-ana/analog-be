package dto

type UserCreateRequest struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name" binding:"required"`
	ProfileImage string   `json:"profileImage"`
	PartOf       string   `json:"partOf"`
	Generation   uint16   `json:"generation"`
	Connections  []string `json:"connections"`
}

type UserUpdateRequest struct {
	Name         *string   `json:"name"`
	ProfileImage *string   `json:"profileImage"`
	PartOf       *string   `json:"partOf"`
	Generation   *uint16   `json:"generation"`
	Connections  *[]string `json:"connections"`
}

type UserResponse struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	ProfileImage string   `json:"profileImage"`
	JoinedAt     string   `json:"joinedAt"`
	PartOf       string   `json:"partOf"`
	Generation   uint16   `json:"generation"`
	Connections  []string `json:"connections"`
}

type UserListResponse struct {
	Users  []*UserResponse `json:"users"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}
