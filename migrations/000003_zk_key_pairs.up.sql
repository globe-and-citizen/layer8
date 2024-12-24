CREATE TABLE zk_snarks_key_pairs (
    id BIGSERIAL,
    proving_key bytea NOT NULL,
    verifying_key bytea NOT NULL,

    PRIMARY KEY (id)
);

ALTER TABLE users
    ADD COLUMN zk_key_pair_id integer NOT NULL;
