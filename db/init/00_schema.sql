CREATE DATABASE rakutan;

\c rakutan

CREATE TYPE user_action AS ENUM ('search', 'rakutan', 'onitan', 'set_fav', 'unset_fav', 'get_fav', 'info', 'help');

CREATE TABLE users (
    uid VARCHAR(64) NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    registered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    verified_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (uid)
);

CREATE TABLE user_logs (
    uid VARCHAR(64) NOT NULL,
    action user_action NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (uid, timestamp)
);

CREATE TABLE verification_tokens (
    uid VARCHAR(64) NOT NULL,
    token VARCHAR(128) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (uid)
);

CREATE TABLE favorites (
    uid VARCHAR(64) NOT NULL,
    id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (uid, id)
);

CREATE TABLE rakutan2021 (
    id INTEGER NOT NULL,
    faculty_name VARCHAR(64) NOT NULL,
    lecture_name VARCHAR(64) NOT NULL,
    register SMALLINT[],
    passed SMALLINT[],
    PRIMARY KEY (id)
);