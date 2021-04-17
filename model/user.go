package model

// 数据库用户实体
type User struct {
	Uid      string `gorm:"type: varchar(60); not null"`  // 用户ID
	Username string `gorm:"type: varchar(60); not null"`  // 用户名
	Password string `gorm:"type: varchar(100); not null"` // 用户密码
}
