CREATE TABLE password_reset_tokens (
    id BIGSERIAL,
    username character varying(50) NOT NULL UNIQUE,
    token bytea NOT NULL,
    expires_at timestamp without time zone NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE zk_snarks_key_pairs (
    id BIGSERIAL,
    proving_key bytea NOT NULL,
    verification_key bytea NOT NULL,

    PRIMARY KEY (id)
);