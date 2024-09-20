-- name: UpdateUserToken :exec
UPDATE users SET refresh_token = $2 WHERE user_id = $1;


-- name: CreateOrUpdateUser :one
INSERT INTO users (
  user_id,
	app_id,
	email,
  picture,
	role,
	name,
	refresh_token,
	is_deleted,
	is_active,
	created_at,
	updated_at
) VALUES (
  @userid,@appid,@mail,@pic,@roles,@name,@refreshtoken,@isdeleted,@isactive,@createdat,@updatedat
) ON CONFLICT (app_id) DO UPDATE SET 
picture = @picture,
name = @name
RETURNING *, CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation;


-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetRefreshToken :one
SELECT refresh_token FROM users
WHERE user_id = $1 LIMIT 1;
