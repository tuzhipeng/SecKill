package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const configFilePath string = "config-prod.yaml"
//const configFilePath string = "config-dev.yaml"
//const configFilePath string = "config-test.yaml"

var configFile []byte

// 根据yaml文件中的配置，定义对应结构体
type AppConfig struct {
	App App `yaml:"app"`
}

// 根据yaml文件中的配置，定义对应结构体
type App struct {
	Database        Database `yaml:"database"`
	Redis           Redis    `yaml:"redis"`
	RabbitMQ        RabbitMQ `yaml:"rabbitmq"`
	FlushAllForTest bool     `yaml:"flushAllForTest"`
}

type Database struct {
	Type     string `yaml:"type"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbName"`
	Address  string `yaml:"address"`
	MaxIdle  int    `yaml:"maxIdle"`
	MaxOpen  int    `yaml:"maxOpen"`
}

type Redis struct {
	Address     string `yaml:"address"`
	Network     string `yaml:"network"`
	Password    string `yaml:"password"`
	MaxIdle     int    `yaml:"maxIdle"`
	MaxActive   int    `yaml:"maxActive"`
	IdleTimeout int    `yaml:"idleTimeout"`
}

type RabbitMQ struct {
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Address     string `yaml:"address"`
	VirtualHost string `yaml:"virtualHost"`
	QueenName   string `yaml:"queenName"`
}

// 根据预先定义好的结构体，将yaml文件中内容unmarshal出来
func GetAppConfig() (appConfig AppConfig, err error) {
	err = yaml.Unmarshal(configFile, &appConfig)
	return appConfig, err
}

// 通过io流将文件写入预先定义好的字节流configFile中
func init() {
	var err error
	configFile, err = ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("yamlFile get err %v", err)
	}

}
