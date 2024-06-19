ALTER TABLE users 
    ADD COLUMN email_proof character varying(255),
    ADD COLUMN verification_code character varying(10),
    DROP COLUMN email;

CREATE TABLE email_verification_data (
     id BIGSERIAL,
     user_id integer NOT NULL,
     verification_code character varying(10) NOT NULL,
     expires_at timestamp without time zone NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);