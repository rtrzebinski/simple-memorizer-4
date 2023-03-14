CREATE TABLE exercise_result
(
    id SERIAL PRIMARY KEY,
    exercise_id INT NOT NULL,
    bad_answers INT NOT NULL DEFAULT 0,
    good_answers INT NOT NULL DEFAULT 0
);
