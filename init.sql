CREATE DATABASE hg_bot_config;

CREATE TABLE users (
	id				SERIAL 			NOT NULL UNIQUE,
	username 		VARCHAR(255) 	NOT NULL,
	password_hash 	BYTEA 			NOT NULL
);

CREATE TABLE lists_config (
    id              SERIAL          NOT NULL UNIQUE,
    keyword         VARCHAR(255)    NOT NULL,
    ageMin          VARCHAR(4)      NOT NULL,
    countPlayers    VARCHAR(6)      NOT NULL,
    timeSession     VARCHAR(5)      NOT NULL
);

CREATE TABLE lists_users (
	id			SERIAL			NOT NULL UNIQUE,
	user_id	INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
	list_id INT REFERENCES lists_config (id) ON DELETE CASCADE NOT NULL
);