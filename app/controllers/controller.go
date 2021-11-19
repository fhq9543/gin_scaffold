package controllers

import (
	ctlDemo "baseFrame/app/controllers/demo"
	ctlUser "baseFrame/app/controllers/user"

	"github.com/google/wire"
)

var ControllerSet = wire.NewSet(
	wire.Struct(new(ctlDemo.DemoCtl), "*"),
	wire.Struct(new(ctlUser.UserCtl), "*"),
)
