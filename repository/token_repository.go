package repository

import (
	"analog-be/entity"
	"context"

	"github.com/uptrace/bun"
)

type TokenRepository struct {
	db bun.IDB
}

func NewTokenRepository(db bun.IDB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) FindByID(ctx context.Context, refreshTokenID string) (*entity.RefreshToken, error) {
	refToken := new(entity.RefreshToken)

	err := r.db.NewSelect().
		Model(refToken).
		Where("id = ?", refreshTokenID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return refToken, nil
}

func (r *TokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	_, err := r.db.NewInsert().
		Model((*entity.RefreshToken)(nil)).
		Returning("id").
		Exec(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (r *TokenRepository) Delete(ctx context.Context, id *entity.ID) error {
	_, err := r.db.NewDelete().
		Model((*entity.RefreshToken)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
