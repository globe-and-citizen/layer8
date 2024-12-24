DROP TABLE zk_snarks_key_pairs;

ALTER TABLE users
    DROP COLUMN zk_key_pair_id;
