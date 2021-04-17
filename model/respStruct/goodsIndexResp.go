package respStruct

// 该结构体可以通过json串自动生成，先编写好接口文档，声明好要传的json串，然后自动生成即可
// 自动生成工具：https://mholt.github.io/json-to-go/
type GoodsIndexRespJson struct {
	Data    Data `json:"data"`
	Success bool `json:"success"`
}
type List struct {
	Image string `json:"image"`
	Title string `json:"title"`
	Price string `json:"price"`
	Cfav  int    `json:"cfav"`
	Iid   string `json:"iid"`
}
type Data struct {
	List []List `json:"list"`
}
