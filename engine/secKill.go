package engine

import (
	"GraduateDesign/api/controller"
	"GraduateDesign/conf"
	"GraduateDesign/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

//const SessionHeaderKey = "Authorization"

func SecKillEngine() *gin.Engine {
	router := gin.New()
	router.Use(middleware.Cors())

	config, err := conf.GetAppConfig()
	if err != nil {
		log.Panic("无法加载config, 错误：" + err.Error())
	}
	log.Println(config)

	//store, _ := redis.NewStore(config.App.Redis.MaxIdle, config.App.Redis.Network,
	//	config.App.Redis.Address, config.App.Redis.Password, []byte("seckill"))
	//router.Use(sessions.Sessions(SessionHeaderKey, store))
	//gob.Register(&model.User{})

	// 插入商品
	goodsRouter := router.Group("/api/goods")
	{
		goodsRouter.POST("", controller.PubGoods)
		goodsRouter.GET("", controller.GetGoodsDetail)
		goodsRouter.GET("/recommend", controller.GetGoodsRecommend)
	}
	//获取首页商品列表
	homeRouter := router.Group("/api/home")
	{
		homeRouter.GET("/data", controller.GetGoodsIndexList)
		homeRouter.GET("/multidata", controller.GetHomeMultiData)
	}
	// 生成订单
	//orderRouter := router.Group("/api/order")
	//{
	//	orderRouter.POST("", controller.GenerateOrder)
	//}
	// 用户登录
	userRouter := router.Group("api/user")
	{
		userRouter.POST("", controller.Login)
		//userRouter.GET("", controller.GetUserInfo)
	}

	authRouter := router.Group("").Use(middleware.Auth())
	//authRouter.Use(middleware.Auth())
	{
		authRouter.GET("api/user", controller.GetUserInfo)
		authRouter.POST("/api/order", controller.GenerateOrder)
		authRouter.GET("/api/order", controller.GetOrder)
	}

	// 测试路由
	/*testRouter := router.Group("/api/test")
	//testRouter.Use(middleware.Auth())
	{
		testRouter.GET("", func(c *gin.Context) {
			//spider.InitGoods()
			c.JSON(200, gin.H{
				"hello": "world",
			})

		})
		testRouter.GET("/flush", func(context *gin.Context) {
			if _, err := data.FlushAll(); err != nil {
				fmt.Println("Error when flushAll. " + err.Error())
			} else {
				fmt.Println("Flushall succeed.")
			}
		})
	}*/
	return router

}
