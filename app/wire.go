// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package app

import (
	"baseFrame/app/controllers"
	"baseFrame/app/db"
	"baseFrame/app/middleware"
	"baseFrame/app/router"
	"baseFrame/pkg/auth"
	"baseFrame/pkg/config"
	"baseFrame/pkg/response"
	"baseFrame/pkg/route"
	"baseFrame/pkg/sms"

	"github.com/google/wire"
)

func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		config.Init,
		//redis.InitRedis,
		sms.InitSMS,
		response.InitResponse,
		db.DBSet,
		auth.InitAuth,
		route.InitRouter,
		controllers.ControllerSet,
		router.SetOfRouter,
		middleware.InitMiddlewareAndRouter,
		route.NewHttpServer,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
