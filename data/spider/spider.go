package spider

import (
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

// 本文件的作用： 爬取后端接口的json数据并保存至本地文件
// 这个是那个Vue项目的后端接口
//var baseInterfaceUrl = "http://152.136.185.210:7878/api/m5"
//var baseJsonPath = "./data/spider/json"

func fetch(url string) string {
	fmt.Println("Fetch Url", url)
	client := &http.Client{}
	// 构建请求，使用浏览器引擎
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		return ""
	}
	defer resp.Body.Close()

	//将请求的body转换成字符串返回
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error", err)
		return ""
	}
	return string(body)
}

//typeList := []string{"pop", "new", "sell"}
//startPage := 1
//endPage := 50
// 根据商品类型和分页，爬取各主页的内容
func SpiderPage(typeList []string, startPage int, endPage int, baseInterfaceUrl string, baseJsonPath string) {
	//basePageUrl := "http://152.136.185.210:7878/api/m5/home/data?type=%s&page=%d"
	basePageUrl := baseInterfaceUrl + "/home/data?type=%s&page=%d"
	basePageFileName := baseJsonPath + "/pageJson/%s_%d.json"
	// 构建详情URL，分别爬取["pop","new","sell"] 1~50页的内容

	// 开始计时，计算整个爬取环节耗时
	start := time.Now()

	// 并发爬取，用waitGroup做阻塞
	var wg sync.WaitGroup

	for _, goodsType := range typeList {
		// 由于目标URL的sell数据只有20页，所以要特殊处理一下
		if goodsType == "sell" && endPage > 20 {
			endPage = 20
		}
		for page := startPage; page <= endPage; page++ {
			wg.Add(1)
			// 采用goroutine并发爬取，主意将page和type的值传入，不然值会错乱
			go func(goodsType string, page int) {
				defer wg.Done()
				targetUrl := fmt.Sprintf(basePageUrl, goodsType, page)
				result := fetch(targetUrl)

				//创建文件，保存爬取下来的数据
				fileName := fmt.Sprintf(basePageFileName, goodsType, page)
				dstFile, err := os.Create(fileName)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				// 写入后关闭文件
				dstFile.WriteString(result)
				dstFile.Close()
			}(goodsType, page)

		}
	}
	// 阻塞，让子协程完成
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("耗时： ", elapsed)
	// 经过测试，并发版爬虫耗时4.55秒， 串行爬虫写入用时13.04秒

}

//typeList := []string{"pop", "new", "sell"}
//startPage := 1
//endPage := 50
// 根据商品的iid爬取详情页的内容
func SpiderDetail(typeList []string, startPage int, endPage int, baseInterfaceUrl string, baseJsonPath string) {
	// 根据每一页的内容提取其中商品列表中的iid，再去爬取详情内容
	//baseDetailUrl := "http://152.136.185.210:7878/api/m5/detail?iid=%s"
	baseDetailUrl := baseInterfaceUrl + "/detail?iid=%s"
	basePageFileName := baseJsonPath + "/pageJson/%s_%d.json"
	baseDetailFileName := baseJsonPath + "/detailJson/%s.json"

	// 开始计时，计算整个爬取环节耗时
	start := time.Now()

	// 并发爬取，用waitGroup做阻塞
	//var wg sync.WaitGroup

	for _, goodsType := range typeList {
		// 由于目标URL的sell数据只有20页，所以要特殊处理一下
		if goodsType == "sell" && endPage > 20 {
			endPage = 20
		}
		for i := startPage; i <= endPage; i++ {
			// 通过gojsonq提取出json中想要的字段，简单使用见：https://blog.csdn.net/yjp19871013/article/details/83035588
			// 补充：其实这里用gjson更简单，懒得改了，可以参考：https://zhuanlan.zhihu.com/p/116076040
			targetFilePath := fmt.Sprintf(basePageFileName, goodsType, i)
			jq := gojsonq.New().File(targetFilePath).From("data.list").Select("iid")

			itemList, _ := jq.Get().([]interface{})
			for _, item := range itemList {
				//fmt.Println( fmt.Sprintf(baseDetailUrl, item.(map[string]interface{}) ["iid"].(string)) )
				//wg.Add(1)
				iid := item.(map[string]interface{})["iid"].(string)
				targetUrl := fmt.Sprintf(baseDetailUrl, iid)

				// 并发爬取数据
				//go func(targetUrl string, iid string) {
				//	defer wg.Done()
				result := fetch(targetUrl)

				//创建文件，保存爬取下来的数据
				fileName := fmt.Sprintf(baseDetailFileName, iid)
				dstFile, err := os.Create(fileName)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				// 写入后关闭文件
				dstFile.WriteString(result)
				dstFile.Close()
				//}(targetUrl, iid)
			}

		}
	}

	// 阻塞，让子协程完成
	//wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("耗时： ", elapsed)
	// 经过测试串行爬虫耗时166s，
	// 并发爬虫有点小问题哈哈哈哈，爬取数据写入文件是没问题的就是fetch URL时候会无响应，导致无法计时

}
