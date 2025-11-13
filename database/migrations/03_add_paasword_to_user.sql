-- +goose Up
ALTER TABLE Users ADD password TEXT NOT NULL DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique
ON Users(email);