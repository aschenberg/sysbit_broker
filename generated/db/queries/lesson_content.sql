-- name: AddOrUpdateLessonContent :one
INSERT INTO lesson_content (
	lesson_id,
	page_id,
	step,
  	lang,
	content
) VALUES (
  $1,$2,$3,$4,$5
) ON CONFLICT (lesson_id,page_id,step,lang) 
DO UPDATE SET content=EXCLUDED.content RETURNING content_id,
CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation;

-- -- name: GetLesson :many
-- SELECT * FROM lesson_content
-- WHERE lessoncode = $1 AND lang=$2;

-- name: DeleteLessonContent :exec
DELETE FROM lesson_content WHERE content_id = $1;
