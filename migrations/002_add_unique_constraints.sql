ALTER TABLE urls
ADD CONSTRAINT unique_short_id
UNIQUE (short_id);

ALTER TABLE urls 
ADD CONSTRAINT unique_original_url
UNIQUE (original_url);