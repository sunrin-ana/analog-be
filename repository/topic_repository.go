package repository

import (
	"analog-be/entity"
	"context"
	"github.com/uptrace/bun"
)

type TopicRepository struct {
	db bun.IDB
}

func NewTopicRepository(db bun.IDB) *TopicRepository {
	return &TopicRepository{db: db}
}

func (r *TopicRepository) Create(ctx context.Context, topic *entity.Topic) (*entity.Topic, error) {
	_, err := r.db.NewInsert().Model(topic).Exec(ctx)
	return topic, err
}

func (r *TopicRepository) FindAll(ctx context.Context, limit int, offset int) ([]*entity.Topic, error) {
	var topics []*entity.Topic

	err := r.db.NewSelect().
		Model(&topics).
		Column("topic.id", "topic.name").
		ColumnExpr("COUNT(ltt.log_id) AS count").
		Join("LEFT JOIN log_to_topics AS ltt ON ltt.topic_id = topic.id").
		Group("topic.id", "topic.name").
		Order("topic.name ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *TopicRepository) Search(ctx context.Context, query string, limit int, offset int) ([]*entity.Topic, error) {
	var topics []*entity.Topic

	err := r.db.NewSelect().
		Model(&topics).
		Where("topic.name ILIKE ?", "%"+query+"%").
		Column("topic.id", "topic.name").
		ColumnExpr("COUNT(ltt.log_id) AS count").
		Join("LEFT JOIN log_to_topics AS ltt ON ltt.topic_id = topic.id").
		Group("topic.id", "topic.name").
		Order("topic.name ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *TopicRepository) Delete(ctx context.Context, id *entity.ID) error {
	_, err := r.db.NewDelete().
		Model((*entity.Topic)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
