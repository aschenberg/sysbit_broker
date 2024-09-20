package handler

import (
	"errors"
	"net/http"
	"sysbitBroker/config"
	"sysbitBroker/domain/entity"
	"sysbitBroker/domain/resp"
	"sysbitBroker/usecase"
	"sysbitBroker/utils"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type userHandler struct {
	User     usecase.IUserUsecase
	Cfg      *config.Config
	OidcCfg  *oauth2.Config
	Provider *oidc.Provider
}

func NewUserHandler(user usecase.IUserUsecase, cfg *config.Config, OidcCfg *oauth2.Config, Provider *oidc.Provider) *userHandler {
	return &userHandler{
		User:     user,
		Cfg:      cfg,
		OidcCfg:  OidcCfg,
		Provider: Provider,
	}
}

func (h *userHandler) LoginOrRegister(c *gin.Context) {
	// Define a struct to receive JSON data
	var request entity.LoginRequest

	// Bind JSON payload to the struct
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
		return
	}

	verifier := h.Provider.Verifier(&oidc.Config{ClientID: h.OidcCfg.ClientID})
	idToken, err := verifier.Verify(c, request.IdToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithError(nil, false, resp.AuthError, err))
		return
	}
	var claims entity.GoogleClaims
	if err := idToken.Claims(&claims); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithError(nil, false, resp.AuthError, err))
		return
	}

	//User Usecase
	user, status, err := h.User.LoginOrRegister(c, claims)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithAnyError(nil, true, resp.InternalError, err.Error()))
		return
	}
	if status == "inserted" {
		c.JSON(http.StatusCreated, resp.GenerateBaseResponse(user, true, resp.Success))

	} else if status == "updated" {
		c.JSON(http.StatusOK, resp.GenerateBaseResponse(user, true, resp.Success))

	}

}

func (h *userHandler) LoginOrRegisterWeb(c *gin.Context) {
	// Define a struct to receive JSON data
	var request entity.LoginRequest

	// Bind JSON payload to the struct
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
		return
	}

	if request.Email == "" || request.Os == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, errors.New("request field must not empty")))
		return
	}

	// verifier := lc.Provider.Verifier(&oidc.Config{ClientID: lc.OidcCfg.ClientID})

	// idToken, err := verifier.Verify(c, request.IDToken)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithError(nil, false, helper.AuthError, err))
	// 	return
	// }
	// var claims models.GoogleClaims
	// if err := idToken.Claims(&claims); err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithError(nil, false, helper.AuthError, err))
	// 	return
	// }

	claims := entity.GoogleClaims{
		Email: request.Email,
	}
	//User Usecase
	user, status, err := h.User.LoginOrRegister(c, claims)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.GenerateBaseResponseWithAnyError(nil, true, resp.InternalError, err.Error()))
		return
	}
	if status == "inserted" {
		c.JSON(http.StatusCreated, resp.GenerateBaseResponse(user, true, resp.Success))

	} else if status == "updated" {
		c.JSON(http.StatusOK, resp.GenerateBaseResponse(user, true, resp.Success))

	}

}

func (h *userHandler) GenerateAccessToken(c *gin.Context) {

	// Define a map to receive JSON data
	var request entity.RequestAccessToken
	// Bind JSON payload to the map
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
		return
	}

	// Extract the "refreshtoken" value
	if request.RefreshToken == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid refreshtoken"})
		return
	}

	//Extract ID from refresh token
	user, err := utils.ExtractIDFromToken(request.RefreshToken, h.Cfg.JWT.RefreshTokenSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithValidationError(nil, false, resp.ValidationError, err))
	}
	//Checked refresh token valid with database
	responseData, err := h.User.CheckRefreshToken(c, request.RefreshToken, user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.GenerateBaseResponseWithAnyError(nil, false, resp.ValidationError, err))
		return
	}
	if responseData.AccessToken == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, resp.GenerateBaseResponseWithAnyError(nil, false, resp.ValidationError, err))
		return
	}

	c.JSON(http.StatusOK, resp.GenerateBaseResponse(responseData, true, resp.Success))

}
