package data

import (
	"GraduateDesign/conf"
	"GraduateDesign/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

var Db *gorm.DB

// 初始化连接，确保创建表、索引等
func initMysql(config conf.AppConfig) {
	fmt.Println("加载数据库配置中...")

	// 设置连接的参数
	dbType := config.App.Database.Type

	usr := config.App.Database.User
	pwd := config.App.Database.Password
	address := config.App.Database.Address
	dbName := config.App.Database.DbName
	dbLink := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		usr, pwd, address, dbName)

	//创建一个数据库的连接，因为docker中的mysql服务启动时延，一开始需要尝试重试连接
	log.Println("正在初始化连接，准备连接...")
	var err error
	for Db, err = gorm.Open(dbType, dbLink); err != nil; Db, err = gorm.Open(dbType, dbLink) {
		log.Println("无法连接数据库，错误: ", err.Error())
		log.Println("5秒后重连...")
		time.Sleep(5 * time.Second)
	}

	// 设置连接池连接数
	Db.DB().SetMaxIdleConns(config.App.Database.MaxIdle)
	Db.DB().SetMaxOpenConns(config.App.Database.MaxOpen)
	Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8")

	// 准备结构体映射
	user := model.User{}
	goods := model.Goods{}
	goodsComment := model.GoodsComment{}
	goodsImage := model.GoodsImage{}
	goodsShop := model.GoodsShop{}
	goodsSecKill := model.GoodsSecKill{}
	order := model.Order{}

	// 创建表
	tables := []interface{}{user, goods, goodsComment,
		goodsImage, goodsShop, goodsSecKill, order}

	for _, table := range tables {
		if !Db.HasTable(table) {
			Db.AutoMigrate(table)
		} else if config.App.FlushAllForTest {
			log.Println("FlushAllForTest为true，清空所有的表做测试...")
			Db.Delete(table)
			//Db.DropTable(table)
			//Db.AutoMigrate(table)
		}
	}

	// TODO：创建索引
	//"iid, cfav, price, title, image")

}
