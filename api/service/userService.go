package service

import (
	"GraduateDesign/data"
	"GraduateDesign/model"
	"GraduateDesign/model/reqStruct"
	"GraduateDesign/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"time"
)

const (
	SmsCodeSize        = 6
	UserTokenSize      = 64
	LoginSmsCodeKey    = "login:sms_code:%s"
	LoginSmsCodeExpire = 30
	UserTokenKey       = "user:token:%s"
	UserTokenExpire    = 3600 * 24
	GoodsStockKey      = "goods_stock:%d"
	OrderUidGidKey     = "order:%d:%d"

)

//var Ctx = context.Background()

func LoginService(postUser reqStruct.UserJson, ctx *gin.Context) (token string, err error) {

	// 查找该用户

	queryUser := model.User{}
	dbErr := data.Db.Where("uid = ?", postUser.Uid).Find(&queryUser).Error
	//dbErr := data.Db.Where(&queryUser).First(&queryUser).Error
	if dbErr != nil && gorm.IsRecordNotFoundError(dbErr) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errMsg": "No such queryUser."})
		return "", dbErr
	}

	// 匹配密码
	if queryUser.Password != postUser.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errMsg": "Password mismatched."})
		return "", errors.New("password mismatched")
	}
	//data.Db.First(queryUser)
	// 生成token并保存
	token = utils.CreateKey(utils.AlphabetAndNumber, UserTokenSize)
	if err := saveUserToken(token, queryUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errMsg": "generate token err."})
		return "", err
	}
	//fmt.Println("生成的token：", token)
	return token, nil
}

//func GetUserInfoService()  {
//
//}

// 将token存入redis，格式"user:token": "{user结构体}"
func saveUserToken(token string, user model.User) (err error) {
	var userData []byte
	if userData, err = json.Marshal(user); err != nil {
		log.Printf("json.Marshal() failed, err: %v", err)
		return
	}
	if err = data.Client.Set(fmt.Sprintf(UserTokenKey, token), string(userData), time.Duration(UserTokenExpire)*time.Second).Err(); err != nil {
		log.Printf("redis.Set() failed, err: %v", err)
	}
	return
}

func GetUserByToken(token string) (user model.User, err error) {
	var userData string
	if userData, err = data.Client.Get(fmt.Sprintf(UserTokenKey, token)).Result(); err != nil {
		if err == redis.Nil {
			err = nil
		} else {
			log.Printf("redis.Get() failed, err: %v", err)
			return
		}
	}
	if len(userData) == 0 {
		err = errors.New("token is empty")
		return
	}
	if err = json.Unmarshal([]byte(userData), &user); err != nil {
		log.Printf("json.Unmarshal() failed, err: %v", err)
	}
	return
}

func RenewUserToken(token string) (err error) {
	if err = data.Client.Expire(fmt.Sprintf(UserTokenKey, token), time.Duration(UserTokenExpire)*time.Second).Err(); err != nil {
		log.Printf("redis.Expire() failed, err: %v", err)
	}
	return
}
