package reqStruct

// 该结构体可以通过json串自动生成，先编写好接口文档，声明好要传的json串，然后自动生成即可
// 自动生成工具：https://mholt.github.io/json-to-go/
type GoodsJson struct {
	Goodsimage   Goodsimage   `json:"goodsImage"`
	Goodsshop    Goodsshop    `json:"goodsShop"`
	Goodscomment Goodscomment `json:"goodsComment"`
	Goods        Goods        `json:"goods"`
}
type Goodsimage struct {
	Topimages   []string `json:"topImages"`
	Detailimage []string `json:"detailImage"`
}
type Score struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type Goodsshop struct {
	Score    []Score `json:"score"`
	Cfans    int64   `json:"cFans"`
	Csells   int64   `json:"cSells"`
	Shoplogo string  `json:"shopLogo"`
	Name     string  `json:"name"`
	Cgoods   int64   `json:"cGoods"`
	Shopid   string  `json:"shopId"`
	Userid   string  `json:"userId"`
}
type User struct {
	Uname  string `json:"uname"`
	Avatar string `json:"avatar"`
}
type List struct {
	User    User   `json:"user"`
	Content string `json:"content"`
	Created int64  `json:"created"`
	Style   string `json:"style"`
}
type Goodscomment struct {
	List []List `json:"list"`
}
type Goods struct {
	Iid          string `json:"iid"`
	Shopid       string `json:"shopId"`
	Cfav         int64  `json:"cfav"`
	Sells        int64  `json:"sells"`
	Image        string `json:"image"`
	Desc         string `json:"desc"`
	Title        string `json:"title"`
	Discountdesc string `json:"discountDesc"`
	Oldprice     string `json:"oldPrice"`
	Price        string `json:"price"`
	Lownowprice  string `json:"lowNowPrice"`
	Stock        int64  `json:"stock"`
	Delivery     string `json:"delivery"`
}
