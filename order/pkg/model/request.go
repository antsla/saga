package model

type OrderData struct {
	UserID   int64   `json:"user_id"`
	GoodsIds []int64 `json:"goods_ids"`
}
