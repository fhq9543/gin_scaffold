package middleware

import (
	"baseFrame/app/router"
	"baseFrame/pkg/cors"
	"baseFrame/pkg/logger"
	"baseFrame/pkg/response"
	corsLib "github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
)

type Middleware struct{}

func InitMiddlewareAndRouter(router *router.Router) *Middleware {
	// middlewares
	router.Use(corsLib.New(cors.Cors()))
	router.Use(logger.Logrus())
	router.Use(response.ResponseHeaders())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.RegisterRouter()

	return &Middleware{}
}
