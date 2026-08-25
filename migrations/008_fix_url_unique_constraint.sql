ALTER TABLE urls
DROP CONSTRAINT IF EXISTS unique_original_url;

ALTER TABLE urls ADD CONSTRAINT unique_user_original_url UNIQUE (user_id, original_url);