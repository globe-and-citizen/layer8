ALTER TABLE IF EXISTS clients
    ADD COLUMN iteration_count integer DEFAULT 4096,
    ADD COLUMN server_key VARCHAR(255),
    ADD COLUMN stored_key VARCHAR(255);