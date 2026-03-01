package repository

import (
	"analog-be/entity"
	"context"

	"github.com/uptrace/bun"
)

type CommentRepository interface {
	FindByID(ctx context.Context, id *entity.ID) (*entity.Comment, error)
	FindByLogID(ctx context.Context, logID *entity.ID) ([]*entity.Comment, *int, error)
	Create(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
	Update(ctx context.Context, comment *entity.Comment) error
	Delete(ctx context.Context, id *entity.ID) error
	DeleteByLogID(ctx context.Context, logID *entity.ID) error
}

type CommentRepositoryImpl struct {
	db bun.IDB
}

func NewCommentRepository(db bun.IDB) CommentRepository {
	return &CommentRepositoryImpl{
		db: db,
	}
}

func (r *CommentRepositoryImpl) FindByID(ctx context.Context, id *entity.ID) (*entity.Comment, error) {
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

func (r *CommentRepositoryImpl) FindByLogID(ctx context.Context, logID *entity.ID) ([]*entity.Comment, *int, error) {
	var comments []*entity.Comment

	count, err := r.db.NewSelect().
		Model(&comments).
		Where("log_id = ?", logID).
		Order("created_at ASC").
		ScanAndCount(ctx)

	if err != nil {
		return nil, nil, err
	}

	return comments, &count, nil
}

func (r *CommentRepositoryImpl) Create(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	_, err := r.db.NewInsert().
		Model(comment).
		Exec(ctx)
	return comment, err
}

func (r *CommentRepositoryImpl) Update(ctx context.Context, comment *entity.Comment) error {
	_, err := r.db.NewUpdate().
		Model(comment).
		Where("id = ?", comment.ID).
		Exec(ctx)
	return err
}

func (r *CommentRepositoryImpl) Delete(ctx context.Context, id *entity.ID) error {
	_, err := r.db.NewDelete().
		Model((*entity.Comment)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *CommentRepositoryImpl) DeleteByLogID(ctx context.Context, logID *entity.ID) error {
	_, err := r.db.NewDelete().
		Model((*entity.Comment)(nil)).
		Where("log_id = ?", logID).
		Exec(ctx)
	return err
}
