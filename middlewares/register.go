package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	//设置cors
	engine.Use(CorsMiddleware())
	//全局错误处理
	engine.Use(ExceptionMiddleware())

	//自定义中间件
	//engine.Use(SessionMiddleware(store))
}
