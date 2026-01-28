package dto

import (
	"analog-be/entity"
	"time"
)

type UserCreateRequest struct {
	ID           entity.ID `json:"id"`
	Name         string    `json:"name" binding:"required"`
	ProfileImage string    `json:"profileImage"`
	PartOf       string    `json:"partOf"`
	Generation   uint16    `json:"generation"`
	Connections  []string  `json:"connections"`
}

type UserUpdateRequest struct {
	Name         *string   `json:"name"`
	ProfileImage *string   `json:"profileImage"`
	PartOf       *string   `json:"partOf"`
	Generation   *uint16   `json:"generation"`
	Connections  *[]string `json:"connections"`
}

type UserResponse struct {
	ID           entity.ID `json:"id"`
	Name         string    `json:"name"`
	ProfileImage string    `json:"profileImage"`
	JoinedAt     time.Time `json:"joinedAt"`
	PartOf       string    `json:"partOf"`
	Generation   uint16    `json:"generation"`
	Connections  []string  `json:"connections"`
}

type UserListResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

func NewUserResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		ProfileImage: user.ProfileImage,
		JoinedAt:     user.JoinedAt,
		PartOf:       user.PartOf,
		Generation:   user.Generation,
		Connections:  user.Connections,
	}
}
