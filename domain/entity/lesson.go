package entity

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type AddLessonParams struct {
	LessonCode string `json:"lesson_code"`
	Title      string `json:"title"`
}

type AddLessonContentParams struct {
	PageID   int32  `json:"page_id"`
	Step     int32  `json:"step"`
	Lang     string `json:"lang"`
	Content  string `json:"content"`
	LessonID int32  `json:"lesson_id"`
}

type AddLessonPageParams struct {
	LessonID int32 `json:"lesson_id"`
	Page     int32 `json:"page"`
}

type LessonPage struct {
	PageID   int32  `json:"page_id"`
	LessonID int32  `json:"lesson_id"`
	Page     string `json:"page"`
}

type LessonContent struct {
	ContentID int32  `json:"content_id"`
	PageID    int32  `json:"page_id,omitempty"`
	Lang      string `json:"lang"`
	Step      int32  `json:"step"`
	Content   string `json:"content"`
}
type LessonDial struct {
	PrimaryID int32  `json:"primary_id"`
	PrimaryLang    string  `json:"primary_lang,omitempty"`
	SecondaryID      int32 `json:"secondary_id"`
	SecondaryLang      string  `json:"secondary_lang"`
}

type OKReply struct {
	Status  string
	Message string
}

type GetLessonDetailByCodeRow struct {
	LessonID   int32
	LessonCode string
	Title      string
	ContentID  pgtype.Int4
	LessonID_2 pgtype.Int4
	PageID     pgtype.Int4
	Lang       pgtype.Text
	Step       pgtype.Int4
	Content    pgtype.Text
	PageID_2   pgtype.Int4
	LessonID_3 pgtype.Int4
	Page       pgtype.Int4
}

type GetLessonDetail struct {
	LessonID   int32
	LessonCode string
	Title      string
	ContentID  int32
	PageID     int32
	Lang       string
	Step       string
	Content    string
	Page       string
}

type NOKReply struct {
	Status string
	Errors string
}

type AppId struct {
	AppId string `json:"appid"`
	Pin   string `json:"pin"`
}

type LessonDetails struct {
	LessonDetails []LessonDetail `json:"lessonDetail"`
}

type LessonDetail struct {
	StepCode string `json:"stepCode"`
	SpeakTxt string `json:"speakTxt"`
}

type Lesson struct {
	LessonID   int32  `json:"lesson_id"`
	LessonCode string `json:"lesson_code"`
	Title      string `json:"title"`
}

type Progress struct {
	Done []Lesson `json:"Done"`
}

// ----------------------------
type QuizData struct {
	Quizes []Quiz `json:"quizes"`
}

type Selection struct {
	Choice     string `json:"choice"`
	Desription string `json:"description"`
}

type Quiz struct {
	Nbr        int         `json:"nbr"`
	Context    string      `json:"context"`
	Question   string      `json:"question"`
	Selections []Selection `json:"selections"`
	Answer     string      `json:"answer"`
	Reason     string      `json:"reason"`
}

// -----------------------------------------------
// type LessonHeaders struct {
// 	LessonHeader []LessonHeader
// }

type LessonData struct {
	Pages []Page `json:"pages"`
}

type Page struct {
	Page    int             `json:"page"`
	Steps []Step `json:"steps"`
}

type Step struct {
	Step int `json:"step"`
	Content LessonDial `json:"content"`
	
}
