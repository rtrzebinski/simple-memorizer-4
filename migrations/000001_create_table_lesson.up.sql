CREATE TABLE lesson
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR NOT NULL,
    description    VARCHAR NOT NULL,
    exercise_count INT DEFAULT 0,
    CONSTRAINT lesson_name_unique UNIQUE (name)
);
