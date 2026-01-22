package entity

import "time"

// Analog에서 'log'는 article을 의미합니다.
type Log struct {
	ID          string     `bun:"id,pk"`
	Title       string     `bun:"title"`
	Topics      []string   `bun:"topics,array"`
	Generations []uint16   `bun:"generations,array"`
	Content     string     `bun:"content"`
	CreatedAt   time.Time  `bun:"created_at"`
	LoggedBy    []int64    `bun:"logged_by,array"`
	Comments    []*Comment `bun:"rel:has-many,join:log_id=id"`
}

type Comment struct {
	ID        string    `bun:"id,pk"`
	LogID     string    `bun:"log_id"`
	Author    string    `bun:"author"`
	Content   string    `bun:"content"`
	CreatedAt time.Time `bun:"created_at"`
}

type Topic struct {
	Name  string `bun:"name,pk"`
	Count int64  `bun:"count"`
}
