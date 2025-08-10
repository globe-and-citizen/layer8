ALTER TABLE IF EXISTS users
    ADD COLUMN password character varying(255);

ALTER TABLE IF EXISTS clients
    ADD COLUMN password character varying(255);
