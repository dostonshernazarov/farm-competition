package models

type FoodReq struct {
	FoodName string `json:"food_name"`
	Union string `json:"union" example:"piece"`
	Description string `json:"description"`
	TotalCapacity int64 `json:"total_capacity"`
}

type FoodRes struct {
	Id string `json:"id"`
	FoodName string `json:"food_name"`
	Union string `json:"union"`
	Description string `json:"description"`
	TotalCapacity int64 `json:"total_capacity"`
}

type FoodFieldValues struct {
	Name string `json:"name"`
	Union string `json:"union"`
}

type ListFoodsRes struct {
	Foods []*FoodRes `json:"foods"`
}