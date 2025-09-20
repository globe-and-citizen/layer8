CREATE TABLE phone_number_verification_data (
    id BIGSERIAL,
    user_id integer NOT NULL,
    verification_code character varying(10) NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    zk_proof bytea NOT NULL,
    phone_number_zk_pair_id integer NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

ALTER TABLE users
    ADD COLUMN phone_number_verification_code character varying(10),
    ADD COLUMN phone_number_zk_proof bytea,
    ADD COLUMN phone_number_zk_pair_id integer,
    ADD COLUMN telegram_session_id_hash bytea;
