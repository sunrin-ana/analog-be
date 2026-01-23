package entity

import "time"

type OAuthState struct {
	ID           ID        `bun:"id,pk,autoincrement" json:"id"`
	State        string    `bun:"state,unique,notnull" json:"state"`
	CodeVerifier string    `bun:"code_verifier,notnull" json:"codeVerifier"`
	RedirectUri  string    `bun:"redirect_uri" json:"redirectUri"`
	IsSignup     bool      `bun:"is_signup,notnull,default:false" json:"isSignup"`
	ExpiresAt    time.Time `bun:"expires_at,notnull" json:"expiresAt"`
	CreatedAt    time.Time `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
}

type Session struct {
	ID           ID        `bun:"id,pk,autoincrement" json:"id"`
	SessionToken string    `bun:"session_token,unique,notnull" json:"sessionToken"`
	UserID       ID        `bun:"user_id,notnull" json:"userId"`
	ExpiresAt    time.Time `bun:"expires_at,notnull" json:"expiresAt"`
	CreatedAt    time.Time `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`

	User *User `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}
