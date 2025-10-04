-- Seed data for testing
-- Passwords: demo@example.com -> "password", test@example.com -> "test123"
INSERT OR IGNORE INTO users (id, name, email, password_hash, created_at) VALUES
    ('user-1', 'Demo User', 'demo@example.com', '$2a$10$xISzNCI9qDPuAdqvBDRVfeZ3P2vM6uNoDwfBRPhVPnQ3pEtM9fCwK', datetime('now')),
    ('user-2', 'Test User', 'test@example.com', '$2a$10$CU5wKCo4fgbJycXUtrEMKeoLt5sCB1YLoySJ26Z6wHPa7YkMcnDOm', datetime('now'));
