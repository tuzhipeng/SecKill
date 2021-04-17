package model

// 数据库商品实体
type Goods struct {
	Iid          string `gorm:"type:varchar(60); not null"`  // 商品itemID
	Sid          string `gorm:"type:varchar(60); not null"`  // 商品所属店铺ID
	Image        string `gorm:"type:varchar(100); not null"` // 概览图片imageURL
	Title        string `gorm:"type:text; not null"`         // 商品标题
	Cfav         int64  `gorm:"type:int; not null"`          // 收藏数
	OldPrice     string `gorm:"type:varchar(60); not null"`  // 原价
	Price        string `gorm:"type:varchar(60); not null"`  // 现价
	LowNowPrice  string `gorm:"type:varchar(60); not null"`  // 折扣价
	Desc         string `gorm:"type:text; not null"`         // 详细描述
	DiscountDesc string `gorm:"type:varchar(60); not null"`  // 折扣描述 “优惠价”
	Stock        int64  `gorm:"type:int; not null"`          // 库存 默认10000
	Sells        int64  `gorm:"type:int; not null"`          // 销量
	Delivery     string `gorm:"type:varchar(100); not null"` // 默认快递
}

// 商品展示图片表
type GoodsImage struct {
	Iid   string `gorm:"type:varchar(60); not null"`  // 商品itemID
	Image string `gorm:"type:varchar(100); not null"` // 商品图片URL
	Ikind int64  `gorm:"type:int; not null"`          // 区分滚动图还是竖栏图片
}

// 店铺表
type GoodsShop struct {
	Uid          string  `gorm:"type:varchar(60); not null"`  // 商家ID
	Sid          string  `gorm:"type:varchar(60); not null"`  // 店铺ID
	ShopLogo     string  `gorm:"type:varchar(100); not null"` // 店铺logoURL
	Name         string  `gorm:"type:varchar(100); not null"` // 店铺名称
	CFans        int64   `gorm:"type:int; not null"`          // 店铺粉丝数
	CSells       int64   `gorm:"type:int; not null"`          // 总销量
	CGoods       int64   `gorm:"type:int; not null"`          // 总商品数
	QualityScore float64 `gorm:"type:float(24); not null"`    // 质量满意评分
	PriceScore   float64 `gorm:"type:float(24); not null"`    // 价格合理评分
	DescScore    float64 `gorm:"type:float(24); not null"`    // 相关描述评分

}

// 商品评价表
type GoodsComment struct {
	Iid     string `gorm:"type:varchar(60); not null"`  // 商品ID
	Uname   string `gorm:"type:varchar(60); not null"`  // 用户名
	Avatar  string `gorm:"type:varchar(100); not null"` // 用户名头像URL
	Content string `gorm:"type:text; not null"`         // 评语
	Style   string `gorm:"type:varchar(60); not null"`  // 衣服款式
	Created int64  `gorm:"type:int; not null"`          // 创建时间
}

// 秒杀商品表
type GoodsSecKill struct {
	Iid       string `gorm:"type:varchar(60); not null"` // 商品ID
	Stock     int64  `gorm:"type:int; not null"`         // 秒杀商品剩余库存
	StartTime int64  `gorm:"type:int; not null"`         // 秒杀开始时间
	EndTime   int64  `gorm:"type:int; not null"`         // 秒杀结束时间
}
