ALTER TABLE user_metadata
    DROP CONSTRAINT user_metadata_pkey,
    DROP CONSTRAINT user_metadata_id_fkey,

    DROP COLUMN id,
    DROP COLUMN is_email_verified,
    DROP COLUMN is_phone_number_verified,
    DROP COLUMN display_name,
    DROP COLUMN color,
    DROP COLUMN bio,

    ADD COLUMN id BIGSERIAL,
    ADD COLUMN user_id integer NOT NULL,
    ADD COLUMN key character varying(255) NOT NULL,
    ADD COLUMN value character varying(255) NOT NULL,

    ADD CONSTRAINT user_metadata_pkey PRIMARY KEY (id),
    ADD CONSTRAINT user_metadata_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE users
    ADD COLUMN first_name character varying(50) NOT NULL,
    ADD COLUMN last_name character varying(50) NOT NULL;
