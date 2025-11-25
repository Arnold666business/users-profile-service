CREATE TABLE users_profile.users_block_status (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users_profile.users(id) ON DELETE CASCADE,
    block_type_id INT NOT NULL REFERENCES users_profile.block_type_dictionary(type_id) ON DELETE CASCADE,
    forever_flag BOOLEAN DEFAULT FALSE,
    unblock_date DATE,
    is_active BOOLEAN,
    unblock_event_sent BOOLEAN DEFAULT FALSE
);

CREATE INDEX user_id_and_block_type_id_idx ON users_profile.users_block_status(user_id, block_type_id);