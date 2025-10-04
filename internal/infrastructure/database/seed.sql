-- Seed data for testing
INSERT OR IGNORE INTO users (id, name, email, password_hash, created_at) VALUES
    ('user-1', 'Demo User', 'demo@example.com', '$2a$10$placeholder', datetime('now')),
    ('user-2', 'Test User', 'test@example.com', '$2a$10$placeholder', datetime('now'));
