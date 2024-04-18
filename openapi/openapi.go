package openapi

import (
	"jassue-gin/entity/response"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name  string
	phone string
	// 其他用户信息...
}

var users = map[string]User{
	"1":  {"Alice", "123456789"},
	"2":  {"Bob", "987654321"},
	"3":  {"Ziying", "111111111"},
	"4":  {"David", "222222222"},
	"5":  {"Emily", "333333333"},
	"6":  {"Frank", "444444444"},
	"7":  {"Grace", "555555555"},
	"8":  {"Henry", "666666666"},
	"9":  {"Ivy", "777777777"},
	"10": {"Jack", "888888888"},
}

func GetUserByID(userID string) (User, error) {
	user, ok := users[userID]
	if !ok {
		return User{}, nil // 或者返回一个错误，表示用户不存在
	}
	return user, nil
}

func GetUserNameRESTful(c *gin.Context) {
	userID := c.Param("id")
	user, err := GetUserByID(userID)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	} else {
		response.Success(c, user.Name)
	}
}
