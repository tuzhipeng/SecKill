package main

import (
	"GraduateDesign/api/cache"
	"GraduateDesign/api/dao"
	"GraduateDesign/api/service"
	"GraduateDesign/data"
	"GraduateDesign/middleware"
	"GraduateDesign/model/reqStruct"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const port = 7777
var secKillGoodsStockInfo []dao.SecKillGoodsItem
var done chan int
var visitCount = -1
func initSecKillGoodsStockInfo()  {
	cache.GetHotDataInsertToRedis()
	secKillGoodsStockInfo = dao.SelectGoodsInfoFromSecKills()
	//for _, goodsItem := range secKillGoodsStockInfo{
	//	goodsItem.Stock = goodsItem.Stock/3 + 1
	//}
	done = make(chan int, 1)
	done <- 1

}

func localDecrementGoodsStock(iid string) bool  {
	fmt.Println("secKillGoodsStockInfo : ",secKillGoodsStockInfo)
	for index, tempGoods := range secKillGoodsStockInfo{
		if tempGoods.Iid == iid{
			visitCount ++
			// 本地预扣库存，如果扣取成功再远程访问redis扣库存
			if secKillGoodsStockInfo[index].Stock >0 && visitCount %10 == 0  {
				fmt.Println("visitCount is ", visitCount)
				secKillGoodsStockInfo[index].Stock --
				return true
			}else {
				return false
			}
		}
	}
	log.Fatal("secKillGoodsStockInfo中没有该商品ID:", iid)
	return false 
	//return true
}

func main() {

	initSecKillGoodsStockInfo()
	fmt.Println("secKillGoodsStockInfo : ",secKillGoodsStockInfo)
	router := gin.New()
	router.Use(middleware.Cors())
	router.Use(middleware.Auth())

	router.POST("api/secKillOrder", secKillOrderHandleFunc)
	defer data.Close()

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Panic("router运行时出错： " + err.Error())
	}

}

func secKillOrderHandleFunc(ctx *gin.Context)  {
	<-done
	//全局读写锁

	var postSecKillOrder reqStruct.OrderListItem
	ctx.BindJSON(&postSecKillOrder)
	goodsIid := postSecKillOrder.Iid
	if localDecrementGoodsStock(goodsIid) {
		res , err :=cache.RemoteDecrementGoodsStock(goodsIid)
		// 本地预扣库存成功，远程扣取失败，商品售罄
		if err == nil && res == 0{
			done <- 1
			ctx.JSON(http.StatusOK, gin.H{
				"success":      false,
				"message": "商品已售罄",
			})
		// 本地预扣库存成功，远程扣取也成功，生成订单，返回成功
		}else if err== nil && res == 1{
		
			uid, ok := ctx.Get("uid")
			if !ok {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			service.GenerateOrderService(uid.(string), []reqStruct.OrderListItem{postSecKillOrder})
			done <- 1
			ctx.JSON(http.StatusOK, gin.H{
				"success":      true,
				"message": "抢购成功",
			})
		}
		// 本地库存扣取失败，直接返回无货
	}else {
		done <- 1
		ctx.JSON(http.StatusOK, gin.H{
			"success":      false,
			"message": "商品已售罄",
		})
	}

}

