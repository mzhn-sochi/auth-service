CREATE TYPE roles AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS users
(
    id                TEXT PRIMARY KEY,
    last_name         TEXT NOT NULL,
    first_name        TEXT NOT NULL,
    middle_name       TEXT NOT NULL,
    role              roles           NOT NULL DEFAULT 'user',
    phone             CHAR(11) UNIQUE NOT NULL,
    password          TEXT            NOT NULL,
    email             TEXT UNIQUE,
    is_phone_verified BOOLEAN         NOT NULL DEFAULT FALSE,
    is_email_verified BOOLEAN         NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP
);
