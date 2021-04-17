package controller

import (
	"GraduateDesign/api/service"
	"GraduateDesign/model/reqStruct"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(ctx *gin.Context) {
	var postUser = reqStruct.UserJson{}
	ctx.BindJSON(&postUser)

	token, err := service.LoginService(postUser, ctx)
	if err != nil {
		fmt.Println("loginService err : ", err)
	}

	fmt.Println("生成的token：", token)
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func GetUserInfo(ctx *gin.Context) {
	uid, ok := ctx.Get("uid")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	username, ok := ctx.Get("username")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"uid":      uid,
		"username": username,
	})

}
