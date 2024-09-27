package entity

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/jackc/pgx/v5/pgtype"
)

// Lesson
// ==========================================================================
type GetLessonDetailByCodeRow struct {
	LessonID   int32
	LessonCode string
	Title      string
	Image      string
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

type LessonData struct {
	SceneUrl string `json:"scene_url"`
	Pages    []Page `json:"pages"`
}

type Page struct {
	Page  int    `json:"page"`
	Steps []Step `json:"steps"`
}

type Step struct {
	Step     int        `json:"step"`
	Contents LessonDial `json:"contents"`
	Assets   StepAsset  `json:"assets"`
}

type Lesson struct {
	LessonID   int32  `json:"lesson_id"`
	LessonCode string `json:"lesson_code"`
	Title      string `json:"title"`
	Image      string `json:"image"`
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

type AddOrUpdateLessonRow struct {
	LessonID   int32
	LessonCode string
	Title      string
	Image      string
	Operation  string
}

// LessonPage
// =============================================================================
type LessonPage struct {
	PageID   int32  `json:"page_id"`
	LessonID int32  `json:"lesson_id"`
	Page     string `json:"page"`
}

//LessonContent
//============================================================================

type LessonContent struct {
	ContentID int32  `json:"content_id"`
	PageID    int32  `json:"page_id,omitempty"`
	Lang      string `json:"lang"`
	Step      int32  `json:"step"`
	Content   string `json:"content"`
}
type LessonDial struct {
	PrimaryID     int32  `json:"primary_id"`
	PrimaryLang   string `json:"primary_lang,omitempty"`
	SecondaryID   int32  `json:"secondary_id"`
	SecondaryLang string `json:"secondary_lang"`
}
type AddOrUpdateLessonContentRow struct {
	ContentID int32  `json:"content_id"`
	Operation string `json:"operation"`
}

type AddOrUpdateLessonAudioRow struct {
	VideoID   int32  `json:"video_id"`
	LessonID  int32  `json:"lesson_id"`
	Type      string `json:"type"`
	Url       string `json:"url"`
	Operation string `json:"operation"`
}

// LessonAssets
// ==========================================================================
type AddOrUpdateLessonAssetRow struct {
	AssetID   int32  `json:"asset_id"`
	Operation string `json:"operation"`
}

type GetLessonAssetByCodeRow struct {
	LessonID   int32
	LessonCode string
	Title      string
	AssetID    pgtype.Int4
	LessonID_2 pgtype.Int4
	PageID     pgtype.Int4
	Step       pgtype.Int4
	Type       pgtype.Text
	Url        pgtype.Text
	PageID_2   pgtype.Int4
	LessonID_3 pgtype.Int4
	Page       pgtype.Int4
}

type LessonAsset struct {
	LessonID   int32  `json:"lesson_id"`
	LessonCode string `json:"lesson_code"`
	AssetID    int32  `json:"asset_id"`
	PageID     int32  `json:"page_id"`
	Step       int32  `json:"step"`
	Type       string `json:"type"`
	Url        string `json:"url"`
	Page       int32  `json:"page"`
}

type LessonAssetUrl struct {
	Url string `json:"url"`
}

type StepAsset struct {
	Image    string `json:"image"`
	VoiceUrl string `json:"voice"`
}

// LessonVideo
// =============================================================================
type GetLessonVideoRow struct {
	LessonID   int32
	LessonCode string
	Title      string
	VideoID    pgtype.Int4
	LessonID_2 pgtype.Int4
	Type       pgtype.Text
	Url        pgtype.Text
}

type LessonVideo struct {
	LessonID   int32  `json:"lesson_id"`
	LessonCode string `json:"lesson_code"`
	VideoID    int32  `json:"video_id"`
	Type       string `json:"type"`
	Url        string `json:"url"`
}

// Quiz
// ==============================================================================
type GetQuizByCodeRow struct {
	LessonID   int32
	LessonCode string
	Title      string
	Image      string
	QuizID     pgtype.Int4
	LessonID_2 pgtype.Int4
	Num        pgtype.Int4
	Lang       pgtype.Text
	Quiz       []byte
	Answer     []byte
}

type AddOrUpdateQuizeRow struct {
	QuizID    int32  `json:"quiz_id"`
	Operation string `json:"operation"`
}

type Selection struct {
	Choice     string `json:"choice"`
	Desription string `json:"description"`
}

type QuizQuestion struct {
	Context    string      `json:"context"`
	Question   string      `json:"question"`
	Selections []Selection `json:"selections"`
}

type QuizAnswer struct {
	Answer string `json:"answer"`
	Reason string `json:"reason"`
}

type QuizData struct {
	LessonCode string `json:"lesson_code"`
	QuizID     int32  `json:"quiz_id"`
	Num        int32  `json:"num"`
	Lang       string `json:"lang"`
	QuizQuestion
	QuizAnswer
}

type Progress struct {
	Done []Lesson `json:"Done"`
}

//Func Resp
//======================================================================

func QuizResp(val []GetQuizByCodeRow) []QuizData {
	var result []QuizData

	for _, v := range val {
		result = append(result, QuizData{
			LessonCode:   v.LessonCode,
			QuizID:       v.QuizID.Int32,
			Num:          v.Num.Int32,
			Lang:         v.Lang.String,
			QuizQuestion: ByteToJson[QuizQuestion](v.Quiz),
			QuizAnswer:   ByteToJson[QuizAnswer](v.Answer),
		})
	}
	return result
}

func StructToBytes[T any](v T) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println("error marshaling struct to bytes:", err)

	}
	return data
}

func ByteToJson[T any](data []byte) T {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Println("error marshaling", err)
	}
	return result
}

func LessonResponse(val []GetLessonAssetByCodeRow, baseurl string, path string) []LessonAsset {
	var result []LessonAsset
	for _, v := range val {
		result = append(result, LessonAsset{
			LessonID:   v.LessonID,
			LessonCode: v.LessonCode,
			AssetID:    v.AssetID.Int32,
			PageID:     v.PageID.Int32,
			Step:       v.Step.Int32,
			Type:       v.Type.String,
			Url:        URLJoin(baseurl, path, v.Url.String),
			Page:       v.Page.Int32})
	}
	return result
}

func LessonVideoResponse(val []GetLessonVideoRow, baseurl string, path string) []LessonVideo {
	var result []LessonVideo

	for _, v := range val {
		result = append(result, LessonVideo{
			LessonID:   v.LessonID,
			LessonCode: v.LessonCode,
			VideoID:    v.VideoID.Int32,
			Type:       v.Type.String,
			Url:        URLJoin(baseurl, path, v.Url.String)})
	}
	return result
}

func URLJoin(baseurl string, path ...string) string {
	result, err := url.JoinPath(baseurl, path...)
	if err != nil {
		log.Println("Error join url:", err)
	}
	return result
}
