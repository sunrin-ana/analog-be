package entity

import "time"

type RefreshToken struct {
	Token     string    `bun:"token,unique,notnull"`
	IssuedAt  time.Time `bun:"issued_at,notnull"`
	ExpiresAt time.Time `bun:"expires_at,notnull"`
	UserID    ID        `bun:"user_id,notnull"`
	User      *User     `bun:"rel:belongs-to,join:user_id=id"`
}
