package cache

import (
	"fmt"
	"testing"
)

//func TestMain(m *testing.M) {
//	m.Run()
//
//}
//func TestReduceStockWorkFlow(t *testing.T)  {
//	t.Run("测试远程扣除秒杀商品的库存", testRemoteDecrementGoodsStock)
//
//}
func testRemoteDecrementGoodsStock(t *testing.T) {
	res, err := RemoteDecrementGoodsStock("secKillGoodsStock:iid:1m7k5lw")
	fmt.Println(res, err)
}