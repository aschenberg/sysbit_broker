CREATE TABLE IF NOT EXISTS quizes (
	quiz_id INTEGER PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    lesson_id INTEGER NOT NULL,
    num INTEGER NOT NULL,
	lang VARCHAR NOT NULL,
    quiz JSONB,
    answer JSONB
);


ALTER TABLE quizes
ADD CONSTRAINT lesson_num_lang_unique UNIQUE (lesson_id,num,lang);