package service

import (
	"context"
	"errors"
	"fmt"
	"jassue-gin/global"
	"jassue-gin/repository"

	"jassue-gin/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
}

var JwtService = new(jwtService)

// 所有需要颁发 token 的用户模型必须实现这个接口
type JwtUser interface {
	GetUid() string
}

// CustomClaims 自定义 Claims
type CustomClaims struct {
	jwt.RegisteredClaims
}

const (
	TokenType    = "bearer"
	AppGuardName = "app"
)

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (jwtService *jwtService) CreateToken(GuardName string, user JwtUser) (tokenData TokenOutPut, err error, token *jwt.Token) {
	token = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.App.Config.Jwt.JwtTtl) * time.Second)),
				ID:        user.GetUid(),
				Issuer:    GuardName,
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		},
	)
	tokenStr, err := token.SignedString([]byte(global.App.Config.Jwt.Secret))
	tokenData = TokenOutPut{
		tokenStr,
		int(global.App.Config.Jwt.JwtTtl),
		TokenType,
	}
	return
}

// 获取黑名单缓存 key
func (jwtService *jwtService) getBlackListKey(tokenStr string) string {
	result := "jwt_black_list:" + util.MD5([]byte(tokenStr))
	return result
}

// JoinBlackList token 加入黑名单
func (jwtService *jwtService) JoinBlackList(token *jwt.Token) (err error) {
	expiresAt := token.Claims.(*CustomClaims).ExpiresAt.Time
	now := time.Now()
	timer := time.Duration(int64(expiresAt.Sub(now).Seconds())) * time.Second
	// 将 token 剩余时间设置为缓存有效期，并将当前时间作为缓存 value 值
	err = global.App.Redis.SetNX(context.Background(), jwtService.getBlackListKey(token.Raw), now.Unix(), timer).Err()
	return
}

// IsInBlacklist token 是否在黑名单中
func (jwtService *jwtService) IsInBlacklist(tokenStr string) bool {
	joinUnixStr, err := global.App.Redis.Get(context.Background(), jwtService.getBlackListKey(tokenStr)).Result()
	fmt.Println("joinUnixStr", joinUnixStr)

	joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	if joinUnixStr == "" || err != nil {
		return false
	}
	// JwtBlacklistGracePeriod 为黑名单宽限时间，避免并发请求失效
	if time.Now().Unix()-joinUnix < global.App.Config.Jwt.JwtBlacklistGracePeriod {
		return false
	}
	return true
}

// GetUserInfo() 根据不同客户端 token ，查询不同用户表数据
func (jwtService *jwtService) GetUserInfo(GuardName string, id string) (err error, user JwtUser) {
	switch GuardName {
	case AppGuardName:
		repo := repository.NewUserRepository()
		userService := NewUserService(repo)
		return userService.GetUserInfo(id)
	default:
		err = errors.New("guard " + GuardName + " does not exist")
	}
	return
}
