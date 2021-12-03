package model

type Goods struct {
	OrderID int64 `json:"order_id"`
}

type CreatedGoodsMsg struct {
	Data Goods `json:"data"`
}

type RejectedGoodsMsg struct {
	Data Goods `json:"data"`
}
