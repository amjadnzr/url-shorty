-- +goose up
ALTER TABLE Users RENAME COLUMN password TO passwordHash;