package repository

import (
	"analog-be/entity"
	"context"
	"time"

	"github.com/uptrace/bun"
)

type OAuthStateRepository interface {
	Create(ctx context.Context, state *entity.OAuthState) error
	FindByState(ctx context.Context, state string) (*entity.OAuthState, error)
	Delete(ctx context.Context, state string) error
	DeleteExpired(ctx context.Context) error
}

type OAuthStateRepositoryImpl struct {
	db bun.IDB
}

func NewOAuthStateRepository(db bun.IDB) OAuthStateRepository {
	return &OAuthStateRepositoryImpl{
		db: db,
	}
}

func (r *OAuthStateRepositoryImpl) Create(ctx context.Context, state *entity.OAuthState) error {
	_, err := r.db.NewInsert().
		Model(state).
		Exec(ctx)
	return err
}

func (r *OAuthStateRepositoryImpl) FindByState(ctx context.Context, state string) (*entity.OAuthState, error) {
	oauthState := new(entity.OAuthState)
	err := r.db.NewSelect().
		Model(oauthState).
		Where("state = ?", state).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return oauthState, nil
}

func (r *OAuthStateRepositoryImpl) Delete(ctx context.Context, state string) error {
	_, err := r.db.NewDelete().
		Model((*entity.OAuthState)(nil)).
		Where("state = ?", state).
		Exec(ctx)
	return err
}

func (r *OAuthStateRepositoryImpl) DeleteExpired(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*entity.OAuthState)(nil)).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)
	return err
}

type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) error
	FindByToken(ctx context.Context, token string) (*entity.Session, error)
	Delete(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID int64) error
	DeleteExpired(ctx context.Context) error
}

type SessionRepositoryImpl struct {
	db bun.IDB
}

func NewSessionRepository(db bun.IDB) SessionRepository {
	return &SessionRepositoryImpl{
		db: db,
	}
}

func (r *SessionRepositoryImpl) Create(ctx context.Context, session *entity.Session) error {
	_, err := r.db.NewInsert().
		Model(session).
		Exec(ctx)
	return err
}

func (r *SessionRepositoryImpl) FindByToken(ctx context.Context, token string) (*entity.Session, error) {
	session := new(entity.Session)
	err := r.db.NewSelect().
		Model(session).
		Where("session_token = ?", token).
		Relation("User").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepositoryImpl) Delete(ctx context.Context, token string) error {
	_, err := r.db.NewDelete().
		Model((*entity.Session)(nil)).
		Where("session_token = ?", token).
		Exec(ctx)
	return err
}

func (r *SessionRepositoryImpl) DeleteByUserID(ctx context.Context, userID int64) error {
	_, err := r.db.NewDelete().
		Model((*entity.Session)(nil)).
		Where("user_id = ?", userID).
		Exec(ctx)
	return err
}

func (r *SessionRepositoryImpl) DeleteExpired(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*entity.Session)(nil)).
		Where("expires_at < ?", time.Now().UTC()).
		Exec(ctx)
	return err
}
