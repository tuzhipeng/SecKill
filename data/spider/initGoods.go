package spider

import (
	"GraduateDesign/model/reqStruct"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

//var postUrl = "127.0.0.1:9999/api/goods"
var postUrl = "http://47.93.12.71:9999/api/goods"
var successCount = 1
var fatalCount = 1
var wg sync.WaitGroup

// 插入一些爬取得来的json数据到数据库
func InitGoods(fileList []string) {

	//fmt.Println(fileList)
	// 开始计时，计算整个插入环节耗时
	start := time.Now()

	fmt.Println("开始读取json文件中的内容赋值到结构体...")
	for _, filePath := range fileList[900:] {
		wg.Add(1)
		//
		//filePath := fileList[0]
		go func(filePath string) {
			defer wg.Done()

			var goodsImage = reqStruct.Goodsimage{}
			var goodsShop = reqStruct.Goodsshop{}
			var goodsComment = reqStruct.Goodscomment{}
			var goods = reqStruct.Goods{}

			goodsImage = buildGoodsImage(filePath, goodsImage)
			goodsShop = buildGoodsShop(filePath, goodsShop)
			goodsComment = buildGoodsComment(filePath, goodsComment)
			goods = buildGoods(filePath, goods)
			//fmt.Println(goodsImage)
			// 构建json
			var goodsJson = reqStruct.GoodsJson{goodsImage, goodsShop, goodsComment, goods}
			goodsJsonStr, err := json.Marshal(goodsJson)
			if err != nil {
				log.Println("goodsJsonStr生成错误： ", err.Error())
			}

			// 将构建好的json体post到接口上，从而插入到数据库
			err = DoPostJson(goodsJsonStr)
			if err != nil {
				log.Printf("第 %d 次doPostJson时发生错误：%v\n", fatalCount, err)
				fatalCount++
			}

		}(filePath)

	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("耗时： ", elapsed)
	// 串行插入耗时50分钟，并发插入数据库耗时1分40秒！

}

//构建商品图片结构体
func buildGoodsImage(filePath string, goodsImage reqStruct.Goodsimage) reqStruct.Goodsimage {
	// gojsonq每次只能挂载一个node，所以每次都得reset一下,好坑！
	// 构建topImages
	jq := gojsonq.New().File(filePath)
	topImagesRes := jq.From("result.itemInfo.topImages").Get().([]interface{})
	for _, imageUrl := range topImagesRes {
		goodsImage.Topimages = append(goodsImage.Topimages, imageUrl.(string))
	}
	// 构建detailImages

	jq.Reset()
	detailImageRes, _ := jq.From("result.detailInfo.detailImage").Select("list").Get().([]interface{})
	//jq := gojsonq.New().File(filePath).From("result.detailInfo.detailImage").Select("list")
	//detailImageRes, _ := jq.Get().([]interface{})
	//fmt.Println("detailImageRes: ", detailImageRes)

	for _, detailInfo := range detailImageRes {

		detailMap := detailInfo.(map[string]interface{})
		//fmt.Println("detailMap", detailMap)
		//fmt.Println("detailMap KEY: ", detailMap["list"])
		//遍历list数组
		for _, imageUrl := range detailMap["list"].([]interface{}) {
			//fmt.Println("url :"+ imageUrl.(string))
			goodsImage.Detailimage = append(goodsImage.Detailimage, imageUrl.(string))
		}
	}
	//fmt.Println(goodsImage)
	return goodsImage

}

//构建商品店铺结构体
func buildGoodsShop(filePath string, goodsShop reqStruct.Goodsshop) reqStruct.Goodsshop {
	jq := gojsonq.New().File(filePath)
	//jq := gojsonq.New().File(filePath).From("result.shopInfo.score").Select("name","score")
	//scoreListRes, _ := jq.Get().([]interface{})
	scoreListRes, _ := jq.From("result.shopInfo.score").Select("name", "score").Get().([]interface{})
	var temScore reqStruct.Score

	for _, scoreInfo := range scoreListRes {
		scoreMap := scoreInfo.(map[string]interface{})

		temScore.Name = scoreMap["name"].(string)
		temScore.Score = scoreMap["score"].(float64)
		goodsShop.Score = append(goodsShop.Score, temScore)
	}

	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("读取文件失败:%#v", err)
	}
	jsonStr := string(jsonFile)
	// 临时使用gjson是因为运行下面这条代码的时候读取到的是float64类型，转不了int64
	// 所以才调研用到gjson来救火哈哈哈哈，前面用的gojsonq就不改啦，没时间~
	//goodsShop.Cfans = gojsonq.New().File(filePath).Find("result.shopInfo.cFans")
	goodsShop.Cfans = gjson.Get(jsonStr, "result.shopInfo.cFans").Int()
	goodsShop.Csells = gjson.Get(jsonStr, "result.shopInfo.cSells").Int()
	goodsShop.Cgoods = gjson.Get(jsonStr, "result.shopInfo.cGoods").Int()

	goodsShop.Shopid = gjson.Get(jsonStr, "result.shopInfo.shopId").String()
	goodsShop.Shoplogo = gjson.Get(jsonStr, "result.shopInfo.shopLogo").String()
	goodsShop.Name = gjson.Get(jsonStr, "result.shopInfo.name").String()
	goodsShop.Userid = gjson.Get(jsonStr, "result.shopInfo.userId").String()
	return goodsShop

}

//构建商品评论结构体
func buildGoodsComment(filePath string, goodsComment reqStruct.Goodscomment) reqStruct.Goodscomment {
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("读取文件失败:%#v", err)
	}
	jsonStr := string(jsonFile)

	var tempComment reqStruct.List
	tempComment.User.Uname = gjson.Get(jsonStr, "result.rate.list.0.user.uname").String()
	tempComment.User.Avatar = gjson.Get(jsonStr, "result.rate.list.0.user.avatar").String()
	tempComment.Created = gjson.Get(jsonStr, "result.rate.list.0.created").Int()
	tempComment.Style = gjson.Get(jsonStr, "result.rate.list.0.style").String()
	tempComment.Content = gjson.Get(jsonStr, "result.rate.list.0.content").String()
	goodsComment.List = append(goodsComment.List, tempComment)

	return goodsComment
}

// 构建商品结构体
func buildGoods(filePath string, goods reqStruct.Goods) reqStruct.Goods {
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("读取文件失败:%#v", err)
	}
	jsonStr := string(jsonFile)

	goods.Iid = gjson.Get(jsonStr, "result.itemInfo.iid").String()
	goods.Shopid = gjson.Get(jsonStr, "result.itemInfo.shopId").String()
	goods.Discountdesc = gjson.Get(jsonStr, "result.itemInfo.discountDesc").String()
	goods.Price = gjson.Get(jsonStr, "result.itemInfo.price").String()
	goods.Oldprice = gjson.Get(jsonStr, "result.itemInfo.oldPrice").String()
	goods.Lownowprice = gjson.Get(jsonStr, "result.itemInfo.lowNowPrice").String()
	goods.Desc = gjson.Get(jsonStr, "result.itemInfo.desc").String()
	goods.Title = gjson.Get(jsonStr, "result.itemInfo.title").String()
	goods.Image = gjson.Get(jsonStr, "result.itemInfo.topImages.0").String()

	var valid = regexp.MustCompile("[0-9]+")
	// [
	//    "销量 5013",
	//    "收藏49人",    原始数据如左边，要提取销量数和收藏数
	//    "默认快递"
	//]
	sellsStr := gjson.Get(jsonStr, "result.columns.0").String()
	sells, _ := strconv.Atoi(valid.FindAllString(sellsStr, -1)[0])
	goods.Sells = int64(sells) // 这个参数

	cfavStr := gjson.Get(jsonStr, "result.columns.1").String()
	cfav, _ := strconv.Atoi(valid.FindAllString(cfavStr, -1)[0])
	goods.Cfav = int64(cfav)

	goods.Delivery = gjson.Get(jsonStr, "result.columns.2").String()
	// 库存参数默认给着
	goods.Stock = int64(10000)

	return goods
}

// 发送带有json体的post请求
func DoPostJson(postJsonStr []byte) error {
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(postJsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")




	// 发送请求
	client := &http.Client{Timeout: time.Duration( 100*time.Second )}
	resp, err := client.Do(req)


	if err != nil {
		log.Panic("doPost err :", err)
		return err
	}
	defer resp.Body.Close()
	//fmt.Println("response Status:", resp.StatusCode)
	//fmt.Println("response Headers:", resp.Header)
	if resp.StatusCode != http.StatusOK {
		log.Panic("statusCodeErr : ", resp.StatusCode)
		return errors.New("statusCodeErr")
	}
	fmt.Printf("成功插入第 %d 条商品\n", successCount)
	successCount++
	return nil
}
