package data

import (
	"GraduateDesign/conf"
	"GraduateDesign/model/reqStruct"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

//rabbitMQ结构体
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl string
	sync.Mutex
}

var Rabbitmq *RabbitMQ
//var MQRUL string

func initRabbitMQ(config conf.AppConfig)  {
	username := config.App.RabbitMQ.Username
	password := config.App.RabbitMQ.Password
	address := config.App.RabbitMQ.Address
	virtualHost := config.App.RabbitMQ.VirtualHost
	queenName := config.App.RabbitMQ.QueenName

	MQURL := fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, address,virtualHost)

	Rabbitmq = &RabbitMQ{QueueName: queenName, Exchange:"", Key:"", Mqurl:MQURL}
	var err error
	//获取connection
	Rabbitmq.Conn, err = amqp.Dial(Rabbitmq.Mqurl)
	Rabbitmq.failOnErr(err, "failed to connect Rabbitmq!")
	//获取channel
	Rabbitmq.Channel, err = Rabbitmq.Conn.Channel()
	Rabbitmq.failOnErr(err, "failed to open a channel")

}

//断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.Channel.Close()
	r.Conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}


//直接模式队列生产
func (r *RabbitMQ) PublishSimple(message string){
	//r.Lock()
	//defer r.Unlock()
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.Channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil)
	if err != nil {
		fmt.Println("r.channel.QueueDeclare err:", err)
		//return err
	}
	//调用channel 发送消息到队列中
	r.Channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	fmt.Println("publish 完成")
	//return nil
}
/*func (r *RabbitMQ) PublishSimple (message string)  {
	// 1. 申请队列
	_, err := r.channel.QueueDeclare(r.QueueName, false, false,
		false,false, nil)
	if err != nil {
		fmt.Println(err)
	}

	r.channel.Publish(r.Exchange, r.QueueName, false,false,
		amqp.Publishing{
			ContentType:"text/plain",
			Body: []byte(message),
		})
}*/
//simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.Channel.QueueDeclare(r.QueueName,false,false,false,false,nil)
	if err != nil {
		fmt.Println(err)
	}
	//消费者流控
	r.Channel.Qos(
		1, //当前消费者一次能接受的最大消息数量
		0, //服务器传递的最大容量（以八位字节为单位）
		false, //如果设置为true 对channel可用
	)

	//接收消息
	msgs, err := r.Channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"", // consumer
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
			//err = dao.InsertOrderByMessage(message)
			//if err !=nil {
			//	fmt.Println(err)
			//}
			//
			////扣除商品数量
			//err = dao.SubStockByIid(message.Iid)
			//if err !=nil {
			//	fmt.Println(err)
			//}
			//如果为true表示确认所有未确认的消息，
			//为false表示确认当前消息
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

/*
// 连接信息
const MQURL  = "amqp://admin:admin@47.93.12.71:5672/testGraduateDesign"

// rabbitMQ 结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	// 队列名称
	QueueName string
	// 交换机名称
	Exchange string
	// bind key 4名称
	Key string
	// 连接信息
	Mqurl string
}

// 断开channel和coon
func (r *RabbitMQ) Destroy()  {
	r.channel.Close()
	r.conn.Close()
}
// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string)  {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}

}
//  创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName:queueName, Exchange:exchange, Key:key, Mqurl:MQURL}

}

//  创建简单模式下的RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitMQSimple := NewRabbitMQ(queueName, "", "")

	var err error
	// 获取connection
	rabbitMQSimple.conn, err = amqp.Dial(rabbitMQSimple.Mqurl)
	rabbitMQSimple.failOnErr(err, "Failed to connect the rabbitMQ")

	// 获取channel
	rabbitMQSimple.channel, err = rabbitMQSimple.conn.Channel()
	rabbitMQSimple.failOnErr(err, "Failed to open a channel")

	return rabbitMQSimple

}

//  Simple模式生产方法绑定队列
func (r *RabbitMQ) PublishSimple (message string)  {
	// 1. 申请队列
	_, err := r.channel.QueueDeclare(r.QueueName, false, false,
		false,false, nil)
	if err != nil {
		fmt.Println(err)
	}

	r.channel.Publish(r.Exchange, r.QueueName, false,false,
		amqp.Publishing{
			ContentType:"text/plain",
			Body: []byte(message),
		})
}

//  simple模式消费方法
func (r *RabbitMQ) ConsumeSimple()  {
	q, err := r.channel.QueueDeclare(r.QueueName, false, false,
		false,false, nil)
	if err != nil {
		fmt.Println(err)
	}
	// 接受消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"",     // consumer
		//是否自动应答
		true,   // auto-ack
		//是否独有
		false,  // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false,  // no-local
		//列是否阻塞
		false,  // no-wait
		nil,    // args)
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	// 启动协程处理消息
	go func() {
		for d := range msgs{
			log.Printf("Received a message: %s ", d.Body)
		}
	}()

	log.Printf("Waiting for messages...")
	<- forever
}
*/