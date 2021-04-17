package controller

import (
	"GraduateDesign/api/cache"
	"GraduateDesign/api/service"
	"GraduateDesign/model/reqStruct"
	"GraduateDesign/model/respStruct"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// 发布商品的处理函数
func PubGoods(ctx *gin.Context) {
	// 将发送过来的json转成预先定义好的结构体
	goodsReq := &reqStruct.GoodsJson{}
	ctx.BindJSON(&goodsReq)
	service.PubGoodsService(goodsReq)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// 获取首页数据的处理函数
func GetGoodsIndexList(ctx *gin.Context) {
	//name := c.DefaultQuery("name", "枯藤")
	//c.String(http.StatusOK, fmt.Sprintf("hello %s", name))
	goodsType := ctx.DefaultQuery("type", "sell")
	page := ctx.DefaultQuery("page", "1")

	var goodsIndexRespJson respStruct.GoodsIndexRespJson
	var isEmpty = true
	if page == "1"{
		goodsIndexRespJson.Data.List, isEmpty = cache.GetGoodsListFromRedis(goodsType)
	}
	if isEmpty==true{
		log.Printf("Redis没有缓存%s商品第%s页数据 从数据库中取", goodsType, page)
		goodsIndexRespJson.Data.List = service.GetGoodsIndexListService(goodsType, page)
	}
	goodsIndexRespJson.Success = true

	goodsIndexRespJsonBytes, err := json.Marshal(goodsIndexRespJson)
	//goodsIndexRespJsonStr := string(goodsIndexRespJsonBytes)
	//fmt.Println(goodsIndexRespJsonStr)
	if err != nil {
		log.Panic("生成goodsIndexRespJsonBytes时错误： ", err)
	}
	//ctx.JSON(http.StatusOK, goodsIndexRespJsonStr)
	// 上面这种返回方式有问题，会返回一个字符串，浏览器上会看到一堆的转义符
	// 正确的方式可见：https://ask.csdn.net/questions/1028522
	ctx.Data(http.StatusOK, "application/json", goodsIndexRespJsonBytes)
	//c.JSON(http.StatusOK, gin.H{
	//	"status" :200,
	//	"error": nil,
	//	"data": persons, 当然，可以直接写上变量不用Marshal，但是变量名容易写错，所以用data方式返回还是更好
	//})
}

// 获取首页的滚动栏和推荐栏，由于这部分数据不是系统的重点，所以没存入数据库，直接硬编码返回哈哈哈哈，偷个懒~
func GetHomeMultiData(ctx *gin.Context) {
	jsonData := []byte(`{
    "data": {
        "recommend": {
            "list": [{
                "image": "https://s10.mogucdn.com/mlcdn/c45406/180913_036dli57aah85cb82l1jj722g887g_225x225.png",
                "link": "http://act.meilishuo.com/10dianlingquan?acm=3.mce.2_10_1dggc.13730.0.ccy5br4OlfK0U.pos_0-m_313898-sd_119",
                "title": "十点抢券"
            }, {
                "image": "https://s10.mogucdn.com/mlcdn/c45406/180913_25e804lk773hdk695c60cai492111_225x225.png",
                "link": "https://act.mogujie.com/tejiazhuanmai09?acm=3.mce.2_10_1dgge.13730.0.ccy5br4OlfK0V.pos_1-m_313899-sd_119",
                "title": "好物特卖"
            }, {
                "image": "https://s10.mogucdn.com/mlcdn/c45406/180913_387kgl3j21ff29lh04181iek48a6h_225x225.png",
                "link": "http://act.meilishuo.com/neigouful001?acm=3.mce.2_10_1b610.13730.0.ccy5br4OlfK0W.pos_2-m_260486-sd_119",
                "title": "内购福利"
            }, {
                "image": "https://s10.mogucdn.com/mlcdn/c45406/180913_8d4e5adi8llg7c47lgh2291akiec7_225x225.png",
                "link": "http://act.meilishuo.com/wap/yxzc1?acm=3.mce.2_10_1dggg.13730.0.ccy5br4OlfK0X.pos_3-m_313900-sd_119",
                "title": "初秋上新"
            }]
        },
        "banner": {
            "list": [{
                "image": "https://s10.mogucdn.com/mlcdn/c45406/180926_45fkj8ifdj4l824l42dgf9hd0h495_750x390.jpg",
                "link": "https://act.mogujie.com/huanxin0001?acm=3.mce.2_10_1jhwa.43542.0.ccy5br4OlfK0Q.pos_0-m_454801-sd_119"
            }, {
                "image": "https://s10.mogucdn.com/mlcdn/c45406/180926_31eb9h75jc217k7iej24i2dd0jba3_750x390.jpg",
                "link": "https://act.mogujie.com/ruqiu00001?acm=3.mce.2_10_1ji16.43542.0.ccy5br4OlfK0R.pos_1-m_454889-sd_119"
            }, {
                "image": "https://s10.mogucdn.com/mlcdn/c45406/180919_3f62ijgkj656k2lj03dh0di4iflea_750x390.jpg",
                "link": "https://act.mogujie.com/huanji001?acm=3.mce.2_10_1jfj8.43542.0.ccy5br4OlfK0S.pos_2-m_453270-sd_119"
            }, {
                "image": "https://s10.mogucdn.com/mlcdn/c45406/180917_18l981g6clk33fbl3833ja357aaa0_750x390.jpg",
                "link": "https://act.mogujie.com/liuxing00001?acm=3.mce.2_10_1jepe.43542.0.ccy5br4OlfK0T.pos_3-m_452733-sd_119"
            }]
        }
    },
    "success": true
}`)

	ctx.Data(http.StatusOK, "application/json", jsonData)
}

// 获取商品详情数据
func GetGoodsDetail(ctx *gin.Context) {
	iid := ctx.Query("iid")
	var goodsDetailRespJson respStruct.GoodsDetailRespJson
	var isEmpty = true
	// 从redis中取商品详情，如果取到了会将isEmpty置为false
	goodsDetailRespJson, isEmpty = cache.GetGoodsDetailFromRedis(iid)
	if isEmpty == true{
		// 说明redis中没有缓存该商品的详情信息
		log.Println("redis中没有缓存该商品的详情信息，从MySQL中取...")
		goodsDetailRespJson = service.GetGoodsDetailService(iid)
	}
	goodsDetailRespJson.Success = true

	goodsDetailRespJsonByte, err := json.Marshal(goodsDetailRespJson)
	if err != nil {
		log.Panic("生成goodsDetailRespJsonByte时错误： ", err)
	}
	ctx.Data(http.StatusOK, "application/json", goodsDetailRespJsonByte)

}

// 获取推荐商品数据
func GetGoodsRecommend(ctx *gin.Context) {
	var goodsIndexRespJson respStruct.GoodsIndexRespJson
	goodsIndexRespJson.Data.List = service.GetGoodsIndexListService("sell", "10")
	goodsIndexRespJson.Success = true

	goodsIndexRespJsonBytes, err := json.Marshal(goodsIndexRespJson)

	if err != nil {
		log.Panic("生成goodsIndexRespJsonBytes时错误： ", err)
	}

	ctx.Data(http.StatusOK, "application/json", goodsIndexRespJsonBytes)
}
