-- Analog Backend Database Schema
-- PostgreSQL 필요

CREATE TABLE users (
    id BIGINT PRIMARY KEY,  -- AnAccount와 같음
    name VARCHAR(255) NOT NULL,
    handle VARCHAR(255) UNIQUE NOT NULL,
    profile_image VARCHAR(500),
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    part_of VARCHAR(100),
    generation SMALLINT,
    connections TEXT[] DEFAULT '{}',  -- Array of connections
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- Article 개념
CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    description VARCHAR(100) NOT NULL,
    generations SMALLINT[] DEFAULT '{}',
    content TEXT NOT NULL,
    pre_rendered TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_logs_created_at ON logs(created_at DESC);

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    log_id BIGINT NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE topics (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE INDEX idx_comments_log_id ON comments(log_id);
CREATE INDEX idx_comments_created_at ON comments(created_at);

CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

CREATE TABLE o_auth_states (
    id BIGSERIAL PRIMARY KEY,
    state VARCHAR(255) UNIQUE NOT NULL,
    code_verifier VARCHAR(255) NOT NULL,
    redirect_uri VARCHAR(500),
    is_signup BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_o_auth_states_expires_at ON o_auth_states(expires_at);

CREATE TABLE tokens (
    id BIGSERIAL PRIMARY KEY,
    value VARCHAR(255) UNIQUE NOT NULL,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    issued_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
    );

CREATE TABLE log_to_users (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_id BIGINT NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, log_id)
);

CREATE TABLE log_to_topics (
    log_id BIGINT NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    topic_id BIGINT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    PRIMARY KEY (log_id, topic_id)
);

CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_refresh_token ON tokens(refresh_token);
CREATE INDEX idx_tokens_expires_at ON tokens(expires_at);


-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for users table
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Cleanup function for expired sessions and oauth states
CREATE OR REPLACE FUNCTION cleanup_expired_records()
RETURNS void AS $$
BEGIN
DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;
DELETE FROM o_auth_states WHERE expires_at < CURRENT_TIMESTAMP;
DELETE FROM tokens WHERE expires_at < CURRENT_TIMESTAMP;
END;
$$ language 'plpgsql';