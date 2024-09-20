-- name: AddLanguage :one
INSERT INTO lang (language,lang_code) 
VALUES($1,$2) RETURNING *;