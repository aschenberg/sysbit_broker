package repository

import (
	"context"

	"sysbitBroker/domain/entity"
	"sysbitBroker/domain/param"
	"sysbitBroker/pkg"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ILessonRepository interface {
	CreateOrUpdateLesson(c context.Context, arg param.AddLessonParams) (entity.AddOrUpdateLessonRow, error)
	GetLesson(c context.Context) ([]entity.Lesson, error)
	GetLessonByCode(c context.Context, lessonCode string) (entity.Lesson, error)
	DeleteLesson(ctx context.Context, lessonID int32) error
	CreateLessonPage(c context.Context, arg param.AddLessonPageParams) (entity.LessonPage, error)
	GetLessonPage(c context.Context, lessonID int32) ([]entity.LessonPage, error)
	DeleteLessonPage(ctx context.Context, pageID int32) error
	GetLessonDetail(c context.Context, lessonCode string, lang []string) ([]entity.GetLessonDetailByCodeRow, error)
	CreateOrUpdateLessonContent(c context.Context, arg param.AddLessonContentParams) (entity.AddOrUpdateLessonContentRow, error)
	DeleteLessonContent(c context.Context, contentID int32) error
	CreateOrUpdateLessonAsset(c context.Context, arg param.AddOrUpdateLessonAssetParams) (entity.AddOrUpdateLessonAssetRow, error)
	GetLessonAsset(c context.Context, lessonCode string) ([]entity.GetLessonAssetByCodeRow, error)
	DeleteLessonAsset(c context.Context, assetID int32) error
	CreateOrUpdateLessonVideo(c context.Context, arg param.AddOrUpdateLessonAudioParams) (entity.AddOrUpdateLessonAudioRow, error)
	GetLessonVideo(c context.Context, lessonCode string) ([]entity.GetLessonVideoRow, error)
	DeleteLessonVideo(c context.Context, videoID int32) error
	CreateOrUpdateQuiz(c context.Context, arg param.AddOrUpdateQuizeParams) (entity.AddOrUpdateQuizeRow, error)
	GetQuizByCode(c context.Context, lessonCode string) ([]entity.GetQuizByCodeRow, error)
	DeleteQuiz(c context.Context, quizID int32) error
}

type lessonRepository struct {
	db *pgxpool.Pool
}

func NewLessonRepository(pg *pkg.Postgres) ILessonRepository {
	return &lessonRepository{
		db: pg.Pool,
	}
}

func (rp *lessonRepository) CreateOrUpdateLesson(c context.Context, arg param.AddLessonParams) (entity.AddOrUpdateLessonRow, error) {
	const addOrUpdateLesson = `-- name: AddOrUpdateLesson :one
INSERT INTO lesson (lesson_code,title,image) 
VALUES($1,$2,$3) ON CONFLICT (lesson_code) 
DO UPDATE SET title=EXCLUDED.title RETURNING lesson_id, lesson_code, title, image,
CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation
`

	row := rp.db.QueryRow(c, addOrUpdateLesson, arg.LessonCode, arg.Title)
	var i entity.AddOrUpdateLessonRow
	err := row.Scan(
		&i.LessonID,
		&i.LessonCode,
		&i.Title,
		&i.Image,
		&i.Operation,
	)
	return i, err
}

func (rp *lessonRepository) GetLesson(c context.Context) ([]entity.Lesson, error) {
	const getLessonHeader = `-- name: GetLessonHeader :many
SELECT lesson_id, lesson_code, title, image FROM lesson
ORDER BY lesson_code ASC
`
	rows, err := rp.db.Query(c, getLessonHeader)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []entity.Lesson{}
	for rows.Next() {
		var i entity.Lesson
		if err := rows.Scan(
			&i.LessonID,
			&i.LessonCode,
			&i.Title,
			&i.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (rp *lessonRepository) GetLessonByCode(c context.Context, lessonCode string) (entity.Lesson, error) {
	const getLessonHeaderByCode = `-- name: GetLessonHeaderByCode :one
SELECT lesson_id, lesson_code, title FROM lesson WHERE lesson_code = $1
`
	row := rp.db.QueryRow(c, getLessonHeaderByCode, lessonCode)
	var i entity.Lesson
	err := row.Scan(&i.LessonID, &i.LessonCode, &i.Title)
	return i, err
}

func (rp *lessonRepository) DeleteLesson(ctx context.Context, lessonID int32) error {
	const deleteLesson = `-- name: DeleteLesson :exec
	DELETE FROM lesson WHERE lesson_id = $1
	`
	_, err := rp.db.Exec(ctx, deleteLesson, lessonID)
	return err
}

func (rp *lessonRepository) CreateOrUpdateLessonContent(c context.Context, arg param.AddLessonContentParams) (entity.AddOrUpdateLessonContentRow, error) {
	const addOrUpdateLessonContent = `-- name: AddOrUpdateLessonContent :one
	INSERT INTO lesson_content (
	lesson_id,
	page_id,
	step,
  	lang,
	content) 
	VALUES ($1,$2,$3,$4,$5) ON CONFLICT (lesson_id,page_id,step,lang) 
	DO UPDATE SET content=EXCLUDED.content RETURNING content_id,
	CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation
	`
	row := rp.db.QueryRow(c, addOrUpdateLessonContent,
		arg.LessonID,
		arg.PageID,
		arg.Step,
		arg.Lang,
		arg.Content,
	)
	var i entity.AddOrUpdateLessonContentRow
	err := row.Scan(&i.ContentID, &i.Operation)
	return i, err
}

func (rp *lessonRepository) CreateLessonPage(c context.Context, arg param.AddLessonPageParams) (entity.LessonPage, error) {
	const addLessonPage = `-- name: AddLessonPage :one
	INSERT INTO lesson_page (lesson_id,page) 
	VALUES($1,$2)ON CONFLICT (lesson_id,page) 
	DO NOTHING RETURNING page_id, lesson_id, page
	`
	row := rp.db.QueryRow(c, addLessonPage, arg.LessonID, arg.Page)
	var i entity.LessonPage
	err := row.Scan(&i.PageID, &i.LessonID, &i.Page)
	return i, err
}

func (rp *lessonRepository) GetLessonPage(c context.Context, lessonID int32) ([]entity.LessonPage, error) {
	const getLessonPage = `-- name: GetLessonPage :many
	SELECT page_id, lesson_id, page FROM lesson_page WHERE lesson_id= $1
	ORDER BY page ASC
	`
	rows, err := rp.db.Query(c, getLessonPage, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []entity.LessonPage{}
	for rows.Next() {
		var i entity.LessonPage
		if err := rows.Scan(&i.PageID, &i.LessonID, &i.Page); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
func (rp *lessonRepository) DeleteLessonPage(ctx context.Context, pageID int32) error {
	const deleteLessonPage = `-- name: DeleteLessonPage :exec
	DELETE FROM lesson_page WHERE page_id = $1
	`
	_, err := rp.db.Exec(ctx, deleteLessonPage, pageID)
	return err
}
func (rp *lessonRepository) GetLessonDetail(c context.Context, lessonCode string, lang []string) ([]entity.GetLessonDetailByCodeRow, error) {
	const getLessonDetailByCode = `-- name: GetLessonDetailByCode :many
SELECT a.lesson_id, lesson_code, title, image, content_id, b.lesson_id, b.page_id, lang, step, content, c.page_id, c.lesson_id, page FROM lesson as a 
LEFT JOIN lesson_content as b ON a.lesson_id = b.lesson_id 
LEFT JOIN lesson_page as c ON b.page_id = c.page_id 
WHERE a.lesson_code = $1 AND b.lang = ANY ($2) ORDER BY c.page
`
	rows, err := rp.db.Query(c, getLessonDetailByCode, lessonCode, lang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []entity.GetLessonDetailByCodeRow{}
	for rows.Next() {
		var i entity.GetLessonDetailByCodeRow
		if err := rows.Scan(
			&i.LessonID,
			&i.LessonCode,
			&i.Title,
			&i.Image,
			&i.ContentID,
			&i.LessonID_2,
			&i.PageID,
			&i.Lang,
			&i.Step,
			&i.Content,
			&i.PageID_2,
			&i.LessonID_3,
			&i.Page,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (rp *lessonRepository) DeleteLessonContent(c context.Context, contentID int32) error {
	const deleteLessonContent = `-- name: DeleteLessonContent :exec
DELETE FROM lesson_content WHERE content_id = $1
`
	_, err := rp.db.Exec(c, deleteLessonContent, contentID)
	return err
}

func (rp *lessonRepository) CreateOrUpdateLessonAsset(c context.Context, arg param.AddOrUpdateLessonAssetParams) (entity.AddOrUpdateLessonAssetRow, error) {
	const addOrUpdateLessonAsset = `-- name: AddOrUpdateLessonAsset :one
INSERT INTO lesson_assets (
	page_id,
	lesson_id,
	step,
  	type,
	url
) VALUES (
  $1,$2,$3,$4,$5
)  ON CONFLICT (lesson_id,page_id,step,type) DO UPDATE SET url=EXCLUDED.url RETURNING asset_id
`
	row := rp.db.QueryRow(c, addOrUpdateLessonAsset,
		arg.PageID,
		arg.LessonID,
		arg.Step,
		arg.Type,
		arg.Url,
	)
	var i entity.AddOrUpdateLessonAssetRow
	err := row.Scan(&i.AssetID, &i.Operation)
	return i, err
}

func (rp *lessonRepository) GetLessonAsset(c context.Context, lessonCode string) ([]entity.GetLessonAssetByCodeRow, error) {
	const getLessonAssetByCode = `-- name: GetLessonAssetByCode :many
SELECT a.lesson_id, lesson_code, title, asset_id, b.lesson_id, b.page_id, step, type, url, c.page_id, c.lesson_id, page FROM lesson as a 
LEFT JOIN lesson_assets as b ON a.lesson_id = b.lesson_id 
LEFT JOIN lesson_page as c ON b.page_id = c.page_id 
WHERE a.lesson_code = $1 ORDER BY c.page
`
	rows, err := rp.db.Query(c, getLessonAssetByCode, lessonCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []entity.GetLessonAssetByCodeRow{}
	for rows.Next() {
		var i entity.GetLessonAssetByCodeRow
		if err := rows.Scan(
			&i.LessonID,
			&i.LessonCode,
			&i.Title,
			&i.AssetID,
			&i.LessonID_2,
			&i.PageID,
			&i.Step,
			&i.Type,
			&i.Url,
			&i.PageID_2,
			&i.LessonID_3,
			&i.Page,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (rp *lessonRepository) DeleteLessonAsset(c context.Context, assetID int32) error {
	const deleteLessonAsset = `-- name: DeleteLessonAsset :exec
DELETE FROM lesson_assets WHERE asset_id = $1
`
	_, err := rp.db.Exec(c, deleteLessonAsset, assetID)
	return err
}

func (rp *lessonRepository) CreateOrUpdateLessonVideo(c context.Context, arg param.AddOrUpdateLessonAudioParams) (entity.AddOrUpdateLessonAudioRow, error) {
	const addOrUpdateLessonAudio = `-- name: AddOrUpdateLessonAudio :one
INSERT INTO lesson_video (
  lesson_id,
	type,
	url
) VALUES (
  $1,$2,$3
) ON CONFLICT (lesson_id,type) DO UPDATE SET 
url=EXCLUDED.url RETURNING video_id, lesson_id, type, url,
CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation
`

	row := rp.db.QueryRow(c, addOrUpdateLessonAudio, arg.LessonID, arg.Type, arg.Url)
	var i entity.AddOrUpdateLessonAudioRow
	err := row.Scan(
		&i.VideoID,
		&i.LessonID,
		&i.Type,
		&i.Url,
		&i.Operation,
	)
	return i, err
}

func (rp *lessonRepository) GetLessonVideo(c context.Context, lessonCode string) ([]entity.GetLessonVideoRow, error) {
	const getLessonVideo = `-- name: GetLessonVideo :many
SELECT a.lesson_id, lesson_code, title, video_id, b.lesson_id, type, url FROM lesson as a 
LEFT JOIN lesson_video as b ON a.lesson_id = b.lesson_id
WHERE a.lesson_code = $1
`
	rows, err := rp.db.Query(c, getLessonVideo, lessonCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []entity.GetLessonVideoRow{}
	for rows.Next() {
		var i entity.GetLessonVideoRow
		if err := rows.Scan(
			&i.LessonID,
			&i.LessonCode,
			&i.Title,
			&i.VideoID,
			&i.LessonID_2,
			&i.Type,
			&i.Url,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (rp *lessonRepository) DeleteLessonVideo(c context.Context, videoID int32) error {
	const deleteLessonVideo = `-- name: DeleteLessonVideo :exec
DELETE FROM lesson_video WHERE video_id = $1
`
	_, err := rp.db.Exec(c, deleteLessonVideo, videoID)
	return err
}

func (rp *lessonRepository) CreateOrUpdateQuiz(c context.Context, arg param.AddOrUpdateQuizeParams) (entity.AddOrUpdateQuizeRow, error) {
	const addOrUpdateQuize = `-- name: AddOrUpdateQuize :one
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
CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation
`
	row := rp.db.QueryRow(c, addOrUpdateQuize,
		arg.LessonID,
		arg.Num,
		arg.Lang,
		arg.Quiz,
		arg.Answer,
	)
	var i entity.AddOrUpdateQuizeRow
	err := row.Scan(&i.QuizID, &i.Operation)
	return i, err
}

func (rp *lessonRepository) GetQuizByCode(c context.Context, lessonCode string) ([]entity.GetQuizByCodeRow, error) {
	const getQuizByCode = `-- name: GetQuizByCode :many
SELECT a.lesson_id, lesson_code, title, image, quiz_id, b.lesson_id, num, lang, quiz, answer FROM lesson as a 
LEFT JOIN quizes as b ON a.lesson_id = b.lesson_id
WHERE a.lesson_code = $1 ORDER BY b.num
`
	rows, err := rp.db.Query(c, getQuizByCode, lessonCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []entity.GetQuizByCodeRow{}
	for rows.Next() {
		var i entity.GetQuizByCodeRow
		if err := rows.Scan(
			&i.LessonID,
			&i.LessonCode,
			&i.Title,
			&i.Image,
			&i.QuizID,
			&i.LessonID_2,
			&i.Num,
			&i.Lang,
			&i.Quiz,
			&i.Answer,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (rp *lessonRepository) DeleteQuiz(c context.Context, quizID int32) error {
	const deleteQuiz = `-- name: DeleteQuiz :exec
DELETE FROM quizes WHERE quiz_id = $1
`
	_, err := rp.db.Exec(c, deleteQuiz, quizID)
	return err
}
