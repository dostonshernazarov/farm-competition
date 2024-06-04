package models

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type DeliveryReq struct {
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Capacity    int64  `json:"capacity"`
	Union       string `json:"union"`
	Time        string `json:"time" example:"2024-01-01"`
}

type DeliveryRes struct {
	ID          string `json:"id"`
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Capacity    int64  `json:"capacity"`
	Union       string `json:"union"`
	Time        string `json:"time" example:"2024-01-01"`
}

type ListDeliverysRes struct {
	Delivery []*DeliveryRes `json:"deliveries"`
}

func (t *DeliveryReq) Validate() error {
	t.ProductName = strings.ToLower(t.ProductName)
	t.Category = strings.ToLower(t.Category)
	t.Union = strings.ToLower(t.Union)
	t.Time = strings.ToLower(t.Time)
	return validation.ValidateStruct(t,
		validation.Field(
			&t.ProductName,
			validation.Required,
		),
		validation.Field(
			&t.Category,
			validation.Required,
		),
		validation.Field(
			&t.Union,
			validation.Required,
		),

		validation.Field(
			&t.Time,
			validation.Required,
			validation.Date(time.ANSIC),
		),
	)
}
