package controller

import (
	"GraduateDesign/api/service"
	"GraduateDesign/model/reqStruct"
	"GraduateDesign/model/respStruct"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GenerateOrder(ctx *gin.Context) {
	var postOrder = reqStruct.OrderJson{}
	ctx.BindJSON(&postOrder)
	uid, ok := ctx.Get("uid")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	service.GenerateOrderService(uid.(string), postOrder.List)

	ctx.JSON(http.StatusOK, gin.H{
		"success":      true,
		"getPostOrder": postOrder,
	})
}

func GetOrder(ctx *gin.Context)  {
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	uid, ok := ctx.Get("uid")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var orderRespJson respStruct.OrderRespJson
	orderRespJson = service.GetOrderService(uid.(string), page, 10)

	orderRespJsonBytes, err := json.Marshal(orderRespJson)
	if err != nil {
		log.Panic("json.Marshal(orderRespJson) err :", err)
	}

	ctx.Data(http.StatusOK, "application/json", orderRespJsonBytes)

}
