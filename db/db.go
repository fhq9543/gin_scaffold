package db

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"go_base/config"
)

var Xorm *xorm.Engine

func InitDB() {
	var err error
	Xorm, err = xorm.NewEngine("mysql", config.Viper.GetString("MYSQL_DSN"))

	if err != nil {
		panic(fmt.Sprintf("mysql init error. %s", err.Error()))
	}

	if gin.Mode() == gin.DebugMode {
		Xorm.ShowSQL(true)
		Xorm.Logger().SetLevel(core.LOG_WARNING)
	}
	//DB.SetMapper(core.SnakeMapper{})
	Xorm.DB().SetMaxIdleConns(0)
	Xorm.DB().SetMaxOpenConns(500)
}

func CloseDB() {
	Xorm.Close()
}
