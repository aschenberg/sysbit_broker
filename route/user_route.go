package route

import (
	"sysbitBroker/config"
	"sysbitBroker/handler"
	"sysbitBroker/repository"
	"sysbitBroker/usecase"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func User(group *gin.RouterGroup, cfg *config.Config, pg *config.Postgres, provider *oidc.Provider, oidcCfg *oauth2.Config) {
	ar := repository.NewUserRepository(pg)
	au := usecase.NewUserUsecase(ar, cfg, 10*time.Second)
	ah := handler.NewUserHandler(au, cfg, oidcCfg, provider)
	group.POST("/login", ah.LoginOrRegister)
	group.POST("/loginweb", ah.LoginOrRegisterWeb)
	group.PUT("/refreshtoken", ah.GenerateAccessToken)

}
