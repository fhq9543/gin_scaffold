package router

func (r Router) DemoRouter() {
	demo := r.Group("/demo")
	{
		demo.POST("info", r.AuthMiddleware(), r.Demo.DemoCreate)
		demo.GET("info", r.Demo.DemoList)
		demo.GET("info/:id", r.Demo.DemoDetail)
		demo.PUT("info/:id", r.Demo.DemoUpdate)
		demo.DELETE("info/:id", r.Demo.DemoDelete)
	}
}
