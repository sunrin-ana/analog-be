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

func (r *LogRepository) FindByID(ctx context.Context, id string) (*entity.Log, error) {
	log := new(entity.Log)

	err := r.db.NewSelect().
		Model(log).
		Where("id = ?", id).
		Relation("Comments").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return log, nil
}

func (r *LogRepository) FindAll(ctx context.Context, limit int, offset int) ([]*entity.Log, error) {
	var logs []*entity.Log

	err := r.db.NewSelect().
		Model(&logs).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *LogRepository) Search(ctx context.Context, query string, limit int, offset int) ([]*entity.Log, error) {
	var logs []*entity.Log

	err := r.db.NewSelect().
		Model(&logs).
		Where("title ILIKE ? OR content ILIKE ?", "%"+query+"%", "%"+query+"%").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *LogRepository) Create(ctx context.Context, log *entity.Log) error {
	_, err := r.db.NewInsert().
		Model(log).
		Exec(ctx)
	return err
}

func (r *LogRepository) Update(ctx context.Context, log *entity.Log) error {
	_, err := r.db.NewUpdate().
		Model(log).
		Where("id = ?", log.ID).
		Exec(ctx)
	return err
}

func (r *LogRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*entity.Log)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *LogRepository) Count(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().
		Model((*entity.Log)(nil)).
		Count(ctx)
	return count, err
}
