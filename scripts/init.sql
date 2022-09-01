CREATE TABLE users(
    u_id SERIAL PRIMARY KEY,
    u_password VARCHAR(200) NOT NULL,
    u_name VARCHAR(200) NOT NULL
);