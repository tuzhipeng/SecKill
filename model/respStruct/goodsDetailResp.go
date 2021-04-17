package respStruct

// 该结构体可以通过json串自动生成，先编写好接口文档，声明好要传的json串，然后自动生成即可
// 自动生成工具：https://mholt.github.io/json-to-go/
type GoodsDetailRespJson struct {
	Result  Result `json:"result"`
	Success bool   `json:"success"`
}
type Iteminfo struct {
	Topimages    []string `json:"topImages"`
	Desc         string   `json:"desc"`
	Title        string   `json:"title"`
	Discountdesc string   `json:"discountDesc"`
	Oldprice     string   `json:"oldPrice"`
	Price        string   `json:"price"`
	Lownowprice  string   `json:"lowNowPrice"`
	Iid          string   `json:"iid"`
}
type Score struct {
	Name     string  `json:"name"`
	Score    float64 `json:"score"`
	Isbetter bool    `json:"isBetter"`
}
type Shopinfo struct {
	Score    []Score `json:"score"`
	Cfans    int64   `json:"cFans"`
	Csells   int64   `json:"cSells"`
	Shoplogo string  `json:"shopLogo"`
	Name     string  `json:"name"`
	Cgoods   int64   `json:"cGoods"`
	Shopid   string  `json:"shopId"`
}
type User struct {
	Uname  string `json:"uname"`
	Avatar string `json:"avatar"`
}
type CommentItem struct {
	User    User   `json:"user"`
	Content string `json:"content"`
	Created int64  `json:"created"`
	Style   string `json:"style"`
}
type Rate struct {
	List []CommentItem `json:"list"`
}
type Detailinfo struct {
	Desc        string        `json:"desc"`
	Detailimage []Detailimage `json:"detailImage"`
}

type Detailimage struct {
	Desc   string   `json:"desc"`
	Key    string   `json:"key"`
	Anchor string   `json:"anchor"`
	List   []string `json:"list"`
}
type Result struct {
	Iteminfo   Iteminfo   `json:"itemInfo"`
	Columns    []string   `json:"columns"`
	Shopinfo   Shopinfo   `json:"shopInfo"`
	Rate       Rate       `json:"rate"`
	Detailinfo Detailinfo `json:"detailInfo"`
}
