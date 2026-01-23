package entity

import (
	"github.com/uptrace/bun"
	"time"
)

// Log 는 Analog에서 article을 의미합니다.
type Log struct {
	bun.BaseModel `bun:"table:logs"`

	ID          ID        `bun:"id,pk,autoincrement"`
	Title       string    `bun:"title"`
	Topics      []*Topic  `bun:"m2m:log_to_topics,join:Log=Topic"`
	Generations []uint16  `bun:"generations,array"`
	Content     string    `bun:"content"`
	CreatedAt   time.Time `bun:"created_at"`
	LoggedBy    []*User   `bun:"m2m:log_to_users,join:Log=User"`
}

type Comment struct {
	bun.BaseModel `bun:"table:comments"`

	ID        ID        `bun:"id,pk,autoincrement"`
	LogID     ID        `bun:"log_id"`
	Log       *Log      `bun:"rel:belongs-to,join:log_id=id"`
	AuthorID  ID        `bun:"author_id"`
	Author    *User     `bun:"rel:belongs-to,join:author_id=id"`
	Content   string    `bun:"content"`
	CreatedAt time.Time `bun:"created_at"`
}

type Topic struct {
	bun.BaseModel `bun:"table:topics"`

	ID   ID     `bun:"id,pk,autoincrement"`
	Name string `bun:"name,unique"`

	Count int64 `bun:"-"`
}

type LogToUser struct {
	bun.BaseModel `bun:"table:log_to_users"`

	UserID ID `bun:"user_id,pk"`
	LogID  ID `bun:"log_id,pk"`

	Log  *Log  `bun:"rel:belongs-to,join:log_id=id"`
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

type LogToTopic struct {
	bun.BaseModel `bun:"table:log_to_topics"`

	LogID   ID `bun:"log_id,pk"`
	TopicID ID `bun:"topic_id,pk"`

	Log   *Log   `bun:"rel:belongs-to,join:log_id=id"`
	Topic *Topic `bun:"rel:belongs-to,join:topic_id=id"`
}
