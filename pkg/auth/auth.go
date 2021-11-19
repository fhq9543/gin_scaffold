package auth

import (
	"baseFrame/pkg/config"
	"baseFrame/pkg/logger"
	"baseFrame/pkg/response"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	jwtSession *JWToken
}

func InitAuth(cfg *config.Config) (a *Auth) {
	a = &Auth{
		jwtSession: initJWT(cfg),
	}
	return
}

func (a *Auth) AppLogin(userID uint) (jwttoken string, err error) {
	jwttoken, err = a.NewAppJwt(userID)
	if !logger.Check(err) {
		return
	}
	return
}

func (a *Auth) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := request.OAuth2Extractor.ExtractToken(c.Request)
		if err != nil {
			c.Abort()
			response.Error(c, 10002, "token无效", gin.H{})
			return
		}
		claims, e := a.jwtSession.JwtVerify(tokenString)
		if !logger.Check(e) {
			logger.Debug("Expired token:", tokenString)
			c.Abort()
			response.Error(c, 10003, "token无效", gin.H{})
			return
		}
		userID, ok := claims["uid"].(float64)
		if !ok {
			c.Abort()
			response.Error(c, 10004, "未登录", gin.H{})
			return
		}

		c.Set("userID", int(userID))
		c.Next()
	}
}

func (a *Auth) GetUserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := request.OAuth2Extractor.ExtractToken(c.Request)
		if err == nil {
			claims, e := a.jwtSession.JwtVerify(tokenString)
			if logger.Check(e) {
				userID, ok := claims["uid"].(float64)
				if ok {
					c.Set("userID", int(userID))
				}
			}
		}
		c.Next()
	}
}

func (a *Auth) NewAppJwt(uid uint) (string, error) {
	if uid == 0 {
		return "", errors.New("uid is empty")
	}
	exp := time.Now().Add(time.Second * time.Duration(a.jwtSession.timeout)).Unix()
	return a.jwtSession.JwtCreate(jwt.MapClaims{
		"uid": uid,
		"exp": exp,
	})
}
