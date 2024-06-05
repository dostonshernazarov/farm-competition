package models

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AnimaEatablesInfoReq struct {
	AnimalID   string   `json:"animal_id"`
	EatablesID string   `json:"eatables_id"`
	Daily      []*Daily `json:"daily"`
	Category   string   `json:"category"`
}

type Daily struct {
	Time     string `json:"time"`
	Capacity int64  `json:"capacity"`
}

type AnimalDrugInfoRes struct {
	ID       string  `json:"id"`
	AnimalID string  `json:"animal_id"`
	Eatables DrugRes `json:"eatables"`
	Daily    []Daily `json:"daily"`
	Category string  `json:"category"`
}

type AnimaEatablesInfoRes struct {
	ID         string   `json:"id"`
	AnimalID   string   `json:"animal_id"`
	EatablesID string   `json:"eatables_id"`
	Daily      []*Daily `json:"daily"`
	Category   string   `json:"category"`
}

type AnimaFoodInfoRes struct {
	ID       string  `json:"id"`
	AnimalID string  `json:"animal_id"`
	Eatables FoodRes `json:"eatables"`
	Daily    []Daily `json:"daily"`
	Category string  `json:"category"`
}

type ListDrugInfoByAnimalRes struct {
	Eatables []*AnimalDrugInfoRes `json:"eatables"`
	Count    int64                `json:"count"`
}

type ListFootInfoByAnimalRes struct {
	Eatables []*AnimaFoodInfoRes `json:"eatables"`
	Count    int64               `json:"count"`
}

type ListEatablesInfoByAnimalReq struct {
	AnimalID string `json:"animal_id"`
}

func (t *AnimaEatablesInfoReq) Validate() error {
	t.Category = strings.ToLower(t.Category)

	return validation.ValidateStruct(t,
		validation.Field(
			&t.AnimalID,
			validation.Required,
		),
		validation.Field(
			&t.Category,
			validation.Required,
			validation.In("food", "drug"),
		),
		validation.Field(
			&t.EatablesID,
			validation.Required,
		),
	)
}
