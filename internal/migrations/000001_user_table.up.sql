CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    age INT
    );

CREATE TABLE IF NOT EXISTS auth_users(
    id serial PRIMARY KEY,
    user_id  INT UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
    );

    ALTER TABLE
    "auth_users"
ADD
    CONSTRAINT "fk_auth_users_users" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;