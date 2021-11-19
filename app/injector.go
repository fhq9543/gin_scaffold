package app

import (
	"baseFrame/app/middleware"
	"baseFrame/app/router"
	"baseFrame/pkg/config"
	"baseFrame/pkg/logger"
	"baseFrame/pkg/sms"
	"context"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"github.com/google/wire"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

type Injector struct {
	Router     *router.Router
	Middleware *middleware.Middleware
	Cfg        *config.Config
	DB         *gorm.DB
	SMS        *sms.SMS
	//Redis      *redis.RedisClient
	Server *http.Server
}

func (i *Injector) Init() (cleanFunc func(), err error) {
	//env := config.GetConfig("", "env")
	loggerCleanFunc, err := logger.InitLoggerWithoutFile()
	if err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return nil, err
	}
	cleanFunc = func() {
		loggerCleanFunc()
	}

	return cleanFunc, nil
}

func (i *Injector) RunTask(ctx context.Context) {

	// 跑定时任务
	if i.Cfg.GetConfig("", "env") != "local" &&
		i.Cfg.GetConfig("", "cronTask") == "1" {
		logger.Debug("Server before StartTimer")
	}

}

func (i *Injector) Start(ctx context.Context) {

	// 启动定时任务
	i.RunTask(ctx)

	logger.Debug("Server listen on port" + i.Server.Addr)

	if err := i.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("listen: %s\n", err)
	}

}

func (i *Injector) Stop(ctx context.Context) {

	if err := i.Server.Shutdown(ctx); err != nil {
		//todo send forced shutdown notification
		logger.Fatalln("Server forced to shutdown:", err)
	}

}
