package router

import (
	ctlDemo "baseFrame/app/controllers/demo"
	ctlUser "baseFrame/app/controllers/user"
	"baseFrame/pkg/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type Router struct {
	*gin.Engine
	*auth.Auth
	Demo ctlDemo.DemoCtl
	User ctlUser.UserCtl
}

var SetOfRouter = wire.NewSet(wire.Struct(new(Router), "*"))

func (r Router) RegisterRouter() {
	r.GET("ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, gin.Mode()+",pong,"+time.Now().String())
	})

	// router
	r.DemoRouter()
	r.UserRouter()
}
