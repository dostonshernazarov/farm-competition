package models

type AnimalProductReq struct {
	AnimalID string `json:"animal_id"`
	ProductID string `json:"product_id"`
	Capacity int64 `json:"capacity"`
	GetTime string `json:"get_time"`
}

type AnimalProductRes struct {
	Id string `json:"id"`
	AnimalID string `json:"animal_id"`
	ProductName string `json:"product_name"`
	Capacity int64 `json:"capacity"`
	GetTime string `json:"get_time"`
}

type AnimalProductFieldValues struct {
	AnimalID string `json:"animal_id"`
}

type ListAnimalProductsRes struct {
	AnimalProducts []*AnimalProductRes `json:"animal_products"`
}