DROP TABLE IF EXISTS log_to_topics;
DROP TABLE IF EXISTS log_to_users;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS oauth_states;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS topics;
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS users;

DROP FUNCTION IF EXISTS cleanup_expired_records;
DROP FUNCTION IF EXISTS update_updated_at_column;

DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP FUNCTION IF EXISTS cleanup_expired_records;
DROP FUNCTION IF EXISTS update_updated_at_column;
