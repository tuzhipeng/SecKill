package cache

import (
	"GraduateDesign/model/respStruct"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()

}

func TestRedisWorkFlow(t *testing.T)  {
	t.Run("测试写入商品详情到redis", testGetHotDataInsertToRedis)
	t.Run("测试从Redis读取商品详情", testGetGoodsDetailFromRedis)
	t.Run("测试从Redis读取商品首页列表", testGetGoodsListFromRedis)

}
func testGetHotDataInsertToRedis(t *testing.T) {
	GetHotDataInsertToRedis()
}


func testGetGoodsDetailFromRedis(t *testing.T) {
	var goodsDetailJson respStruct.GoodsDetailRespJson
	var  isEmpty bool
 	goodsDetailJson, isEmpty = GetGoodsDetailFromRedis("1m70y5k")
	if isEmpty== false {
		fmt.Println(goodsDetailJson)
	}else {
		fmt.Println("isEmpty 为 true")
	}

	goodsDetailJson, isEmpty = GetGoodsDetailFromRedis("123546")
	if isEmpty== false {
		fmt.Println(goodsDetailJson)
	}else {
		fmt.Println("isEmpty 为 true")
	}
}

func testGetGoodsListFromRedis(t *testing.T) {
	var goodsIndexList [] respStruct.List
	var isEmpty bool
	goodsIndexList, isEmpty = GetGoodsListFromRedis("new")
	fmt.Println(goodsIndexList, isEmpty)
}
