ALTER TABLE user_metadata
    DROP CONSTRAINT user_metadata_pkey,
    DROP CONSTRAINT user_metadata_user_id_fkey,

    DROP COLUMN id,
    DROP COLUMN user_id,
    DROP COLUMN key,
    DROP COLUMN value,

    ADD COLUMN id integer NOT NULL,
    ADD COLUMN is_email_verified boolean NOT NULL,
    ADD COLUMN is_phone_number_verified boolean NOT NULL,
    ADD COLUMN bio character varying(350) NOT NULL,
    ADD COLUMN display_name character varying(50) NOT NULL,
    ADD COLUMN color character varying(50) NOT NULL,

    ADD CONSTRAINT user_metadata_pkey PRIMARY KEY (id),
    ADD CONSTRAINT user_metadata_id_fkey FOREIGN KEY (id) REFERENCES users(id);

ALTER TABLE users
    DROP COLUMN first_name,
    DROP COLUMN last_name;
