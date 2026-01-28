package repository

import (
	"analog-be/entity"
	"context"

	"github.com/uptrace/bun"
)

type UserRepository struct {
	db bun.IDB
}

func NewUserRepository(db bun.IDB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByID(ctx context.Context, id *entity.ID) (*entity.User, error) {
	user := new(entity.User)

	err := r.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return user, err
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	_, err := r.db.NewUpdate().
		Model(user).
		Where("id = ?", user.ID).
		Exec(ctx)
	return user, err
}

func (r *UserRepository) Delete(ctx context.Context, id *entity.ID) error {
	_, err := r.db.NewDelete().
		Model((*entity.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, *int, error) {
	var users []*entity.User

	count, err := r.db.NewSelect().
		Model(&users).
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)

	if err != nil {
		return nil, nil, err
	}

	return users, &count, nil
}

func (r *UserRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.User, *int, error) {
	var users []*entity.User

	count, err := r.db.NewSelect().
		Model(&users).
		Where("name ILIKE ?", "%"+query+"%").
		WhereOr("part_of ILIKE ?", "%"+query+"%").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)

	if err != nil {
		return nil, nil, err
	}

	return users, &count, nil
}
