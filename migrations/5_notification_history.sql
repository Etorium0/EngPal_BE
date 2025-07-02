CREATE TABLE notification_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    app_name VARCHAR(50) NOT NULL DEFAULT 'EngPal',
    greeting VARCHAR(100),
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Index để query nhanh
CREATE INDEX idx_notification_history_user_id ON notification_history(user_id);
CREATE INDEX idx_notification_history_created_at ON notification_history(created_at);