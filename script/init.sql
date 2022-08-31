CREATE TABLE users(
    user_id SERIAL PRIMARY KEY,
    user_password VARCHAR(200) NOT NULL,
    user_name VARCHAR(200) NOT NULL
);