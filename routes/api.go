package routes

import (
	"jassue-gin/controller/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetApiGroupRoutes 定义 api 分组路由
func SetApiGroupRoutes(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.POST("/auth/register", app.Register)
	router.POST("/auth/login", app.Login)
	router.POST("/auth/info", app.Info)

	router.POST("/sanbox/java", app.JavaSanBox)

	// authRouter := router.Group("").Use(middleware.JWTAuth(service.AppGuardName))
	// {
	// 	authRouter.POST("/auth/info", app.Info)
	// 	authRouter.POST("/auth/logout", app.Logout)
	// }

	// 示例demo 暂不做权限配置

}
