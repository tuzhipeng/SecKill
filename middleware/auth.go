package middleware

import (
	"GraduateDesign/api/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取header中的Authorization字段拿到token
		token := ctx.GetHeader("Authorization")
		fmt.Println("获取到的token :", token)

		if len(token) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, err := service.GetUserByToken(token)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err = service.RenewUserToken(token); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

		ctx.Set("uid", user.Uid)
		ctx.Set("username", user.Username)
		//ctx.Set("password", user.Password)
		ctx.Next()
	}
}
