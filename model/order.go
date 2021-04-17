package model

// 数据库订单实体
type Order struct {
	Oid        string `gorm:"type:varchar(60); not null"`  // 订单ID
	Iid        string `gorm:"type:varchar(60); not null"`  // 商品ID
	Uid        string `gorm:"type:varchar(60); not null"`  // 用户ID
	GoodsTitle string `gorm:"type:text; not null"`         // 商品标题
	GoodsDesc  string `gorm:"type:text; not null"`         // 商品描述
	GoodsImage string `gorm:"type:varchar(100); not null"` // 商品的概览图URL
	GoodsPrice string `gorm:"type:float(24); not null"`    // 商品单价
	Count      int64  `gorm:"type:int; not null"`          // 购买的数量
	CreatedAt  int64  `gorm:"type:int; not null"`          // 创建时间
}

