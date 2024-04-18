package middleware

import (
	"jassue-gin/entity/response"
	"jassue-gin/global"
	"jassue-gin/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(GuardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.TokenFail(c)
			c.Abort()
			return
		}
		tokenStr = tokenStr[len(service.TokenType)+1:]

		// Token 解析校验
		token, err := jwt.ParseWithClaims(tokenStr, &service.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(global.App.Config.Jwt.Secret), nil
		})
		if err != nil || service.JwtService.IsInBlacklist(tokenStr) {
			response.TokenFail(c)
			c.Abort()
			return
		}
		claims := token.Claims.(*service.CustomClaims)
		// Token 发布者校验
		if claims.Issuer != GuardName {
			response.TokenFail(c)
			c.Abort()
			return
		}
		// token 续签
		timeDiff := int64(claims.ExpiresAt.Time.Sub(time.Now()).Seconds())
		if timeDiff < global.App.Config.Jwt.RefreshGracePeriod {
			lock := global.Lock("refresh_token_lock", global.App.Config.Jwt.JwtBlacklistGracePeriod)
			if lock.Get() {
				err, user := service.JwtService.GetUserInfo(GuardName, claims.ID)
				if err != nil {
					global.App.Log.Error(err.Error())
					lock.Release()
				} else {
					tokenData, _, _ := service.JwtService.CreateToken(GuardName, user)
					c.Header("new-token", tokenData.AccessToken)
					c.Header("new-expires-in", strconv.Itoa(tokenData.ExpiresIn))
					_ = service.JwtService.JoinBlackList(token)
				}
			}
		}
		c.Set("token", token)
		c.Set("id", claims.ID)
	}
}
