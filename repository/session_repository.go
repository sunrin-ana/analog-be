package repository

import (
	"analog-be/entity"
	"context"
	"time"

	"github.com/uptrace/bun"
)

type OAuthStateRepository struct {
	db bun.IDB
}

func NewOAuthStateRepository(db bun.IDB) *OAuthStateRepository {
	return &OAuthStateRepository{
		db: db,
	}
}

func (r *OAuthStateRepository) Create(ctx context.Context, state *entity.OAuthState) error {
	_, err := r.db.NewInsert().
		Model(state).
		Exec(ctx)
	return err
}

func (r *OAuthStateRepository) FindByState(ctx context.Context, state string) (*entity.OAuthState, error) {
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

func (r *OAuthStateRepository) Delete(ctx context.Context, state string) error {
	_, err := r.db.NewDelete().
		Model((*entity.OAuthState)(nil)).
		Where("state = ?", state).
		Exec(ctx)
	return err
}

func (r *OAuthStateRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*entity.OAuthState)(nil)).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)
	return err
}

type SessionRepository struct {
	db bun.IDB
}

func NewSessionRepository(db bun.IDB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *entity.Session) error {
	_, err := r.db.NewInsert().
		Model(session).
		Exec(ctx)
	return err
}

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*entity.Session, error) {
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

func (r *SessionRepository) Delete(ctx context.Context, token string) error {
	_, err := r.db.NewDelete().
		Model((*entity.Session)(nil)).
		Where("session_token = ?", token).
		Exec(ctx)
	return err
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	_, err := r.db.NewDelete().
		Model((*entity.Session)(nil)).
		Where("user_id = ?", userID).
		Exec(ctx)
	return err
}

func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*entity.Session)(nil)).
		Where("expires_at < ?", time.Now().UTC()).
		Exec(ctx)
	return err
}
