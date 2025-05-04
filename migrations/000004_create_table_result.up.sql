CREATE TYPE result_type AS ENUM ('good', 'bad');

CREATE TABLE result
(
    id          SERIAL PRIMARY KEY,
    exercise_id INT         NOT NULL REFERENCES exercise (id) ON DELETE CASCADE,
    type        result_type NOT NULL,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_result_exercise_id ON result (exercise_id);
