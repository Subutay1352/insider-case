DO $$
DECLARE
    table_name CONSTANT TEXT := 'messages';
    status_queued CONSTANT TEXT := 'queued';
BEGIN
    EXECUTE format('CREATE TABLE IF NOT EXISTS %I (
        id SERIAL PRIMARY KEY,
        "to" VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        status VARCHAR(20) NOT NULL DEFAULT %L,
        message_id VARCHAR(255),
        retry_count INT DEFAULT 0 NOT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    )', table_name, status_queued);

    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%I_status ON %I(status)', table_name, table_name);
    EXECUTE format('CREATE INDEX IF NOT EXISTS idx_%I_status_created ON %I(status, created_at) WHERE status = %L', table_name, table_name, status_queued);

    EXECUTE format('INSERT INTO %I ("to", content, status, message_id) VALUES 
        (%L, %L, %L, %L),
        (%L, %L, %L, %L),
        (%L, %L, %L, %L),
        (%L, %L, %L, %L),
        (%L, %L, %L, %L)', 
        table_name,
        '+905551111111', 'Merhaba, bu bir test mesajıdır.', status_queued, 'msg-001',
        '+905552222222', 'İkinci test mesajı - scheduler tarafından gönderilecek.', status_queued, 'msg-002',
        '+905553333333', 'Üçüncü test mesajı - webhook.site üzerinden test edilecek.', status_queued, 'msg-003',
        '+905554444444', 'Dördüncü test mesajı - retry mekanizması testi için.', status_queued, 'msg-004',
        '+905555555555', 'Beşinci test mesajı - production field test.', status_queued, 'msg-005');
END $$;