-- name: AddOrUpdateLessonAsset :one
INSERT INTO lesson_assets (
	page_id,
	lesson_id,
	step,
  	type,
	url
) VALUES (
  $1,$2,$3,$4,$5
)  ON CONFLICT (lesson_id,page_id,step,type) 
DO UPDATE SET url=EXCLUDED.url RETURNING asset_id,
CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation;

-- name: GetLessonAssetByCode :many
SELECT * FROM lesson as a 
LEFT JOIN lesson_assets as b ON a.lesson_id = b.lesson_id 
LEFT JOIN lesson_page as c ON b.page_id = c.page_id 
WHERE a.lesson_code = $1 ORDER BY c.page ; 

-- name: UpdateLessonAsset :exec
UPDATE lesson_assets SET url = $2 WHERE asset_id = $1;

-- name: DeleteLessonAsset :exec
DELETE FROM lesson_assets WHERE asset_id = $1;