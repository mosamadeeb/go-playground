-- Active: 1728656163713@@127.0.0.1@5432@chirpy
-- +goose Up
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
