package dao

import (
	"GraduateDesign/data"
	"GraduateDesign/model"
	"GraduateDesign/model/reqStruct"
	"fmt"
	"github.com/gitstliu/go-id-worker"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

func InsertOrderByMessage(orderMessage *reqStruct.OrderMessage) (err error) {
	fmt.Println("传入的orderMessage：", orderMessage)
	currWoker := &idworker.IdWorker{}
	currWoker.InitIdWorker(1000, 1)
	newID, err := currWoker.NextId()
	if err != nil {
		log.Println("InsertOrderByMessage生成ID错误： ", err)
		return err
	}

	var order = model.Order{}
	order.Uid = orderMessage.Uid
	order.Iid = orderMessage.Iid
	order.Oid = strconv.Itoa(int(newID))
	order.GoodsImage = orderMessage.Imageurl
	order.GoodsPrice = orderMessage.Price
	order.GoodsTitle = orderMessage.Title
	order.GoodsDesc = orderMessage.Desc
	order.Count = orderMessage.Count
	order.CreatedAt = orderMessage.CreatedAt

	err = data.Db.Create(&order).Error
	if err != nil {
		log.Panic("Insert InsertOrderByMessage failed: ", err.Error())
		return err
	}
	return nil
}

func SubStockByIid(iid string) (err error) {
	if err := data.Db.Model(&model.Goods{}).Where("iid = ? ", iid).Update("stock", gorm.Expr("stock- ?", 1)).Error; err != nil {
		log.Panic("UPDATE goods err : ", err)
		return err
	}
	if err := data.Db.Model(&model.GoodsSecKill{}).Where("iid = ? ", iid).Update("stock", gorm.Expr("stock- ?", 1)).Error; err != nil {
		log.Panic("UPDATE GoodsSecKill err : ", err)
		return err
	}

	return nil
}

func GetOrdersByUid(uid string, page int, pageSize int) []model.Order {
	var orderList [] model.Order
	err :=  data.Db.Where("uid = ?", uid).Order("created_at desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&orderList).Error
	if err != nil {
		log.Panic("GetOrdersByIid err: ", err)
	}
	return orderList
}