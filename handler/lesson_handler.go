package handler

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sysbitBroker/config"
	"sysbitBroker/domain/param"
	"sysbitBroker/domain/resp"
	"sysbitBroker/pkg"

	"sysbitBroker/domain/entity"
	"sysbitBroker/repository"
	"sysbitBroker/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type lessonHandler struct {
	Cfg        *config.Config
	LessonRepo repository.ILessonRepository
	Rdb        *pkg.Redis
}

func NewLessonHandler(LessonRepo repository.ILessonRepository, cfg *config.Config, Rdb *pkg.Redis) *lessonHandler {
	return &lessonHandler{
		Cfg:        cfg,
		LessonRepo: LessonRepo,
		Rdb:        Rdb,
	}
}

func (h *lessonHandler) CreateLesson(c *gin.Context) {
	var req param.AddLessonParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
		return
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.CreateOrUpdateLesson(ctx, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	if result.Operation == "inserted" {
		c.JSON(http.StatusCreated, resp.GenerateBaseResponse(result, true, resp.Success))
		return
	} else if result.Operation == "updated" {
		c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))
		return
	}

}
func (h *lessonHandler) GetLessons(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	result, err := h.LessonRepo.GetLesson(ctx)

	var lessons []entity.Lesson
	for _, v := range result {
		v.Image = entity.URLJoin(h.Cfg.S3.Baseurl, h.Cfg.S3.Path, "header", v.Image)
		lessons = append(lessons, v)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(lessons, true, resp.Success))
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
func (h *lessonHandler) CreateLessonPage(c *gin.Context) {
	var req param.AddLessonPageParams

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
func (h *lessonHandler) CreateOrUpdateLessonContent(c *gin.Context) {
	var req param.AddLessonContentParams
	// Bind JSON payload to the struct
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.CreateOrUpdateLessonContent(ctx, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	if result.Operation == "updated" {
		c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))
		return
	} else if result.Operation == "inserted" {
		c.JSON(http.StatusCreated, resp.GenerateBaseResponse(result, true, resp.Success))
		return
	}
}

func (h *lessonHandler) DeleteLessonContent(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	err := h.LessonRepo.DeleteLessonContent(ctx, utils.StrToInt32(id))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(nil, true, resp.Success))
}

func (h *lessonHandler) GetLessonDetail(c *gin.Context) {

	lang := c.GetHeader("Accept-Language")
	lessonCode := c.Query("code")

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	//Generated Query Key
	key := fmt.Sprintf("L%s%s", lessonCode, lang)
	//Get Result from Redis
	value, err := pkg.Get[entity.LessonData](ctx, h.Rdb.Rc, key)
	if err == redis.Nil {
		//Query to database and Cache when nonexist
		//Get Lesson Dialog
		result, err := h.LessonRepo.GetLessonDetail(ctx, lessonCode, []string{"en", lang})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
			return
		}
		var highestPageID int32
		for _, lesson := range result {
			if lesson.Page.Int32 > highestPageID {
				highestPageID = lesson.Page.Int32
			}
		}

		//Get Lesson Image Asset
		result2, err := h.LessonRepo.GetLessonAsset(ctx, lessonCode)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
			return
		}

		//Get Lesson Video
		result3, err := h.LessonRepo.GetLessonVideo(ctx, lessonCode)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
			return
		}

		var sceneUrl string
		for _, v := range result3 {
			if v.Type.String == "s" {
				sceneUrl = entity.URLJoin(h.Cfg.S3.Baseurl, h.Cfg.S3.Path, v.Url.String)
			}
		}

		for _, lesson := range result {
			if lesson.Page.Int32 > highestPageID {
				highestPageID = lesson.Page.Int32
			}
		}

		var page []entity.Page
		for i := 1; i <= int(highestPageID); i++ {
			//Retructure Assets
			var step0assets entity.StepAsset
			var step1assets entity.StepAsset
			for _, v := range entity.LessonResponse(result2, h.Cfg.S3.Baseurl, h.Cfg.S3.Path) {
				if int(v.Page) == i && v.Step == 0 {
					if v.Type == "i" {
						step0assets.Image = v.Url
					} else {
						step0assets.VoiceUrl = v.Url
					}
				}
				if int(v.Page) == i && v.Step == 1 {
					if v.Type == "i" {
						step1assets.Image = v.Url
					} else {
						step1assets.VoiceUrl = v.Url
					}
				}
			}
			//Restructure Dialog
			var step0content []entity.LessonContent
			var step1content []entity.LessonContent
			for _, v := range result {
				if int(v.Page.Int32) == i && v.Step.Int32 == 0 {
					step0content = append(step0content, entity.LessonContent{
						ContentID: v.ContentID.Int32,
						PageID:    0,
						Lang:      v.Lang.String,
						Step:      v.Step.Int32,
						Content:   v.Content.String,
					})
				}
				if int(v.Page.Int32) == i && v.Step.Int32 == 1 {
					step1content = append(step1content, entity.LessonContent{
						ContentID: v.ContentID.Int32,
						PageID:    0,
						Lang:      v.Lang.String,
						Step:      v.Step.Int32,
						Content:   v.Content.String,
					})
				}
			}
			steps := []entity.Step{
				{Step: 0,
					Contents: restructure(step0content),
					Assets:   step0assets},
				{Step: 1,
					Contents: restructure(step1content),
					Assets:   step1assets}}
			page = append(page, entity.Page{Page: i, Steps: steps})

		}

		lessonData := entity.LessonData{SceneUrl: sceneUrl, Pages: page}

		err = pkg.Set(ctx, h.Rdb.Rc, key, lessonData, time.Minute*time.Duration(60))
		if err != nil {
			fmt.Printf("error set cache %s", err)
		}
		c.JSON(http.StatusOK, resp.GenerateBaseResponse(lessonData, true, resp.Success))
		return
	} else if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, fmt.Errorf("error cache")))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(value, true, resp.Success))
}

func (h *lessonHandler) CreateOrUpdateLessonAsset(c *gin.Context) {
	var req param.AddOrUpdateLessonAssetParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
		return
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.CreateOrUpdateLessonAsset(ctx, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))
}

func (h *lessonHandler) GetLessonAssets(c *gin.Context) {
	lessonCode := c.Query("code")

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.GetLessonAsset(ctx, lessonCode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}

	result2, err := h.LessonRepo.GetLessonVideo(ctx, lessonCode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}

	var url []entity.LessonAssetUrl

	for _, v := range result {
		url = append(url, entity.LessonAssetUrl{Url: entity.URLJoin(h.Cfg.S3.Baseurl, h.Cfg.S3.Path, v.Url.String)})
	}

	for _, v := range result2 {
		if v.Type.String == "s" {
			url = append(url, entity.LessonAssetUrl{Url: entity.URLJoin(h.Cfg.S3.Baseurl, h.Cfg.S3.Path, v.Url.String)})
		}
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(url, true, resp.Success))
}

func (h *lessonHandler) DeleteLessonAssets(c *gin.Context) {
	id := c.Query("id")

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	err := h.LessonRepo.DeleteLessonAsset(ctx, utils.StrToInt32(id))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(nil, true, resp.Success))
}

func (h *lessonHandler) CreateOrUpdateLessonVideo(c *gin.Context) {
	var req param.AddOrUpdateLessonAudioParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
		return
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.CreateOrUpdateLessonVideo(ctx, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	if result.Operation == "inserted" {
		c.JSON(http.StatusCreated, resp.GenerateBaseResponse(result, true, resp.Success))
		return
	} else if result.Operation == "updated" {
		c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))
	}
}

func (h *lessonHandler) DeleteLessonVideo(c *gin.Context) {
	id := c.Query("id")

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	err := h.LessonRepo.DeleteLessonVideo(ctx, utils.StrToInt32(id))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(nil, true, resp.Success))
}

func (h *lessonHandler) GetLessonVideo(c *gin.Context) {
	lessonCode := c.Query("code")

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.GetLessonVideo(ctx, lessonCode)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(entity.LessonVideoResponse(result, h.Cfg.S3.Baseurl, h.Cfg.S3.Path), true, resp.Success))
}

func (h *lessonHandler) CreateOrUpdateQuiz(c *gin.Context) {
	var req param.AddOrUpdateQuizeStruct
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
		return
	}
	//Marshalling JSON to byte
	quiz := entity.StructToBytes[entity.QuizQuestion](req.Quiz)
	answer := entity.StructToBytes[entity.QuizAnswer](req.Answer)
	param := param.AddOrUpdateQuizeParams{
		LessonID: req.LessonID,
		Num:      req.Num,
		Lang:     req.Lang,
		Quiz:     quiz,
		Answer:   answer,
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.CreateOrUpdateQuiz(ctx, param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	if result.Operation == "inserted" {
		c.JSON(http.StatusCreated, resp.GenerateBaseResponse(result, true, resp.Success))
		return
	} else if result.Operation == "updated" {
		c.JSON(http.StatusOK, resp.GenerateBaseResponse(result, true, resp.Success))
	}
}

func (h *lessonHandler) GetQuizByCode(c *gin.Context) {
	lessonCode := c.Query("code")

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, err := h.LessonRepo.GetQuizByCode(ctx, lessonCode)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	//Sort by Number of Quiz
	sort.Slice(result, func(i, j int) bool {
		return result[i].Num.Int32 < result[j].Num.Int32
	})

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(entity.QuizResp(result), true, resp.Success))
}

func (h *lessonHandler) DeleteQuiz(c *gin.Context) {
	id := c.Query("id")

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	err := h.LessonRepo.DeleteQuiz(ctx, utils.StrToInt32(id))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithError(nil, true, resp.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, resp.GenerateBaseResponse(nil, true, resp.Success))
}

func restructure(val []entity.LessonContent) entity.LessonDial {
	var lesson entity.LessonDial
	for _, v := range val {
		if v.Lang == "en" {
			lesson.PrimaryID = v.ContentID
			lesson.PrimaryLang = v.Content
		} else {
			lesson.SecondaryID = v.ContentID
			lesson.SecondaryLang = v.Content
		}

	}
	return lesson
}
