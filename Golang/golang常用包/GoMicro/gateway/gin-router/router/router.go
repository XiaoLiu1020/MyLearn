package router

import (
	"api-gateway/gin-router/handlers"
	"api-gateway/gin-router/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter(service ...interface{}) *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors(), middleware.InitMiddleware(service), middleware.ErrorMiddleware())
	store := cookie.NewStore([]byte("something-very-secret"))
	ginRouter.Use(sessions.Sessions("mysession", store))
	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})
		userV1 := v1.Group("/user")
		// 用户服务
		userV1.GET("register", handlers.UserRegister)
		userV1.GET("login", handlers.UserLogin)

		/* // 需要登录保护
		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.GET("tasks", handlers.GetTaskList)
			authed.POST("task", handlers.CreateTask)
			authed.GET("task/:id", handlers.GetTaskDetail) // task_id
			authed.PUT("task/:id", handlers.UpdateTask)    // task_id
			authed.DELETE("task/:id", handlers.DeleteTask) // task_id
		} */
	}
	return ginRouter
}
