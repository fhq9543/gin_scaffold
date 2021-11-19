package router

func (r Router) UserRouter() {
	r.GET("users", r.AuthMiddleware(), r.User.UserList)
	user := r.Group("/user")
	{
		user.POST("login", r.User.UserLogin)
		user.PUT(":id/changePwd", r.AuthMiddleware(), r.User.UserChangePwd)
		user.GET(":id", r.AuthMiddleware(), r.User.UserDetail)
		user.PUT(":id", r.AuthMiddleware(), r.User.UserUpdate)
		user.DELETE(":id", r.AuthMiddleware(), r.User.UserDelete)
	}
}
