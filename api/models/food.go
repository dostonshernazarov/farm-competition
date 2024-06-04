package models

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type FoodReq struct {
	FoodName      string `json:"food_name"`
	Union         string `json:"union" example:"piece"`
	Description   string `json:"description"`
	TotalCapacity int64  `json:"total_capacity"`
}

type FoodRes struct {
	Id            string `json:"id"`
	FoodName      string `json:"food_name"`
	Union         string `json:"union"`
	Description   string `json:"description"`
	TotalCapacity int64  `json:"total_capacity"`
}

type FoodFieldValues struct {
	Name  string `json:"name"`
	Union string `json:"union"`
}

type ListFoodsRes struct {
	Foods []*FoodRes `json:"foods"`
}

func (t *FoodReq) Validate() error {
	t.FoodName = strings.ToLower(t.FoodName)
	t.Union = strings.ToLower(t.Union)
	t.Description = strings.ToLower(t.Description)
	return validation.ValidateStruct(t,
		validation.Field(
			&t.FoodName,
			validation.Required,
		),
		validation.Field(
			&t.Union,
			validation.Required,
		),
	)

}
