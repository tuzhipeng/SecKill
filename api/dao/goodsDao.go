package dao

import (
	"GraduateDesign/data"
	"GraduateDesign/model"
	"GraduateDesign/model/reqStruct"
	"GraduateDesign/model/respStruct"
	"log"
	"unicode/utf8"
)

//type Result struct {
//	Name string
//	Age  int
//}
//
//var result Result
//db.Table("users").Select("name, age").Where("name = ?", 3).Scan(&result)
//
//// Raw SQL
//db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)

// 做分页功能，获取取指page，指定pagesize的记录
//分页思路：https://blog.csdn.net/weixin_44854027/article/details/105388471
//gorm实现：https://blog.csdn.net/sss996/article/details/93891503
func SelectGoodsIndexBySells(page int, pageSize int) []respStruct.List {
	var goodsIndexList []respStruct.List
	data.Db.Table("goods").Select("iid, cfav, price, title, image").Order("sells").Limit(pageSize).Offset((page - 1) * pageSize).Scan(&goodsIndexList)
	return goodsIndexList
}

func SelectGoodsIndexBySellsDesc(page int, pageSize int) []respStruct.List {
	var goodsIndexList []respStruct.List
	data.Db.Table("goods").Select("iid, cfav, price, title, image").Order("sells desc").Limit(pageSize).Offset((page - 1) * pageSize).Scan(&goodsIndexList)
	return goodsIndexList
}

func SelectGoodsIndexByIid(iid string) respStruct.List {
	var goodsIndexList respStruct.List
	data.Db.Table("goods").Select("iid, cfav, price, title, image").Where("iid = ?", iid).Scan(&goodsIndexList)
	return goodsIndexList
}

func SelectGoodsIidFromSecKills() []string {
	type iidItem struct{
		Iid string
	}
	var iidList []iidItem
	data.Db.Table("goods_sec_kills").Select("iid").Scan(&iidList)
	var  resList []string
	for _, iiditem := range iidList{
		resList = append(resList, iiditem.Iid)
	}
	return resList
}

type SecKillGoodsItem struct{
	Iid string
	Stock int64
}
func SelectGoodsInfoFromSecKills() []SecKillGoodsItem {

	var secKillGoodsList []SecKillGoodsItem
	data.Db.Table("goods_sec_kills").Select("iid, stock").Scan(&secKillGoodsList)
	//var  resList []string
	//for _, secKillGoods := range secKillGoodsList{
	//	resList = append(resList, secKillGoods.Iid)
	//}
	return secKillGoodsList
}

// 通过商品的iid获取商品图片信息
func SelectGoodsImageByIid(iid string) []model.GoodsImage {
	var goodsImages []model.GoodsImage
	data.Db.Where("iid = ?", iid).Find(&goodsImages)
	return goodsImages
}

// 通过商品的iid获取商品评价信息
func SelectGoodsCommentByIid(iid string) model.GoodsComment {
	var goodsComment model.GoodsComment
	data.Db.Where("iid = ?", iid).First(&goodsComment)
	return goodsComment
}

// 通过商品的iid获取商品店铺信息
func SelectGoodsShopBySid(sid string) model.GoodsShop {
	var goodsShop model.GoodsShop
	data.Db.Where("sid = ?", sid).First(&goodsShop)
	return goodsShop
}

// 通过商品的iid获取商品信息
func SelectGoodsByIid(iid string) model.Goods {
	var goods model.Goods
	data.Db.Where("iid = ?", iid).First(&goods)
	return goods
}



//构建商品图片表中的数据
func InsertGoodsImage(iid string, goodsImage reqStruct.Goodsimage) {
	var goodsImageData = &model.GoodsImage{}
	goodsImageData.Iid = iid
	for _, image := range goodsImage.Topimages {
		// 滚动栏的展示图片，Ikind为0
		goodsImageData.Image = image
		goodsImageData.Ikind = 0
		//fmt.Println("imageURL是： " +image)

		err := data.Db.Create(&goodsImageData).Error
		if err != nil {
			log.Panic("Insert goodsImageData Topimages failed")
		}
	}

	for _, image := range goodsImage.Detailimage {
		//  竖栏的展示图片，Ikind为1
		goodsImageData.Image = image
		goodsImageData.Ikind = 1

		err := data.Db.Create(&goodsImageData).Error
		if err != nil {
			log.Panic("Insert goodsImageData Detailimage failed: ", err.Error())

		}
	}
}

//构建商品评论表中的数据
func InsertGoodsComment(iid string, goodsComment reqStruct.Goodscomment) {
	var goodsCommentData = &model.GoodsComment{}
	goodsCommentData.Iid = iid
	goodsCommentData.Uname = goodsComment.List[0].User.Uname
	goodsCommentData.Avatar = goodsComment.List[0].User.Avatar

	//fmt.Println(goodsComment.DetailList[0].Content)
	//打印出来，字符串是：“客服立马发了链接过来。反正很不错 点赞👍 ”
	//评论中有表情, 无法插入数据库，直接去掉了表情（我太难了。。。。
	//goodsCommentData.Content = string("很满意！")
	goodsCommentData.Content = filterEmoji(goodsComment.List[0].Content)

	goodsCommentData.Style = goodsComment.List[0].Style
	goodsCommentData.Created = goodsComment.List[0].Created

	err := data.Db.Create(&goodsCommentData).Error
	if err != nil {
		log.Panic("Insert InsertGoodsComment failed: ", err.Error())

	}
}

//构建商品店铺表中的数据
func InsertGoodsShop(goodsShop reqStruct.Goodsshop) {
	var goodsShopData = &model.GoodsShop{}
	goodsShopData.Uid = goodsShop.Userid
	goodsShopData.Sid = goodsShop.Shopid
	goodsShopData.ShopLogo = goodsShop.Shoplogo
	goodsShopData.Name = filterEmoji(goodsShop.Name)
	goodsShopData.CFans = goodsShop.Cfans
	goodsShopData.CSells = goodsShop.Csells
	goodsShopData.CGoods = goodsShop.Cgoods

	for _, score := range goodsShop.Score {
		if score.Name == "描述相符" {
			goodsShopData.DescScore = score.Score
		}
		if score.Name == "价格合理" {
			goodsShopData.PriceScore = score.Score
		}
		if score.Name == "质量满意" {
			goodsShopData.QualityScore = score.Score
		}

	}

	// 多个商品可能对应一个店铺，如果店铺已经存在就不用再创建了
	var tempGoodsShopData = &model.GoodsShop{}
	data.Db.Where(" sid = ?", goodsShopData.Sid).First(tempGoodsShopData)
	if tempGoodsShopData.Sid == "" {
		err := data.Db.Create(&goodsShopData).Error
		if err != nil {
			log.Panic("Insert goodsShopData failed: ", err.Error())
		}
	}

}

//构建商品表中的数据
func InsertGoods(sid string, goods reqStruct.Goods) {
	var goodsData = &model.Goods{}
	goodsData.Sid = sid
	goodsData.Iid = goods.Iid
	goodsData.Image = goods.Image
	goodsData.Title = filterEmoji(goods.Title)
	goodsData.Cfav = goods.Cfav
	goodsData.OldPrice = goods.Oldprice
	goodsData.Price = goods.Price
	goodsData.LowNowPrice = goods.Lownowprice
	goodsData.Desc = filterEmoji(goods.Desc)
	goodsData.DiscountDesc = filterEmoji(goods.Discountdesc)
	goodsData.Stock = goods.Stock
	goodsData.Delivery = goods.Delivery
	goodsData.Sells = goods.Sells

	err := data.Db.Create(&goodsData).Error
	if err != nil {
		log.Panic("Insert goodsData failed :", err.Error())
	}
}

// 去除字符串中的表情，由于表情是4个字节，插入数据库时会报错，所以过滤一下
func filterEmoji(content string) string {
	newContent := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			newContent += string(value)
		}
	}
	return newContent
}
