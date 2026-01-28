package entity

import (
	"time"
)

type User struct {
	ID           ID        `bun:"id,pk"` // An-Account와 아이디를 공유합니다
	Name         string    `bun:"name"`
	ProfileImage string    `bun:"profile_image"`
	JoinedAt     time.Time `bun:"joined_at"`
	PartOf       string    `bun:"part_of"`
	Generation   uint16    `bun:"generation"`
	Connections  []string  `bun:"connections,array"`
}
