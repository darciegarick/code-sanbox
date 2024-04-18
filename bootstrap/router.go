package bootstrap

import (
	"jassue-gin/global"
	"jassue-gin/routes"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	// 注册 api 分组路由
	apiGroup := router.Group("/api")
	routes.SetApiGroupRoutes(apiGroup)

	// 注册 openApi 分组路由
	openApiGroup := router.Group("/openapi")
	routes.SetOpenApiGroupRoutes(openApiGroup)

	return router
}

// RunServer 启动服务器
func RunServer() {
	r := setupRouter()
	r.Run(":" + global.App.Config.App.Port)
}
