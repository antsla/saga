package model

type Order struct {
	ID       int64   `json:"id"`
	GoodsIds []int64 `json:"goods_ids"`
}

type CreatedOrderMsg struct {
	Data Order `json:"data"`
}
