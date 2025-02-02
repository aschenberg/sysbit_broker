// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: lesson_page.sql

package pg

import (
	"context"
)

const addLessonPage = `-- name: AddLessonPage :one
INSERT INTO lesson_page (lesson_id,page) 
VALUES($1,$2)ON CONFLICT (lesson_id,page) 
DO NOTHING RETURNING page_id, lesson_id, page
`

type AddLessonPageParams struct {
	LessonID int32
	Page     int32
}

func (q *Queries) AddLessonPage(ctx context.Context, arg AddLessonPageParams) (LessonPage, error) {
	row := q.db.QueryRow(ctx, addLessonPage, arg.LessonID, arg.Page)
	var i LessonPage
	err := row.Scan(&i.PageID, &i.LessonID, &i.Page)
	return i, err
}

const deleteLessonPage = `-- name: DeleteLessonPage :exec
DELETE FROM lesson_page WHERE page_id = $1
`

func (q *Queries) DeleteLessonPage(ctx context.Context, pageID int32) error {
	_, err := q.db.Exec(ctx, deleteLessonPage, pageID)
	return err
}

const getLessonPage = `-- name: GetLessonPage :many
SELECT page_id, lesson_id, page FROM lesson_page WHERE lesson_id= $1
ORDER BY page ASC
`

func (q *Queries) GetLessonPage(ctx context.Context, lessonID int32) ([]LessonPage, error) {
	rows, err := q.db.Query(ctx, getLessonPage, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LessonPage{}
	for rows.Next() {
		var i LessonPage
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

const getLessonPageByCode = `-- name: GetLessonPageByCode :many
SELECT page_id, lesson_id, page FROM lesson_page WHERE lesson_id= $1
ORDER BY page ASC
`

func (q *Queries) GetLessonPageByCode(ctx context.Context, lessonID int32) ([]LessonPage, error) {
	rows, err := q.db.Query(ctx, getLessonPageByCode, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LessonPage{}
	for rows.Next() {
		var i LessonPage
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
