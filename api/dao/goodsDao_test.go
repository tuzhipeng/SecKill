package dao

import (
	"GraduateDesign/data/spider"
	"GraduateDesign/model"
	"GraduateDesign/model/respStruct"
	"GraduateDesign/utils"
	"fmt"
	"testing"
)

var baseInterfaceUrl = "http://152.136.185.210:7878/api/m5"
var baseJsonPath = "testData"

//var baseJsonPath = "./data/spider/testJson"
var typeList = []string{"pop"}
var startPage = 1
var endPage = 1

func TestMain(m *testing.M) {
	m.Run()

}

func TestGoodsWorkFlow(t *testing.T) {
	t.Run("测试：爬取第一页的商品list数据", testSpiderPage)
	t.Run("测试：根据商品list文件爬取里面所有商品的详情数据", testSpiderDetail)
	t.Run("测试：获取商品详情文件,获取所有商品，提取需要的字段插入数据库", testInitGoods)

	t.Run("测试：通过销量获取商品，每次获取三个", testSelectGoodsIndexBySells)
	t.Run("测试：通过iid获取商品的图片", testSelectGoodsImageByIid)
	t.Run("测试：通过iid获取商品基本信息", testSelectGoodsByIid)

}

func testSpiderPage(t *testing.T) {
	// 爬取json数据
	spider.SpiderPage(typeList, startPage, endPage, baseInterfaceUrl, baseJsonPath)
}
func testSpiderDetail(t *testing.T) {
	spider.SpiderDetail(typeList, startPage, endPage, baseInterfaceUrl, baseJsonPath)
}

func testInitGoods(t *testing.T) {
	//获取文件下的所有文件
	var jsonFilePath = baseJsonPath + "/detailJson"
	fileList, err := utils.GetAllFile(jsonFilePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	// 插入数据库
	spider.InitGoods(fileList)
}

func testSelectGoodsIndexBySells(t *testing.T) {
	var goodsIndexList []respStruct.List
	goodsIndexList = SelectGoodsIndexBySells(1, 3)

	for _, goodsIndexListItem := range goodsIndexList {
		fmt.Println(goodsIndexListItem)
	}
}

func testSelectGoodsImageByIid(t *testing.T) {
	var goodsImages []model.GoodsImage
	goodsImages = SelectGoodsImageByIid("1m70y5k")
	for _, goodsImagesItem := range goodsImages {
		fmt.Println("获取的goodsImage: ", goodsImagesItem)
	}
}

func testSelectGoodsByIid(t *testing.T) {
	var goods model.Goods
	goods = SelectGoodsByIid("1m70y5k")
	fmt.Println("获取goods: ", goods)
}
