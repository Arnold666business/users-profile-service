CREATE TABLE IF NOT EXISTS users_profile.users (
    id BIGSERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    access_email_status BOOLEAN DEFAULT FALSE,
    role INT NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    delete_at DATE
)

CREATE INDEX email_idx ON users_profile.users(email);
CREATE INDEX login_idx ON users_profile.users(login);
