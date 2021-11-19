package auth

import (
	"baseFrame/pkg/config"
	"baseFrame/pkg/logger"
	"errors"
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

const algorithm = "HS256"

type JWToken struct {
	timeout int64
	secret  string
}

func initJWT(cfg *config.Config) *JWToken {
	timeOutStr := cfg.GetConfig("jwt", "timeout")
	timeout, _ := strconv.Atoi(timeOutStr)
	token := new(JWToken)
	token.timeout = int64(timeout)
	token.secret = cfg.GetConfig("jwt", "secret")
	return token
}

func (token *JWToken) JwtCreate(mc jwt.MapClaims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(algorithm), mc)
	// Sign and get the complete encoded token as a string
	tokenString, err := jwtToken.SignedString([]byte(token.secret))
	return tokenString, err
}

func (token *JWToken) JwtVerify(reqToken string) (claims jwt.MapClaims, err error) {
	jwtToken, err := jwt.Parse(reqToken, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != algorithm {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(token.secret), nil
	})
	if !logger.Check(err) {
		return nil, err
	}

	var ok bool
	claims, ok = jwtToken.Claims.(jwt.MapClaims)
	if !(ok && jwtToken.Valid) {
		return nil, errors.New("validation failure")
	}

	//uid, ok := claims["uid"].(string)
	//if !ok {
	//	return "", errors.New("not available uid")
	//}
	return
}

func (token *JWToken) Param(reqToken, key string) (interface{}, error) {
	jwtToken, err := jwt.Parse(reqToken, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != algorithm {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(token.secret), nil
	})
	if !logger.Check(err) {
		return "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !(ok && jwtToken.Valid) {
		return "", errors.New("validation failure")
	}

	uid, ok := claims[key]
	if !ok {
		return "", errors.New("not available uid")
	}
	return uid, nil
}
