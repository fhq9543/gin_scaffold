package route

import (
	"baseFrame/pkg/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	gin.DisableConsoleColor()
	router := gin.Default()
	return router
}

func NewHttpServer(conf *config.Config, engine *gin.Engine) *http.Server {
	appPort := conf.GetConfig("", "appPort")
	if appPort == "" {
		appPort = ":9000"
	}

	srv := &http.Server{
		Addr:    appPort,
		Handler: engine,
	}

	return srv
}
