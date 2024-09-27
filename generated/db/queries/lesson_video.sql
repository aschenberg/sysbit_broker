-- name: AddOrUpdateLessonAudio :one
INSERT INTO lesson_video (
  lesson_id,
	type,
	url
) VALUES (
  $1,$2,$3
) ON CONFLICT (lesson_id,type) DO UPDATE SET 
url=EXCLUDED.url RETURNING *,
CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation;

-- name: GetLessonVideo :many
SELECT * FROM lesson as a 
LEFT JOIN lesson_video as b ON a.lesson_id = b.lesson_id
WHERE a.lesson_code = $1;


-- name: DeleteLessonVideo :exec
DELETE FROM lesson_video WHERE video_id = $1;
