package repository

import (
	"context"
	"sysbitBroker/config"
	"sysbitBroker/domain/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ILessonRepository interface {
	CreateLessonPage(c context.Context, arg entity.AddLessonPageParams) (entity.LessonPage, error)
	GetLessonPage(c context.Context, lessonID int32) ([]entity.LessonPage, error)
	DeleteLessonPage(ctx context.Context, pageID int32) error
	CreateLesson(c context.Context, arg entity.AddLessonParams) (entity.Lesson, error)
	GetLesson(c context.Context) ([]entity.Lesson, error)
	GetLessonByCode(c context.Context, lessonCode string) (entity.Lesson, error)
	DeleteLesson(ctx context.Context, lessonID int32) error
	GetLessonDetail(c context.Context, lessonCode string, lang []string) ([]entity.GetLessonDetailByCodeRow, error)
	CreateLessonContent(c context.Context, arg entity.AddLessonContentParams) (int32, error)
	UpdateLessonContent(c context.Context, contentID int32, content string) (error)
	DeleteLessonContent(c context.Context, contentID int32) (error)
}

type lessonRepository struct {
	db *pgxpool.Pool
}

func NewLessonRepository(pg *config.Postgres) ILessonRepository {
	return &lessonRepository{
		db: pg.Pool,
	}
}

func (rp *lessonRepository) CreateLesson(c context.Context, arg entity.AddLessonParams) (entity.Lesson, error) {
	const addLesson = `-- name: AddLesson :one
	INSERT INTO lesson (lesson_code,title) 
	VALUES($1,$2) RETURNING lesson_id, lesson_code, title
	`
	row := rp.db.QueryRow(c, addLesson, arg.LessonCode, arg.Title)
	var i entity.Lesson
	err := row.Scan(&i.LessonID, &i.LessonCode, &i.Title)
	return i, err
}

func (rp *lessonRepository) GetLesson(c context.Context) ([]entity.Lesson, error) {
	const getLessonHeader = `-- name: GetLessonHeader :many
	SELECT lesson_id, lesson_code, title FROM lesson
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
		if err := rows.Scan(&i.LessonID, &i.LessonCode, &i.Title); err != nil {
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

func (rp *lessonRepository) CreateLessonContent(c context.Context, arg entity.AddLessonContentParams) (int32, error) {
	const addLessonContent = `-- name: AddLessonContent :one
	INSERT INTO lesson_content (
	page_id,
	lesson_id,
	step,
  	lang,
	content) 
	VALUES ($1,$2,$3,$4,$5) RETURNING content_id`
	row := rp.db.QueryRow(c, addLessonContent,
		arg.PageID,
		arg.LessonID,
		arg.Step,
		arg.Lang,
		arg.Content,
	)
	var content_id int32
	err := row.Scan(&content_id)
	return content_id, err
}

func (rp *lessonRepository) CreateLessonPage(c context.Context, arg entity.AddLessonPageParams) (entity.LessonPage, error) {
	const addLessonPage = `-- name: AddLessonPage :one
	INSERT INTO lesson_page (lesson_id,page) 
	VALUES($1,$2) RETURNING page_id, lesson_id, page
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
SELECT a.lesson_id, lesson_code, title, content_id, b.lesson_id, b.page_id, lang, step, content, c.page_id, c.lesson_id, page FROM lesson as a 
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


func (rp *lessonRepository) UpdateLessonContent(c context.Context, contentID int32, content string) (error){
	const updateLessonContent = `-- name: UpdateLessonContent :exec
UPDATE lesson_content SET content = $2 WHERE content_id = $1
`
_, err := rp.db.Exec(c, updateLessonContent, contentID, content)
	return err

}

func (rp *lessonRepository) DeleteLessonContent(c context.Context, contentID int32) (error){
	const deleteLessonContent = `-- name: DeleteLessonContent :exec
DELETE FROM lesson_content WHERE content_id = $1
`
_, err := rp.db.Exec(c, deleteLessonContent, contentID)
	return err
}
