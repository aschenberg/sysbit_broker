package usecase

import (
	"context"
	"errors"
	"sysbitBroker/config"
	"sysbitBroker/domain/entity"
	"sysbitBroker/pkg"
	"sysbitBroker/repository"
	"sysbitBroker/utils"

	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type IUserUsecase interface {
	LoginOrRegister(c context.Context, claims entity.GoogleClaims) (entity.LoginResponse, string, error)
	CheckRefreshToken(c context.Context, refresh string, userid string) (entity.RefreshTokenResp, error)
}

type userUsecase struct {
	User           repository.IUserRepository
	contextTimeout time.Duration
	Cfg            *config.Config
	// Log            logging.Logger
}

func NewUserUsecase(user repository.IUserRepository, cfg *config.Config, timeout time.Duration) IUserUsecase {
	return &userUsecase{
		User:           user,
		contextTimeout: timeout,
		Cfg:            cfg,
	}
}

func (uc *userUsecase) LoginOrRegister(c context.Context, claims entity.GoogleClaims) (entity.LoginResponse, string, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	createTime := time.Now().UTC().UnixMilli()
	//Generated Snow Flake ID
	keyID, err := pkg.GenerateID(1, 1, 1)
	if err != nil {
		// uc.Log.Error(logging.Snowflake, logging.CreatedID, err.Error(), nil)
		return entity.LoginResponse{}, "", err
	}

	//Create Or Update User
	created := repository.CreateOrUpdateUserParams{
		Userid: keyID, Appid: claims.Sub, Mail: claims.Email,
		Roles:    []string{"guest"},
		Isactive: true,
		Pic: pgtype.Text{
			String: claims.Picture,
			Valid:  utils.StrIsEmpty(claims.Picture)},
		Name: pgtype.Text{
			String: claims.Name,
			Valid:  utils.StrIsEmpty(claims.Name)},
		Createdat: createTime,
		Updatedat: createTime}

	usr, err := uc.User.CreateOrUpdate(ctx, created)
	if err != nil {
		return entity.LoginResponse{}, "", err
	}

	//Generated Refresh Token
	user := entity.User{
		UserID: utils.Int64ToStr(usr.UserID),
		Email:  usr.Email, Role: usr.Role}

	refreshToken, err := pkg.CreateRefreshToken(
		user, uc.Cfg.JWT.RefreshTokenSecret,
		uc.Cfg.JWT.RefreshTokenExpireHour)
	if err != nil {
		// uc.Log.Error(logging.JWT, logging.GenerateToken, err.Error(), nil)
		return entity.LoginResponse{}, "", err
	}

	//Updated User Refresh Token in Database
	tokenPr := repository.UpdateUserTokenParams{
		RefreshToken: refreshToken, UserID: usr.UserID,
	}
	err = uc.User.UpdateRefreshToken(ctx, tokenPr)
	if err != nil {
		return entity.LoginResponse{}, "", err
	}

	//Generated Access Token
	accessToken, err := pkg.CreateAccessToken(
		user, uc.Cfg.JWT.RefreshTokenSecret,
		uc.Cfg.JWT.AccessTokenExpireMinutes)
	if err != nil {
		// uc.Log.Error(logging.JWT, logging.GenerateToken, err.Error(), nil)
		return entity.LoginResponse{}, "", err
	}

	response := entity.LoginResponse{
		ID: utils.Int64ToStr(usr.UserID),
		// Name:         usr.Name.String,
		// Email:        usr.Email,
		// Picture:      usr.Picture.String,
		Role:         usr.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response, usr.Operation, nil
}

func (uc *userUsecase) CheckRefreshToken(c context.Context, refresh string, userid string) (entity.RefreshTokenResp, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	user, err := uc.User.GetUser(ctx, utils.StrToInt64(userid))
	if err != nil {
		return entity.RefreshTokenResp{}, err
	}
	//Return false when refresh token doesn't match
	if user.RefreshToken != refresh {
		return entity.RefreshTokenResp{}, errors.New("token doesn't match")
	}

	access, err := pkg.CreateAccessToken(entity.User{UserID: utils.Int64ToStr(user.UserID), Role: user.Role}, uc.Cfg.JWT.AccessTokenSecret, uc.Cfg.JWT.AccessTokenExpireMinutes)
	if err != nil {
		return entity.RefreshTokenResp{}, err
	}

	return entity.RefreshTokenResp{
		AccessToken: access,
		ID:          utils.Int64ToStr(user.UserID),
		Role:        user.Role}, nil
}
