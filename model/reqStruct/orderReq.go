package reqStruct

type OrderJson struct {
	List []OrderListItem `json:"list"`
}
type OrderListItem struct {
	Iid      string `json:"iid"`
	Imageurl string `json:"imageUrl"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Price    string `json:"price"`
	Count    int64  `json:"count"`
}

type OrderMessage struct {
	Uid string
	OrderListItem

	CreatedAt int64
}
