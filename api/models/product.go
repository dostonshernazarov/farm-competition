package models

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ProductReq struct {
	ProductName   string `json:"product_name"`
	Union         string `json:"union"`
	Description   string `json:"description"`
	TotalCapacity int64  `json:"total_capacity"`
}

type ProductRes struct {
	Id            string `json:"id"`
	ProductName   string `json:"product_name"`
	Union         string `json:"union"`
	Description   string `json:"description"`
	TotalCapacity int64  `json:"total_capacity"`
}

type ProductFieldValues struct {
	Name  string `json:"name"`
	Union string `json:"union"`
}

type ListProductsRes struct {
	Products []*ProductRes `json:"products"`
}

func (t *ProductReq) Validate() error {
	t.ProductName = strings.ToLower(t.ProductName)
	t.Union = strings.ToLower(t.Union)
	t.Description = strings.ToLower(t.Description)
	return validation.ValidateStruct(t,
		validation.Field(
			&t.ProductName,
			validation.Required,
		),
		validation.Field(
			&t.Union,
			validation.Required,
		),
	)

}
