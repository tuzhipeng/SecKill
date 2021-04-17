package cache

import (
	"GraduateDesign/data"
	"fmt"
	"log"
)

const LuaScript  =`
local goods_key = KEYS[1]
local goods_stock = tonumber(redis.call('GET', goods_key))
if (goods_stock > 0) then
	redis.call('DECR', goods_key)
	return 1
end
return 0
`
// 确保redis加载lua脚本，若未加载则加载
func prepareScript(script string) string {
	// sha := sha1.Sum([]byte(script))
	scriptsExists, err :=data.Client.ScriptExists(script).Result()
	if err != nil {
		panic("Failed to check if script exists: " + err.Error())
	}
	if !scriptsExists[0] {
		scriptSHA, err := data.Client.ScriptLoad(script).Result()
		if err != nil {
			panic("Failed to load script " + script + " err: " + err.Error())
		}
		return scriptSHA
	}
	print("Script Exists.")
	return ""
}
func RemoteDecrementGoodsStock(goodsIid string)(int64, error)  {
	goodsKey := fmt.Sprintf(SecKillGoodsStockKey, goodsIid)
	luaScript := prepareScript(LuaScript)
	res, err := data.Client.EvalSha(luaScript, []string{goodsKey}).Result()
	if err != nil {
		log.Println("RemoteDecrementGoodsStock EvalSha luaScript err ", err)
		return -1, err
	}

	return res.(int64), err

	//fmt.Printf("res:%v,类型是%T, err: %v\n", res, res, err)
	//fmt.Println("res.typeof", reflect.TypeOf(res))0
}
