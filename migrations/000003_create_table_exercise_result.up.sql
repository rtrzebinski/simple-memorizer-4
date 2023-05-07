CREATE TABLE exercise_result
(
    id           SERIAL PRIMARY KEY,
    exercise_id  INT NOT NULL REFERENCES exercise (id) ON DELETE CASCADE,
    bad_answers  INT NOT NULL DEFAULT 0,
    good_answers INT NOT NULL DEFAULT 0,
    CONSTRAINT exercise_result_exercise_id_unique UNIQUE (exercise_id)
);
