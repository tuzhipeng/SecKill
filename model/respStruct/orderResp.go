package respStruct

type OrderRespJson struct {
	List []OrderListItem `json:"list"`
}
type OrderListItem struct {
	Iid       string `json:"iid"`
	Oid       string `json:"oid"`
	Imageurl  string `json:"imageUrl"`
	Title     string `json:"title"`
	Createdat int64    `json:"createdAt"`
	Price     string `json:"price"`
	Count     int64    `json:"count"`
}