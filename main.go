package main

import (
	"github.com/gin-gonic/gin"
	"go_base/config"
	"go_base/db"
	"go_base/middlewares"
	"go_base/models"
	"go_base/routers"
	"go_base/server"
)

func main() {
	engine := gin.Default()

	//配置初始化
	server.Init("config", "logs")

	DeployMode := config.Viper.GetString("ROBO_DEPLOY_MODE")
	// Set gin context
	if DeployMode == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	//数据库初始化
	db.InitDB()
	defer db.CloseDB()

	// 建表
	models.Migrate()

	//注册中间件
	middlewares.Register(engine)
	//注册路由
	routers.Register(engine)

	engine.Run(":" + config.Viper.GetString("PORT"))
}
