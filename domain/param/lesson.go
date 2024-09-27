package param

import "sysbitBroker/domain/entity"

type AddLessonParams struct {
	LessonCode string `json:"lesson_code"`
	Title      string `json:"title"`
	Image      string
}

type AddOrUpdateQuizeStruct struct {
	LessonID int32               `json:"lesson_id"`
	Num      int32               `json:"num"`
	Lang     string              `json:"lang"`
	Quiz     entity.QuizQuestion `json:"quiz"`
	Answer   entity.QuizAnswer   `json:"answer"`
}
type AddOrUpdateQuizeParams struct {
	LessonID int32  `json:"lesson_id"`
	Num      int32  `json:"num"`
	Lang     string `json:"lang"`
	Quiz     []byte `json:"quiz"`
	Answer   []byte `json:"answer"`
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

type UpdateLessonContentParams struct {
	ContentID int32  `json:"content_id"`
	Content   string `json:"content"`
}
type UpdateLessonAssetParams struct {
	AssetID int32  `json:"asset_id"`
	Url     string `json:"url"`
}

type AddOrUpdateLessonAssetParams struct {
	PageID   int32  `json:"page_id"`
	LessonID int32  `json:"lesson_id"`
	Step     int32  `json:"step"`
	Type     string `json:"type"`
	Url      string `json:"url"`
}

type AddOrUpdateLessonAudioParams struct {
	LessonID int32  `json:"lesson_id"`
	Type     string `json:"type"`
	Url      string `json:"url"`
}
