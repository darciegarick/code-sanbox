package app

import (
	"jassue-gin/entity/request"
	"jassue-gin/entity/response"
	"jassue-gin/repository"
	"jassue-gin/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Register 用户注册
func Register(c *gin.Context) {
	var form request.Register
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	service := service.NewUserService(repository.NewUserRepository())
	if err, user := service.Register(form); err != nil {
		response.BusinessFail(c, err.Error())
		return
	} else {
		response.Success(c, user)
	}
}

// Login 用户登录
func Login(c *gin.Context) {
	var form request.Login
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	userservice := service.NewUserService(repository.NewUserRepository())
	if err, user := userservice.Login(form); err != nil {
		response.BusinessFail(c, err.Error())
		return
	} else {
		tokenData, err, _ := service.JwtService.CreateToken(service.AppGuardName, user)
		if err != nil {
			response.BusinessFail(c, err.Error())
			return
		}
		response.Success(c, tokenData)
	}
}

// GetUserInfo 获取用户信息
func Info(c *gin.Context) {
	var form request.Info
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	userservice := service.NewUserService(repository.NewUserRepository())
	err, user := userservice.GetUserInfo(form.Id)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, user)
}

// 用户登出
func Logout(c *gin.Context) {
	err := service.JwtService.JoinBlackList(c.Keys["token"].(*jwt.Token))
	if err != nil {
		response.BusinessFail(c, "登出失败")
		return
	}
	response.Success(c, nil)
}
