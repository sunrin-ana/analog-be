package repository

import (
	"analog-be/entity"
	"context"

	"github.com/uptrace/bun"
)

type LogRepository struct {
	db bun.IDB
}

func NewLogRepository(db bun.IDB) *LogRepository {
	return &LogRepository{
		db: db,
	}
}

func (r *LogRepository) FindByID(ctx context.Context, id *entity.ID) (*entity.Log, error) {
	log := new(entity.Log)

	err := r.db.NewSelect().
		Model(log).
		Where("id = ?", id).
		Relation("Topics").
		Relation("LoggedBy").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return log, nil
}

func (r *LogRepository) FindAll(ctx context.Context, limit int, offset int) ([]*entity.Log, *int, error) {
	var logs []*entity.Log

	count, err := r.db.NewSelect().
		Model(&logs).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)

	if err != nil {
		return nil, nil, err
	}

	return logs, &count, nil
}

func (r *LogRepository) FindAllByTopicID(ctx context.Context, topicID *entity.ID, limit int, offset int) ([]*entity.Log, *int, error) {
	var logs []*entity.Log

	count, err := r.db.NewSelect().
		Model(&logs).
		Join("JOIN log_to_topics ltt ON ltt.log_id = log.id").
		Where("ltt.topic_id = ?", topicID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)

	if err != nil {
		return nil, nil, err
	}

	return logs, &count, nil
}

func (r *LogRepository) FindAllByGeneration(ctx context.Context, generation uint16, limit, offset int) ([]*entity.Log, *int, error) {
	var logs []*entity.Log

	count, err := r.db.NewSelect().
		Model(&logs).
		Where("? = ANY(generations)", generation).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)

	if err != nil {
		return nil, nil, err
	}

	return logs, &count, nil
}

func (r *LogRepository) Search(ctx context.Context, query string, limit int, offset int) ([]*entity.Log, *int, error) {
	var logs []*entity.Log

	count, err := r.db.NewSelect().
		Model(&logs).
		Where("title ILIKE ? OR content ILIKE ?", "%"+query+"%", "%"+query+"%").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)

	if err != nil {
		return nil, nil, err
	}

	return logs, &count, nil
}

func (r *LogRepository) Create(ctx context.Context, log *entity.Log, topicIDs, authorIDs *[]entity.ID) (*entity.Log, error) {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(log).Exec(ctx); err != nil {
			return err
		}

		if len(*topicIDs) > 0 {
			log2topic := make([]entity.LogToTopic, 0, len(*topicIDs))
			for _, tid := range *topicIDs {
				log2topic = append(log2topic, entity.LogToTopic{
					LogID:   log.ID,
					TopicID: tid,
				})
			}
			if _, err := tx.NewInsert().Model(&log2topic).Exec(ctx); err != nil {
				return err
			}
		}

		if len(*authorIDs) > 0 {
			log2user := make([]entity.LogToUser, 0, len(*authorIDs))
			for _, uid := range *authorIDs {
				log2user = append(log2user, entity.LogToUser{
					LogID:  log.ID,
					UserID: uid,
				})
			}
			if _, err := tx.NewInsert().Model(&log2user).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return log, err
}

func (r *LogRepository) Update(ctx context.Context, log *entity.Log, topicIDs, authorIDs *[]entity.ID) (*entity.Log, error) {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(log).WherePK().Exec(ctx); err != nil {
			return err
		}

		logID := log.ID

		if _, err := tx.NewDelete().Model((*entity.LogToTopic)(nil)).Where("log_id = ?", logID).Exec(ctx); err != nil {
			return err
		}

		if topicIDs != nil {
			log2topic := make([]entity.LogToTopic, 0, len(*topicIDs))
			for _, tid := range *topicIDs {
				log2topic = append(log2topic, entity.LogToTopic{
					LogID:   logID,
					TopicID: tid,
				})
			}
			if _, err := tx.NewInsert().Model(&log2topic).Exec(ctx); err != nil {
				return err
			}
		}

		if _, err := tx.NewDelete().Model((*entity.LogToUser)(nil)).Where("log_id = ?", logID).Exec(ctx); err != nil {
			return err
		}

		if authorIDs != nil {
			log2user := make([]entity.LogToUser, 0, len(*authorIDs))
			for _, uid := range *authorIDs {
				log2user = append(log2user, entity.LogToUser{
					LogID:  logID,
					UserID: uid,
				})
			}
			if _, err := tx.NewInsert().Model(&log2user).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return log, err
}

func (r *LogRepository) Delete(ctx context.Context, id *entity.ID) error {
	_, err := r.db.NewDelete().
		Model((*entity.Log)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
