package server

import (
	"go_base/config"
	"go_base/server/sentry"
	"go_base/utils/log"
	"go_base/utils/validation"
)

func Init(configPath string, logPath string) {
	//配置文件初始化
	config.Init(configPath)
	//日志初始化
	log.Init(logPath)

	//sentry.Init()

	//设置beego validation 错误信息
	validation.SetValidationMessage()
}
