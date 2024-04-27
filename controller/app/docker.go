package app

import (
	"jassue-gin/entity/request"
	"jassue-gin/entity/response"
	"jassue-gin/sanbox/docker"

	"github.com/gin-gonic/gin"
)

func JavaDocker(c *gin.Context) {
	// sample demo
	var form request.ExecuteCodeRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		return
	}
	result, err := docker.DockerExecuteCode(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
	}
	response.Success(c, result)
}
