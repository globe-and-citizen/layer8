DROP TABLE phone_number_verification_data;

ALTER TABLE users
    DROP COLUMN phone_number_verification_code,
    DROP COLUMN phone_number_zk_proof,
    DROP COLUMN phone_number_zk_pair_id,
    DROP COLUMN telegram_session_id_hash;
