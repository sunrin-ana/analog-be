package repository

import (
	"analog-be/entity"
	"context"

	"github.com/uptrace/bun"
)

type UserRepository interface {
	FindByID(ctx context.Context, id *entity.ID) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id *entity.ID) error
	FindAll(ctx context.Context, limit, offset int) ([]*entity.User, *int, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.User, *int, error)
}

type UserRepositoryImpl struct {
	db bun.IDB
}

func NewUserRepository(db bun.IDB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id *entity.ID) (*entity.User, error) {
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

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return user, err
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	_, err := r.db.NewUpdate().
		Model(user).
		Where("id = ?", user.ID).
		Exec(ctx)
	return user, err
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id *entity.ID) error {
	_, err := r.db.NewDelete().
		Model((*entity.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, *int, error) {
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

func (r *UserRepositoryImpl) Search(ctx context.Context, query string, limit, offset int) ([]*entity.User, *int, error) {
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
