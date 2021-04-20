package service

import (
	"GraduateDesign/api/dao"
	"GraduateDesign/data"
	"GraduateDesign/model/reqStruct"
	"GraduateDesign/model/respStruct"
	"encoding/json"
	"log"
	"time"
)

func GenerateOrderService(uid string, orderList []reqStruct.OrderListItem) {
	var orderMessage = reqStruct.OrderMessage{}
	for _, tempOrderItem := range orderList {
		orderMessage.Uid = uid
		orderMessage.OrderListItem = tempOrderItem
		orderMessage.CreatedAt = time.Now().Unix()

		orderMessageBytes, err := json.Marshal(orderMessage)
		if err != nil {
			log.Panic("GenerateOrderService err: ", err)
		}
		data.Rabbitmq.PublishSimple(string(orderMessageBytes))
	}

}

func GetOrderService(uid string, page int, pageSize int) respStruct.OrderRespJson {
	var orderResp respStruct.OrderRespJson
	var tempOrderRespItem respStruct.OrderListItem
	orderList := dao.GetOrdersByUid(uid, page, pageSize)
	for _, orderListItem := range orderList{
		tempOrderRespItem.Iid = orderListItem.Iid
		tempOrderRespItem.Oid = orderListItem.Oid
		tempOrderRespItem.Count = orderListItem.Count
		tempOrderRespItem.Title = orderListItem.GoodsTitle
		tempOrderRespItem.Imageurl = orderListItem.GoodsImage
		tempOrderRespItem.Price = orderListItem.GoodsPrice
		tempOrderRespItem.Createdat = orderListItem.CreatedAt
		orderResp.List = append(orderResp.List, tempOrderRespItem)
	}
	return orderResp
}
