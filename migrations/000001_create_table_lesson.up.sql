CREATE TABLE lesson
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    CONSTRAINT lesson_name_unique UNIQUE (name)
);
