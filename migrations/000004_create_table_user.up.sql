CREATE TABLE "user"
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "user"
    ADD CONSTRAINT user_name_unique UNIQUE (name);
ALTER TABLE "user"
    ADD CONSTRAINT user_email_unique UNIQUE (email);

CREATE INDEX idx_user_email ON "user" (email);
