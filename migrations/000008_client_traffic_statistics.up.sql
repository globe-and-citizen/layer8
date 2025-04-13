CREATE TABLE client_traffic_statistics (
    id BIGSERIAL,
    client_id character varying(36) NOT NULL,
    rate_per_byte integer NOT NULL,
    total_usage_bytes integer NOT NULL,
    unpaid_amount integer NOT NULL,
    last_traffic_update_timestamp timestamp without time zone NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (client_id) REFERENCES clients(id)
)