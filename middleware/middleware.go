package middleware

import (
	"net/http"
	"strings"
	"sysbitBroker/domain/resp"
	"sysbitBroker/utils"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := utils.IsAuthorized(authToken, secret)
			if authorized {
				user, err := utils.ExtractIDFromToken(authToken, secret)
				if err != nil {
					if err.Error() == "expire" {
						c.AbortWithStatusJSON(http.StatusForbidden, resp.GenerateBaseResponseWithAnyError(nil, false, resp.InternalError, err.Error()))
						return
					}
					c.AbortWithStatusJSON(http.StatusUnauthorized, resp.GenerateBaseResponseWithAnyError(nil, false, resp.InternalError, err.Error()))
					return
				}
				c.Set("roles", user.Role)

				c.Set("userid", user.ID)
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp.GenerateBaseResponseWithAnyError(nil, false, resp.InternalError, err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.GenerateBaseResponse(nil, false, resp.AuthError))
	}
}
