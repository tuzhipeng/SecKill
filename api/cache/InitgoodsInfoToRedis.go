package cache

import (
	"GraduateDesign/api/dao"
	"GraduateDesign/api/service"
	"GraduateDesign/data"
	"GraduateDesign/model/respStruct"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

var GoodsDetailKey       = "goodsDetail:iid:%s"
var GoodsDetailExpire    = 3600 * 24 * 30
var SellGoodsKey = "sellGoodsList" // 热卖商品列表Key
var SecKillGoodsKey = "secKillGoodsList" // 秒杀商品列表key
var NewGoodsKey = "newGoodsList" // 新品商品列表key
var SecKillGoodsStockKey = "secKillGoodsStock:iid:%s"



//数据预热，从MySQL中读取秒杀商品和60条热卖商品（按销量降序）获取它们的详情，写入Redis
func GetHotDataInsertToRedis()  {

	var detailIidList []string  // 存放所有需要将'商品详情信息'放入Redis的商品iid的列表

	secIidList := dao.SelectGoodsIidFromSecKills()
	sellList := service.GetGoodsIndexListService("sell", "1")
	newList := service.GetGoodsIndexListService("new", "1")
	//sellList := dao.SelectGoodsIndexBySells(1, 10)
	// 首先将首页的商品列表信息存入redis
	for _, iid := range secIidList{
		detailIidList = append(detailIidList, iid)
		//type List struct {
		//	Image string `json:"image"`
		//	Title string `json:"title"` 从商品表中取出这几个字段，下面几个list里面都是存的这个
		//	Price string `json:"price"`
		//	Cfav  int    `json:"cfav"`	tempGoodsIndexItem就是左边这个结构体的一个实例
		//	Iid   string `json:"iid"`
		//}
		tempGoodsIndexItem :=  dao.SelectGoodsIndexByIid(iid)
		saveGoodsListToRedis(SecKillGoodsKey, tempGoodsIndexItem)
	}
	for _, listItem := range sellList{
		saveGoodsListToRedis(SellGoodsKey, listItem)
		detailIidList = append(detailIidList, listItem.Iid)
	}

	for _, listItem := range newList{
		saveGoodsListToRedis(NewGoodsKey, listItem)
		detailIidList = append(detailIidList, listItem.Iid)
	}

	// 根据iid将商品详情页的json结构体拿到，再插入Redis
	var goodsDetailRespJson respStruct.GoodsDetailRespJson
	for _, iid := range detailIidList {
		goodsDetailRespJson = service.GetGoodsDetailService(iid)
		goodsDetailRespJson.Success = true
		saveGoodsDetailToRedis(iid, goodsDetailRespJson)
	}

	// 插入秒杀商品库存到Redis中
	secKillGoodsList := dao.SelectGoodsInfoFromSecKills()

	for _, tempItem := range secKillGoodsList{
		saveGoodsStockToRedis(tempItem.Iid, tempItem.Stock)
	}

}
func saveGoodsStockToRedis(iid string, stock int64)  {
	stockStr := strconv.FormatInt(stock, 10)
	err := data.Client.Set(fmt.Sprintf(SecKillGoodsStockKey, iid), stockStr, -1).Err()
	if err != nil {
		log.Printf("saveGoodsStockToRedis : redis.Set() failed, err: %v", err)
	}

}

// 将商品详情写入Redis
func saveGoodsDetailToRedis(iid string, goodsDetailRespJson respStruct.GoodsDetailRespJson)  {
	goodsDetailRespJsonByte, err := json.Marshal(goodsDetailRespJson)
	if err != nil {
		log.Println("saveGoodsDetailToRedis json.Marshal err : ", err)
	}
	redisErr := data.Client.Set(fmt.Sprintf(GoodsDetailKey, iid),string(goodsDetailRespJsonByte), time.Duration(GoodsDetailExpire)*time.Second).Err()
	if  redisErr!=nil{
		log.Printf("redis.Set() failed, err: %v", redisErr)
	}
}

// 将商品主页的json结构体转成字符串push到redis的list中
func saveGoodsListToRedis(goodsListKey string,goodsListItem respStruct.List)  {
	goodsListItemJsonByte, err := json.Marshal(goodsListItem)
	if err != nil {
		log.Println("saveGoodsListToRedis json.Marshal err : ", err)
	}
	data.Client.RPush(goodsListKey, string(goodsListItemJsonByte))
	//if  redisErr!=nil{
	//	log.Printf("redis.Rpush() failed, err: %v", redisErr)
	//}
}

// 从redis中取出商品详情
func GetGoodsDetailFromRedis(iid string)(goodsDetailJson respStruct.GoodsDetailRespJson, isEmpty bool)  {
	log.Println("GetGoodsDetailFromRedis")
	goodsDetailJsonStr, err := data.Client.Get(fmt.Sprintf(GoodsDetailKey, iid)).Result()
	//if err == redis.Nil {
	//	err = nil
	//}else {
	//	log.Println("redis.Get() failed, err:", err)
	//}
	if len(goodsDetailJsonStr) == 0{
		isEmpty = true
		goodsDetailJson = respStruct.GoodsDetailRespJson{}
		return
	}
	err = json.Unmarshal([]byte(goodsDetailJsonStr), &goodsDetailJson)
	if err != nil {
		log.Println("GetGoodsDetailFromRedis Unmarshal err: ", err)

	}
	return  goodsDetailJson, false
}

//从Redis中取出商品主页列表信息 返回 []respStruct.List
func GetGoodsListFromRedis(goodsType string)(goodsList []respStruct.List, isEmpty bool ) {
	goodsKey := ""
	if goodsType == "new"{
		goodsKey= NewGoodsKey
	}else if goodsType == "sell" {
		goodsKey = SellGoodsKey
	}else if goodsType == "secKill" {
		goodsKey = SecKillGoodsKey
	}

	log.Println("GetGoodsListFromRedis, goodsKey: ", goodsKey)
	rLen, err := data.Client.LLen(goodsKey).Result()
	//log.Println(rLen, err)
	if rLen == 0{
		goodsList = []respStruct.List{}
		isEmpty = true
		return
	}

	//遍历
	lists, err := data.Client.LRange(goodsKey, 0, rLen-1).Result()

	var  tempListItem respStruct.List
	for _, listItem := range lists{
		err = json.Unmarshal([]byte(listItem), &tempListItem)
		if err != nil {
			log.Println("GetGoodsDetailFromRedis Unmarshal err: ", err)
		}
		goodsList = append(goodsList, tempListItem)
	}
	isEmpty = false
	return
}