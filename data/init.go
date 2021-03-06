package data

import (
	"GraduateDesign/conf"
	"log"
)

func init() {
	config, err := conf.GetAppConfig()
	if err != nil {
		log.Panic("failed to load data config: " + err.Error())
	}

	initMysql(config)
	initRedisConnection(config)
	initRabbitMQ(config)
}

func Close() {
	err := Client.Close()
	if err != nil {
		log.Println("Error on closing redisService Client.")
	}
	err = Db.Close()
	if err != nil {
		log.Println("Error on closing dbService Client.")
	}
}
