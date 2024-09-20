package route

import (
	"sysbitBroker/config"
	"sysbitBroker/handler"
	"sysbitBroker/middleware"
	"sysbitBroker/repository"

	"github.com/gin-gonic/gin"
)

func LessonPage(group *gin.RouterGroup, cfg *config.Config, pg *config.Postgres) {
	ar := repository.NewLessonRepository(pg)
	lh := handler.NewLessonHandler(ar, cfg)
	group.POST("", middleware.Role([]string{"admin"}), lh.CreateLessonPage)
	group.GET("/:id", lh.GetLessonPage)
	group.DELETE("/:id", middleware.Role([]string{"admin"}), lh.DeleteLessonPage)

}
