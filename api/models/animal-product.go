package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AnimalProductReq struct {
	AnimalID string `json:"animal_id"`
	ProductID string `json:"product_id"`
	Capacity int64 `json:"capacity"`
	GetTime string `json:"get_time"`
}

type AnimalProductUpdateReq struct{
	ID string `json:"id"`
	AnimalID string `json:"animal_id"`
	ProductID string `json:"product_id"`
	Capacity int64 `json:"capacity"`
	GetTime string `json:"get_time"`
}

type AnimalProductRes struct {
	Id string `json:"id"`
	AnimalID string `json:"animal_id"`
	AnimalName string `json:"animal_name"`
	AnimalCategory string `json:"animal_category"`
	ProductName string `json:"product_name"`
	Capacity int64 `json:"capacity"`
	Union string `json:"union"`
	GetTime string `json:"get_time"`
}

type AnimalProductFieldValues struct {
	GetTime string `json:"get_time"`
}

type AnimalProductByAnimalIdFieldValues struct {
	AnimalID string `json:"animal_id"`
}

type AnimalProductByProductIdFieldValues struct {
	ProductID string `json:"product_id"`
}

type ListAnimalProductsRes struct {
	AnimalProducts []*AnimalProductRes `json:"animal_products"`
	Count int64 `json:"count"`
}

type AnimalProductByAnimalIdRes struct {
	Animal   *AnimalRes `json:"animal"`
	Products []*ProductRes  `json:"products"`
	Count int64 `json:"count"`

}

type AnimalProductByProductIdRes struct {
	Product ProductRes `json:"product"`
	Animals []*AnimalRes `json:"animals"`
	Count int64 `json:"count"`
}


func (t *AnimalProductReq) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(
			&t.AnimalID,
			validation.Required,
		),
		validation.Field(
			&t.GetTime,
			validation.Required,
			validation.Date(time.DateTime),
		),
		validation.Field(
			&t.ProductID,
			validation.Required,
		),
	)
}