package app

import (
	"jassue-gin/entity/request"
	"jassue-gin/entity/response"
	"jassue-gin/sanbox"

	"github.com/gin-gonic/gin"
)

// PingExample godoc
// @Summary  Java 的代码沙箱
// @Schemes
// @Description 用于运行用户的 Java 代码
// @Tags example
// @Accept json
// @Produce json
// @Param    code  body  string  true  "Code to execute"  example({"code": "Java 代码"})
// @Success 200 {string} CodeResult
// @Router /api/sanbox/java [post]
func JavaSanBox(c *gin.Context) {
	// sample demo
	var form request.ExecuteCodeRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		return
	}
	result, err := sanbox.ExecuteCode(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
	}
	response.Success(c, result)
}
