CREATE TABLE IF NOT EXISTS users (
	id				SERIAL 			NOT NULL UNIQUE,
	chat_id			VARCHAR(30)	 	NOT NULL
);

CREATE TABLE IF NOT EXISTS lists_config (
    id              SERIAL          NOT NULL UNIQUE,
    price         	VARCHAR(8)    	NOT NULL,
    age 	        VARCHAR(4)      NOT NULL,
    countplayers    VARCHAR(6)      NOT NULL,
    timesession     VARCHAR(5)      NOT NULL,
    switch          VARCHAR(1)      NOT NULL
);

CREATE TABLE IF NOT EXISTS lists_users (
	id			SERIAL			NOT NULL UNIQUE,
	user_id	INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
	list_id INT REFERENCES lists_config (id) ON DELETE CASCADE NOT NULL
);