ALTER TABLE urls
ADD COLUMN visibility VARCHAR(10) NOT NULL DEFAULT 'private';

ALTER TABLE urls
ADD CONSTRAINT urls_visibility_check
CHECK (visibility IN ('public', 'private'));