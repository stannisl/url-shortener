CREATE TABLE IF NOT EXISTS urls (
                                    id BIGSERIAL PRIMARY KEY,
                                    short_code VARCHAR(10) UNIQUE NOT NULL,
                                    original_url TEXT NOT NULL,
                                    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_urls_short_code
    ON urls(short_code);

CREATE TABLE IF NOT EXISTS url_analytics (
                                             id BIGSERIAL PRIMARY KEY,
                                             url_id BIGINT REFERENCES urls(id) ON DELETE CASCADE,
                                             user_agent TEXT,
                                             ip_address INET,
                                             accessed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_url_analytics_url_id
    ON url_analytics(url_id);

CREATE INDEX idx_url_analytics_accessed_at
    ON url_analytics(accessed_at);