ALTER TABLE urls
ADD COLUMN user_id UUID;

ALTER TABLE urls ADD CONSTRAINT fk_urls_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

CREATE INDEX idx_urls_user_id on urls (user_id)