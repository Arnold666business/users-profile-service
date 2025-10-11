CREATE TABLE users_profile.users_block_status (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users_profile.users(id) ON DELETE CASCADE,
    block_type_id INT NOT NULL REFERENCES users_profile.block_type_dictionary(id) ON DELETE CASCADE,
    forever_flag BOOLEAN DEFAULT FALSE,
    unblock_date DATE
);

CREATE INDEX idx_user_id ON users_profile.users_block_status(user_id);
CREATE INDEX idx_unblock_date_forever_flag ON users_profile.users_block_status(unblock_date, forever_flag)
