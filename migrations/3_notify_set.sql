CREATE TABLE user_notification_setting (
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    enable_notify BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMP DEFAULT NOW()
);
