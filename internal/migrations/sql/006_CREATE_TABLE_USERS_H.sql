CREATE TABLE users_profile.users_h (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users_profile.users(id),
    action VARCHAR(50),
    create_at TIMESTAMP DEFAULT NOW(),
    email VARCHAR(255),
    login VARCHAR(255),
    old_fields JSONB
);

CREATE INDEX idx_user_id ON users_profile.users_h(user_id)