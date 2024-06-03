package models

type DrugReq struct {
	DrugName string `json:"drug_name"`
	Union string `json:"union" example:"piece"`
	Description string `json:"description"`
	TotalCapacity int64 `json:"total_capacity"`
	Status string `json:"status"`
}

type DrugRes struct {
	Id string `json:"id"`
	DrugName string `json:"drug_name"`
	Union string `json:"union"`
	Description string `json:"description"`
	TotalCapacity int64 `json:"total_capacity"`
	Status string `json:"status"`
}

type DrugFieldValues struct {
	Name string `json:"name"`
	Union string `json:"union"`
	Status string `json:"status"`
}

type ListDrugsRes struct {
	Drugs []*DrugRes `json:"drugs"`
}