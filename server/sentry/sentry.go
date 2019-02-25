package sentry

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"go_base/config"
)

func Init() {
	dsn := config.Viper.GetString("SENTRY_DSN")
	raven.SetDSN(dsn)
	raven.SetEnvironment(gin.Mode())
}
