CREATE TABLE exercise
(
    id                           SERIAL PRIMARY KEY,
    lesson_id                    INT     NOT NULL REFERENCES lesson (id) ON DELETE CASCADE,
    question                     VARCHAR NOT NULL,
    answer                       VARCHAR NOT NULL,
    bad_answers                  INT     DEFAULT 0,
    bad_answers_today            INT     DEFAULT 0,
    latest_bad_answer            TIMESTAMPTZ DEFAULT NULL,
    latest_bad_answer_was_today  BOOLEAN DEFAULT FALSE,
    good_answers                 INT     DEFAULT 0,
    good_answers_today           INT     DEFAULT 0,
    latest_good_answer           TIMESTAMPTZ DEFAULT NULL,
    latest_good_answer_was_today BOOLEAN DEFAULT FALSE,
    CONSTRAINT exercise_question_unique_per_lesson UNIQUE (lesson_id, question)
);

CREATE INDEX idx_exercise_lesson_id ON exercise (lesson_id);
