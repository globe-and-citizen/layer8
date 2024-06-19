CREATE TABLE clients (
    id character varying(36) NOT NULL,
    secret character varying NOT NULL,
    name character varying(255) NOT NULL,
    redirect_uri character varying(255) NOT NULL,
    username character varying(50) UNIQUE NOT NULL,
    password character varying(255) NOT NULL,
    salt character varying(255),
    backend_uri character varying(255) UNIQUE NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE users (
    id BIGSERIAL,
    username character varying(50) NOT NULL UNIQUE,
    password character varying(255) NOT NULL,
    first_name character varying(50) NOT NULL,
    last_name character varying(50) NOT NULL,
    salt character varying(255) DEFAULT 'ThisIsARandomSalt123!@#'::character varying NOT NULL,
    email_proof character varying(255) NOT NULL,
    verification_code character varying(10) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE user_metadata (
    id BIGSERIAL,
    user_id integer NOT NULL,
    key character varying(255) NOT NULL,
    value character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE email_verification_data (
     id BIGSERIAL,
     user_id integer NOT NULL,
     verification_code character varying(10) NOT NULL,
     expires_at timestamp without time zone NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
