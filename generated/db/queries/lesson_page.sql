-- name: AddLessonPage :one
INSERT INTO lesson_page (lesson_id,page) 
VALUES($1,$2)ON CONFLICT (lesson_id,page) 
DO NOTHING RETURNING *;


-- name: GetLessonPage :many
SELECT * FROM lesson_page WHERE lesson_id= $1
ORDER BY page ASC;


-- name: DeleteLessonPage :exec
DELETE FROM lesson_page WHERE page_id = $1;


-- name: GetLessonPageByCode :many
SELECT * FROM lesson_page WHERE lesson_id= $1
ORDER BY page ASC;