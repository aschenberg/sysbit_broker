-- name: GetLessonHeader :many
SELECT * FROM lesson
ORDER BY lesson_code ASC;


-- name: GetLessonDetailByCode :many
SELECT * FROM lesson as a 
LEFT JOIN lesson_content as b ON a.lesson_id = b.lesson_id 
LEFT JOIN lesson_page as c ON b.page_id = c.page_id 
WHERE a.lesson_code = $1 AND b.lang = ANY (sqlc.slice('lang')) ORDER BY c.page ; 


-- name: AddLesson :one
INSERT INTO lesson (lesson_code,title) 
VALUES($1,$2) RETURNING *;


-- name: DeleteLesson :exec
DELETE FROM lesson WHERE lesson_id = $1;