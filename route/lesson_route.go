package route

import (
	"sysbitBroker/config"
	"sysbitBroker/handler"
	"sysbitBroker/middleware"
	"sysbitBroker/pkg"
	"sysbitBroker/repository"

	"github.com/gin-gonic/gin"
)

func Lesson(group *gin.RouterGroup, cfg *config.Config, pg *pkg.Postgres, rdb *pkg.Redis) {
	ar := repository.NewLessonRepository(pg)
	lh := handler.NewLessonHandler(ar, cfg, rdb)
	//======================================================================
	//Main Lesson
	group.POST("", middleware.Role([]string{"admin"}), lh.CreateLesson)
	group.GET("", lh.GetLessons)
	group.GET("/detail", lh.GetLessonDetail)
	group.DELETE("/:id", middleware.Role([]string{"admin"}), lh.DeleteLesson)
	//======================================================================
	//Page
	group.POST("/page", middleware.Role([]string{"admin"}), lh.CreateLessonPage)
	group.GET("/page/:id", lh.GetLessonPage)
	group.DELETE("/page/:id", middleware.Role([]string{"admin"}), lh.DeleteLessonPage)
	//======================================================================
	//Content
	group.POST("/content", middleware.Role([]string{"admin"}), lh.CreateOrUpdateLessonContent)
	group.DELETE("/content/:id", middleware.Role([]string{"admin"}), lh.DeleteLessonContent)
	//======================================================================
	//Assets
	group.POST("/assets", middleware.Role([]string{"admin"}), lh.CreateOrUpdateLessonAsset)
	group.GET("/assets", lh.GetLessonAssets)
	group.DELETE("/assets/:id", middleware.Role([]string{"admin"}), lh.DeleteLessonAssets)
	//======================================================================
	//Video
	group.POST("/video", middleware.Role([]string{"admin"}), lh.CreateOrUpdateLessonVideo)
	group.GET("/video", lh.GetLessonVideo)
	group.DELETE("/video/:id", middleware.Role([]string{"admin"}), lh.DeleteLessonVideo)
	//======================================================================
	//Quiz
	group.POST("/quiz", middleware.Role([]string{"admin"}), lh.CreateOrUpdateQuiz)
	group.GET("/quiz", lh.GetQuizByCode)
	group.DELETE("/quiz/:id", middleware.Role([]string{"admin"}), lh.DeleteQuiz)

}
