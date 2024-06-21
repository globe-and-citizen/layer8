ALTER TABLE users 
    ADD COLUMN email character varying(255) NOT NULL,
    DROP COLUMN email_proof,
    DROP COLUMN verification_code;

DROP TABLE email_verification_data;
