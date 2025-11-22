-- Migration: Create messages table
-- Description: Creates the messages table for storing message data
-- Note: This migration works for both PostgreSQL and SQLite
-- For PostgreSQL: Uses SERIAL and TIMESTAMPTZ
-- For SQLite: GORM AutoMigrate handles type conversion automatically

-- PostgreSQL version
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    "to" VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'queued',
    message_id VARCHAR(255),
    retry_count INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_messages_status ON messages(status);
CREATE INDEX IF NOT EXISTS idx_messages_status_created ON messages(status, created_at) WHERE status = 'queued';


INSERT INTO messages ("to", content, status, message_id) VALUES 
('+905551111111', 'Merhaba, bu bir test mesajıdır.', 'queued', 'msg-001'),
('+905552222222', 'İkinci test mesajı - scheduler tarafından gönderilecek.', 'queued', 'msg-002'),
('+905553333333', 'Üçüncü test mesajı - webhook.site üzerinden test edilecek.', 'queued', 'msg-003'),
('+905554444444', 'Dördüncü test mesajı - retry mekanizması testi için.', 'queued', 'msg-004'),
('+905555555555', 'Beşinci test mesajı - production field test.', 'queued', 'msg-005');