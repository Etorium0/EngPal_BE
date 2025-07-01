CREATE TABLE user_push_token (
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    token TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);
