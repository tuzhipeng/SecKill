package main

import (
	"GraduateDesign/api/dao"
	"GraduateDesign/data"
	"GraduateDesign/model/reqStruct"
	"encoding/json"
	"fmt"
	"log"
)

func main()  {

	q, err :=data.Rabbitmq.Channel.QueueDeclare(data.Rabbitmq.QueueName,false,false,false,false,nil)
	if err != nil {
		fmt.Println(err)
	}
	//消费者流控
	/*data.Rabbitmq.Channel.Qos(
		1, //当前消费者一次能接受的最大消息数量
		0, //服务器传递的最大容量（以八位字节为单位）
		false, //如果设置为true 对channel可用
	)*/

	//接收消息
	msgs, err := data.Rabbitmq.Channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"secKill", // consumer
		//是否自动应答
		//这里要改掉，我们用手动应答
		false, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			log.Printf("Received a message: %s", d.Body)
			message := &reqStruct.OrderMessage{}
			err :=json.Unmarshal([]byte(d.Body),message)
			if err !=nil {
				fmt.Println(err)
			}
			//插入订单
			err = dao.InsertOrderByMessage(message)
			if err !=nil {
				fmt.Println(err)
			}

			//扣除商品数量
			err = dao.SubStockByIid(message.Iid)
			if err !=nil {
				fmt.Println(err)
			}
			//如果为true表示确认所有未确认的消息，
			//为false表示确认当前消息
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
