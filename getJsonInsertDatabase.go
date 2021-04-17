package main

import (
	"GraduateDesign/data/spider"
	"GraduateDesign/utils"
	"fmt"
)

// 爬取接口的json数据写入到文件，然后读取文件，提取需要的字段，插入数据库
func main() {

	var baseInterfaceUrl = "http://152.136.185.210:7878/api/m5"
	var baseJsonPath = "./data/spider/json"
	//var baseJsonPath = "./data/spider/testJson"
	typeList := []string{"pop", "new", "sell"}
	startPage := 1
	endPage := 50

	// 爬取json数据
	spider.SpiderPage(typeList, startPage, endPage, baseInterfaceUrl, baseJsonPath)
	spider.SpiderDetail(typeList, startPage, endPage, baseInterfaceUrl, baseJsonPath)
	//获取文件下的所有文件
	var jsonFilePath = baseJsonPath + "/detailJson"
	fileList, err := utils.GetAllFile(jsonFilePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	// 插入数据库
	spider.InitGoods(fileList)
}

//并发爬取主页31.6秒， 串行爬取详情8分钟，并发插入数据库2分16秒
