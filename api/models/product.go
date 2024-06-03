package models

type ProductReq struct {
	ProductName string `json:"product_name"`
	Union string `json:"union"`
	Description string `json:"description"`
	TotalCapacity int64 `json:"total_capacity"`
}

type ProductRes struct {
	Id string `json:"id"`
	ProductName string `json:"product_name"`
	Union string `json:"union"`
	Description string `json:"description"`
	TotalCapacity int64 `json:"total_capacity"`
}

type ProductFieldValues struct {
	Name string `json:"name"`
	Union string `json:"union"`
}

type ListProductsRes struct {
	Products []*ProductRes `json:"products"`
}