package models

type DeliveryReq struct {
	ProductName string `json:"product_name"`
	Category string `json:"category"`
	Capacity int64 `json:"capacity"`
	Union string `json:"union"`
	Time string `json:"time" example:"2024-01-01"`
}


type DeliveryRes struct {
	ID string `json:"id"`
	ProductName string `json:"product_name"`
	Category string `json:"category"`
	Capacity int64 `json:"capacity"`
	Union string `json:"union"`
	Time string `json:"time" example:"2024-01-01"`
}

type ListDeliverysRes struct {
	Delivery []*DeliveryRes `json:"deliveries"`
}