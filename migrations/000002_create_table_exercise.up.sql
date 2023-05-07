CREATE TABLE exercise
(
    id        SERIAL PRIMARY KEY,
    lesson_id INT     NOT NULL REFERENCES lesson (id) ON DELETE CASCADE,
    question  VARCHAR NOT NULL,
    answer    VARCHAR NOT NULL
);
