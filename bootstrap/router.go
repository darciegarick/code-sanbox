package bootstrap

import (
	"jassue-gin/docs"
	"jassue-gin/global"
	"jassue-gin/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/api"
	// 注册 api 分组路由
	apiGroup := router.Group("/api")
	routes.SetApiGroupRoutes(apiGroup)

	// 设置Swagger路由
	url := ginSwagger.URL("http://43.142.28.232:" + global.App.Config.App.Port + "/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}

// RunServer 启动服务器
func RunServer() {
	r := setupRouter()
	r.Run(":" + global.App.Config.App.Port)
}
