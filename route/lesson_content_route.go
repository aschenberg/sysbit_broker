package route

import (
	"sysbitBroker/config"
	"sysbitBroker/handler"
	"sysbitBroker/middleware"
	"sysbitBroker/repository"

	"github.com/gin-gonic/gin"
)

func LessonContent(group *gin.RouterGroup, cfg *config.Config, pg *config.Postgres) {
	ar := repository.NewLessonRepository(pg)
	lh := handler.NewLessonHandler(ar, cfg)
	group.POST("", middleware.Role([]string{"admin"}), lh.CreateLessonContent)
	group.GET("/detail/:id", lh.GetLessonDetail)
	group.DELETE("/:id", middleware.Role([]string{"admin"}), lh.DeleteLessonPage)

}
