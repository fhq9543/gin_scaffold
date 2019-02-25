package routers

import (
	"github.com/gin-gonic/gin"
	"go_base/controllers"
	"go_base/middlewares"
	"go_base/utils/rescode"
	"net/http"
)

func errWrapper(handler gin.HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"rescode": rescode.Error, "data": nil, "msg": r})
			}
		}()
		handler(c)
	}
}

func Register(router *gin.Engine) {

	router.GET("/aggregators_info", controllers.ListAggregator)
	router.GET("/aggregators_info/:id", controllers.DetailAggregator)

	chief := router.Group("/aggregators")
	chief.Use(middlewares.AdminAuthentication())
	{
		chief.GET("", controllers.ListAggregator)
		chief.POST("", controllers.CreateAggregator)
		chief.PUT("/:id", controllers.UpdateAggregator)
		chief.POST("/:id", controllers.UpdateAggregator)
		chief.DELETE("/:id", controllers.DeleteAggregator)
		chief.GET("/:id", controllers.DetailAggregator)
	}

}
