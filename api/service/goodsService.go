package service

import (
	"GraduateDesign/api/dao"
	"GraduateDesign/model/reqStruct"
	"GraduateDesign/model/respStruct"
	"fmt"
	"log"
	"strconv"
)

func PubGoodsService(goodsReq *reqStruct.GoodsJson) {
	// 下面三张表都需要goods的iid
	var iid = goodsReq.Goods.Iid
	dao.InsertGoodsImage(iid, goodsReq.Goodsimage)
	dao.InsertGoodsComment(iid, goodsReq.Goodscomment)
	dao.InsertGoodsShop(goodsReq.Goodsshop)
	// 商品表需要店铺ID
	var sid = goodsReq.Goodsshop.Shopid
	dao.InsertGoods(sid, goodsReq.Goods)

}

// 通过传来的获取商品类型和页码，从数据库中查询出对应的商品
func GetGoodsIndexListService(goodsType string, pageStr string) []respStruct.List {
	var goodsIndexList []respStruct.List
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Panic("pageStr transform to int err: ", err)
	}
	if goodsType == "sell" {
		// 查找热卖商品，默认30条
		//goodsIndexList = dao.SelectGoodsIndexBySellsDesc(page, 30)
		goodsIndexList = dao.SelectGoodsIndexBySellsDesc(page, 30)
	}
	if goodsType == "new" {
		// 查找上新商品，默认30条
		//goodsIndexList = dao.SelectGoodsIndexBySells(page, 30)
		goodsIndexList = dao.SelectGoodsIndexBySells(page, 30)
	}
	//if goodsType == "secKill" {
	//	// 查找秒杀商品
	//	secIidList := dao.SelectGoodsIidFromSecKills()
	//	for _, iid := range secIidList{
	//		tempGoodsIndexItem :=  dao.SelectGoodsIndexByIid(iid)
	//		goodsIndexList = append(goodsIndexList, tempGoodsIndexItem)
	//	}
	//}
	return goodsIndexList

}

// 通过传入的iid，构建应该返回的商品详情json
func GetGoodsDetailService(iid string) respStruct.GoodsDetailRespJson {
	var goodsDetailRespJson respStruct.GoodsDetailRespJson
	// 构建图片信息，后面赋值
	var tempTopImages, tempDetailImages []string
	for _, goodsImageItem := range dao.SelectGoodsImageByIid(iid) {

		if goodsImageItem.Ikind == 0 {
			tempTopImages = append(tempTopImages, goodsImageItem.Image)
		} else if goodsImageItem.Ikind == 1 {
			tempDetailImages = append(tempDetailImages, goodsImageItem.Image)
		}
	}

	// 构建评论信息(rate)
	commentInfo := dao.SelectGoodsCommentByIid(iid)
	var tempCommentItem respStruct.CommentItem
	tempCommentItem.User.Uname = commentInfo.Uname
	tempCommentItem.User.Avatar = commentInfo.Avatar
	tempCommentItem.Content = commentInfo.Content
	tempCommentItem.Style = commentInfo.Style
	tempCommentItem.Created = commentInfo.Created
	goodsDetailRespJson.Result.Rate.List = append(goodsDetailRespJson.Result.Rate.List, tempCommentItem)

	// 构建商品信息(itemInfo和detailInfo)
	var itemInfo respStruct.Iteminfo
	goodsInfo := dao.SelectGoodsByIid(iid)
	itemInfo.Topimages = tempTopImages
	itemInfo.Desc = goodsInfo.Desc
	itemInfo.Title = goodsInfo.Title
	itemInfo.Discountdesc = goodsInfo.DiscountDesc
	itemInfo.Oldprice = goodsInfo.OldPrice
	itemInfo.Price = goodsInfo.Price
	itemInfo.Lownowprice = goodsInfo.LowNowPrice
	goodsDetailRespJson.Result.Iteminfo = itemInfo

	var detailImage respStruct.Detailimage
	detailImage.Desc = itemInfo.Desc
	detailImage.Key = "穿着效果"
	detailImage.Anchor = "model_img"
	detailImage.List = tempDetailImages
	goodsDetailRespJson.Result.Detailinfo.Desc = itemInfo.Desc
	goodsDetailRespJson.Result.Detailinfo.Detailimage = append(goodsDetailRespJson.Result.Detailinfo.Detailimage, detailImage)

	// 构建店铺信息(shopInfo)
	var tempScoreList []respStruct.Score
	shopInfo := dao.SelectGoodsShopBySid(goodsInfo.Sid)
	tempScoreList = append(tempScoreList, respStruct.Score{Name: "质量满意", Score: shopInfo.QualityScore, Isbetter: shopInfo.QualityScore > 4.5})
	tempScoreList = append(tempScoreList, respStruct.Score{Name: "描述相符", Score: shopInfo.DescScore, Isbetter: shopInfo.DescScore > 4.5})
	tempScoreList = append(tempScoreList, respStruct.Score{Name: "价格合理", Score: shopInfo.PriceScore, Isbetter: shopInfo.PriceScore > 4.5})
	goodsDetailRespJson.Result.Shopinfo.Score = tempScoreList
	goodsDetailRespJson.Result.Shopinfo.Cfans = shopInfo.CFans
	goodsDetailRespJson.Result.Shopinfo.Csells = shopInfo.CSells
	goodsDetailRespJson.Result.Shopinfo.Shoplogo = shopInfo.ShopLogo
	goodsDetailRespJson.Result.Shopinfo.Name = shopInfo.Name
	goodsDetailRespJson.Result.Shopinfo.Cgoods = shopInfo.CGoods
	goodsDetailRespJson.Result.Shopinfo.Shopid = shopInfo.Sid

	// 构建column信息
	columns := []string{fmt.Sprintf("销量 %d", goodsInfo.Sells), fmt.Sprintf("收藏数 %d", goodsInfo.Cfav), goodsInfo.Delivery}
	goodsDetailRespJson.Result.Columns = columns

	return goodsDetailRespJson
}

// 获取推荐数据
func GetGoodsRecommendService() {
	// TODO 做一点点推荐

}
