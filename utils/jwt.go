package utils

import (
	"fmt"
	"sysbitBroker/domain/entity"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(user entity.User, secret string, expiry int) (accessToken string, err error) {
	exp := time.Now().Add(time.Minute * time.Duration(expiry)).UnixMilli()
	now := time.Now().UnixMilli()

	claims := &entity.JwtCustomClaims{
		Role: user.Role,
		ID:   user.UserID,

		MapClaims: jwt.MapClaims{
			"sub":  user.UserID,
			"name": user.GivenName + user.FamilyName,
			"iat":  now,
			"exp":  exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(user entity.User, secret string, expiry int) (refreshToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).UnixMilli()
	now := time.Now().UnixMilli()
	claimsRefresh := &entity.JwtCustomRefreshClaims{
		ID:   user.UserID,
		Role: user.Role,
		MapClaims: jwt.MapClaims{
			"sub":  user.UserID,
			"name": user.GivenName + user.FamilyName,
			"iat":  now,
			"exp":  exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractIDFromToken(requestToken string, secret string) (entity.JwtSetClaims, error) {
	token, err := jwt.ParseWithClaims(requestToken, &entity.JwtCustomRefreshClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return entity.JwtSetClaims{}, err
	}

	claims, ok := token.Claims.(*entity.JwtCustomRefreshClaims)

	if !ok && !token.Valid {
		return entity.JwtSetClaims{}, fmt.Errorf("invalid token")
	}
	expire, err := claims.MapClaims.GetExpirationTime()

	if err != nil {
		return entity.JwtSetClaims{}, fmt.Errorf("fail get expire time")
	}

	if expire.Unix() < time.Now().Local().Unix() {
		return entity.JwtSetClaims{}, fmt.Errorf("expire")
	}

	return entity.JwtSetClaims{ID: claims.ID, Role: claims.Role}, nil
}
