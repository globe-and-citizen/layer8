CREATE TABLE clients (
    id character varying(36) NOT NULL,
    secret character varying NOT NULL,
    name character varying(255) NOT NULL,
    redirect_uri character varying(255) NOT NULL,
    username character varying(50),
    password character varying(255),
    salt character varying(255),
    backend_uri character varying(255)
);

CREATE TABLE user_metadata (
    id integer NOT NULL,
    user_id integer NOT NULL,
    key character varying(255) NOT NULL,
    value character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE users (
    id integer NOT NULL,
    email character varying(255) NOT NULL,
    username character varying(50) NOT NULL,
    password character varying(255) NOT NULL,
    first_name character varying(50) NOT NULL,
    last_name character varying(50) NOT NULL,
    salt character varying(255) DEFAULT 'ThisIsARandomSalt123!@#'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE SEQUENCE user_metadata_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE ONLY user_metadata ALTER COLUMN id SET DEFAULT nextval('user_metadata_id_seq'::regclass);
ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);

ALTER TABLE ONLY clients
    ADD CONSTRAINT clients_pkey PRIMARY KEY (id);

ALTER TABLE ONLY user_metadata
    ADD CONSTRAINT user_metadata_pkey PRIMARY KEY (id);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_username_key UNIQUE (username);

ALTER TABLE ONLY user_metadata
    ADD CONSTRAINT user_metadata_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);