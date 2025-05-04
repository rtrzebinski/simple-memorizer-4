CREATE TABLE lesson
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    user_id     INT     NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    CONSTRAINT lesson_name_unique UNIQUE (name)
);
