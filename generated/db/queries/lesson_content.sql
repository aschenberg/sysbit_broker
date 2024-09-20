-- name: AddLessonContent :one
INSERT INTO lesson_content (
	page_id,
	lesson_id,
	step,
  	lang,
	content
) VALUES (
  $1,$2,$3,$4,$5
) RETURNING content_id;

-- -- name: GetLesson :many
-- SELECT * FROM lesson_content
-- WHERE lessoncode = $1 AND lang=$2;

-- name: UpdateLessonContent :exec
UPDATE lesson_content SET content = $2 WHERE content_id = $1;

-- name: DeleteLessonContent :exec
DELETE FROM lesson_content WHERE content_id = $1;
