package repository

import (
	"analog-be/entity"
	"context"

	"github.com/uptrace/bun"
)

type CommentRepository struct {
	db bun.IDB
}

func NewCommentRepository(db bun.IDB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) FindByID(ctx context.Context, id string) (*entity.Comment, error) {
	comment := new(entity.Comment)

	err := r.db.NewSelect().
		Model(comment).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepository) FindByLogID(ctx context.Context, logID string) ([]*entity.Comment, error) {
	var comments []*entity.Comment

	err := r.db.NewSelect().
		Model(&comments).
		Where("log_id = ?", logID).
		Order("created_at ASC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	_, err := r.db.NewInsert().
		Model(comment).
		Exec(ctx)
	return err
}

func (r *CommentRepository) Update(ctx context.Context, comment *entity.Comment) error {
	_, err := r.db.NewUpdate().
		Model(comment).
		Where("id = ?", comment.ID).
		Exec(ctx)
	return err
}

func (r *CommentRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*entity.Comment)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *CommentRepository) DeleteByLogID(ctx context.Context, logID string) error {
	_, err := r.db.NewDelete().
		Model((*entity.Comment)(nil)).
		Where("log_id = ?", logID).
		Exec(ctx)
	return err
}
