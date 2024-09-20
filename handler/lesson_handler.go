package handler

import (
	"context"
	"net/http"
	"sysbitBroker/config"
	"sysbitBroker/domain/entity"
	"sysbitBroker/domain/resp"
	"sysbitBroker/repository"
	"sysbitBroker/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type lessonHandler struct {
	Cfg        *config.Config
	LessonRepo repository.ILessonRepository
	// OidcCfg  *oauth2.Config
	// Provider *oidc.Provider
}

func NewLessonHandler(LessonRepo repository.ILessonRepository, cfg *config.Config) *lessonHandler {
	return &lessonHandler{
		Cfg:        cfg,
		LessonRepo: LessonRepo,
	}
}

func (h *lessonHandler) GetLessons(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	result, err := h.LessonRepo.GetLesson(ctx)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))
}

func (h *lessonHandler) DeleteLesson(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	id := c.Param("id")
	err := h.LessonRepo.DeleteLesson(ctx, utils.StrToInt32(id))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(nil, true, resp.Success))
}

func (h *lessonHandler) CreateLesson(c *gin.Context) {
	var req entity.AddLessonParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
		return
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.CreateLesson(ctx, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))
}

func (h *lessonHandler) CreateLessonPage(c *gin.Context) {
	var req entity.AddLessonPageParams

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
		return
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.CreateLessonPage(ctx, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))
}

func (h *lessonHandler) GetLessonPage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	id := c.Param("id")
	result, err := h.LessonRepo.GetLessonPage(ctx, utils.StrToInt32(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))
}

func (h *lessonHandler) DeleteLessonPage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	id := c.Param("id")
	err := h.LessonRepo.DeleteLessonPage(ctx, utils.StrToInt32(id))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(nil, true, resp.Success))
}

func (h *lessonHandler) CreateLessonContent(c *gin.Context) {
	var req entity.AddLessonContentParams
	// Bind JSON payload to the struct
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.CreateLessonContent(ctx, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))
}

func (h *lessonHandler) GetLessonDetail(c *gin.Context) {

	lang := c.GetHeader("Accept-Language")
	lessonCode := c.Query("code")

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.GetLessonDetail(ctx, lessonCode, []string{"en", lang})
	var highestPageID int32
	for _, lesson := range result {
		if lesson.Page.Int32 > highestPageID {
			highestPageID = lesson.Page.Int32
		}
	}

	var page []entity.Page
	for i := 1; i <= int(highestPageID); i++ {
				var step0content []entity.LessonContent
				var step1content []entity.LessonContent
		for _, v := range result {
			if int(v.Page.Int32) == i && v.Step.Int32 == 0{
                    step0content =  append(step0content, entity.LessonContent{
						ContentID: v.ContentID.Int32,
						PageID:    0,
						Lang:      v.Lang.String,
						Step:      v.Step.Int32,
						Content:   v.Content.String,
					})
			}
            if int(v.Page.Int32) == i && v.Step.Int32 == 1{
				step1content =  append(step1content, entity.LessonContent{
					ContentID: v.ContentID.Int32,
					PageID:    0,
					Lang:      v.Lang.String,
					Step:      v.Step.Int32,
					Content:   v.Content.String,
				})
			}
		}
        steps:= []entity.Step{{
			Step: 0,
			Content: restructure(step0content),
		},{Step: 1,Content: restructure(step1content)}}
		page = append(page, entity.Page{Page: i, Steps:steps })

	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(page, true, resp.Success))

}

func restructure(val []entity.LessonContent)entity.LessonDial{
	var lesson entity.LessonDial
	for _, v := range val {
		if	v.Lang == "en"{
			lesson.PrimaryID = v.ContentID
			lesson.PrimaryLang = v.Content
		}else{
			lesson.SecondaryID = v.ContentID
			lesson.SecondaryLang = v.Content
		}

	}
	return lesson
}

// 	type OKReply struct {
// 		Status     string
// 		LessonData entity.LessonData
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	var okReply OKReply
// 	okReply.Status = "OK"
// 	okReply.LessonData = lessonData
// 	json.NewEncoder(w).Encode(okReply)

// }

// func (h *lessonHandler) SyncApp(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {

// 	fileName := ""
// 	fileNames := []string{}

// 	appId, ok := r.Context().Value("appId").(string)

// 	if !ok {
// 		errHandler.ErrMsg(w, errors.New("application id not in claims"), http.StatusInternalServerError)
// 		return
// 	}

// 	rows, err := conn.Query(context.Background(), "select fileName from sync where updatedOn >= (Select lastSync from app where appId  = $1) and substring(fileName,2,4) <= ( select lessonCode from progress where appId = $2 order by lessonCode desc limit 1)", appId, appId)

// 	if err != nil {
// 		errHandler.ErrMsg(w, err, http.StatusInternalServerError)
// 		return
// 	}

// 	defer rows.Close()

// 	for rows.Next() {

// 		if err := rows.Scan(&fileName); err != nil {
// 			errHandler.ErrMsg(w, err, http.StatusInternalServerError)
// 			return
// 		}

// 		fileNames = append(fileNames, fileName)

// 	}

// 	updAppStmt := "Update app set lastSync = Now() where appId = $1"

// 	_, err2 := conn.Exec(context.Background(), updAppStmt, appId)

// 	if err2 != nil {
// 		errHandler.ErrMsg(w, err2, http.StatusInternalServerError)
// 		return
// 	}

// 	type OKReply struct {
// 		Status    string
// 		FileNames []string
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	var okReply OKReply
// 	okReply.Status = "OK"
// 	okReply.FileNames = fileNames
// 	json.NewEncoder(w).Encode(okReply)

// }

// func (h *lessonHandler) GetQuizDetail(w http.ResponseWriter, r *http.Request, conn *pgx.Conn, quizCode string, lingoCode string) {

// 	var strQuizData string

// 	if err := conn.QueryRow(context.Background(), "select quizData from quiz where quizCode = $1 and lingoCode = $2", quizCode, lingoCode).Scan(&strQuizData); err != nil {
// 		errHandler.ErrMsg(w, err, http.StatusInternalServerError)
// 		return
// 	}

// 	bytQuizData := []byte(strQuizData)
// 	var jsonQuizData QuizData

// 	err := json.Unmarshal(bytQuizData, &jsonQuizData)
// 	if err != nil {
// 		errHandler.ErrMsg(w, err, http.StatusInternalServerError)
// 		return
// 	}

// 	type OKReply struct {
// 		Status   string
// 		QuizData QuizData
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	var okReply OKReply
// 	okReply.Status = "OK"
// 	okReply.QuizData = jsonQuizData
// 	json.NewEncoder(w).Encode(okReply)

// }

// func (h *lessonHandler) GetQuiz(w http.ResponseWriter, r *http.Request, conn *pgx.Conn, quizCode string) {

// 	var strQuizData string

// 	if err := conn.QueryRow(context.Background(), "select quizData from quiz where quizCode = $1", quizCode).Scan(&strQuizData); err != nil {
// 		errHandler.ErrMsg(w, err, http.StatusInternalServerError)
// 		return
// 	}

// 	bytQuizData := []byte(strQuizData)
// 	var jsonQuizData QuizData

// 	err := json.Unmarshal(bytQuizData, &jsonQuizData)
// 	if err != nil {
// 		errHandler.ErrMsg(w, err, http.StatusInternalServerError)
// 		return
// 	}

// 	type OKReply struct {
// 		Status   string
// 		QuizData QuizData
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	var okReply OKReply
// 	okReply.Status = "OK"
// 	okReply.QuizData = jsonQuizData
// 	json.NewEncoder(w).Encode(okReply)

// }

// func (h *lessonHandler) UpdQuiz(w http.ResponseWriter, r *http.Request, conn *pgx.Conn, quizCode string, lingoCode string) {

// 	quizData := QuizData{
// 		Quizes: []Quiz{},
// 	}

// 	json.NewDecoder(r.Body).Decode(&quizData)

// 	updAppStmt := "Update quiz set quizData = $1 where quizCode = $2 and lingoCode = $3"

// 	_, err := conn.Exec(context.Background(), updAppStmt, quizData, quizCode, lingoCode)

// 	if err != nil {
// 		errHandler.ErrMsg(w, err, http.StatusInternalServerError)
// 		return
// 	}

// 	var okReply OKReply
// 	okReply.Status = "OK"
// 	okReply.Message = "Quiz Updated"
// 	okReply.Message = fmt.Sprintf("Quiz %s Updated", quizCode)
// 	json.NewEncoder(w).Encode(okReply)

// }
