-- name: AddOrUpdateQuize :one
INSERT INTO quizes (
	lesson_id,
	num,
	lang,
  	quiz,
	answer
) VALUES (
  $1,$2,$3,$4,$5
) ON CONFLICT (lesson_id,num,lang) 
DO UPDATE SET quiz=EXCLUDED.quiz, answer=EXCLUDED.answer RETURNING quiz_id,
CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation;

-- name: GetQuizByCode :many
SELECT * FROM lesson as a 
LEFT JOIN quizes as b ON a.lesson_id = b.lesson_id
WHERE a.lesson_code = $1 ORDER BY b.num ; 


-- name: DeleteQuiz :exec
DELETE FROM quizes WHERE quiz_id = $1;