CREATE TYPE answer_type AS ENUM ('good', 'bad');

CREATE TABLE answer
(
    id          SERIAL PRIMARY KEY,
    exercise_id INT         NOT NULL REFERENCES exercise (id) ON DELETE CASCADE,
    type        answer_type NOT NULL,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_answer_exercise_id ON answer (exercise_id);
