package models

import (
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type DeliveryCreateReq struct {
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Capacity    int64  `json:"capacity"`
	Union       string `json:"union"`
	Time        string `json:"time" example:"2024-01-01"`
	Status string `json:"status"`
	Description string `json:"description"`
}

type DeliveryReq struct {
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Capacity    int64  `json:"capacity"`
	Union       string `json:"union"`
	Time        string `json:"time" example:"2024-01-01"`
}

type DeliveryCreateRes struct {
	Message          string `json:"message"`
}

type DeliveryRes struct {
	ID string `json:"id"`
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Capacity    int64  `json:"capacity"`
	Union       string `json:"union"`
	Time        string `json:"time" example:"2024-01-01 12:00:00"`
}

type ListDeliverysRes struct {
	Delivery []*DeliveryRes `json:"deliveries"`
}

type DeliveryFieldValues struct {
	Name string `json:"name"`
	Category string `json:"category"`
	Time string `json:"time"`
}

func (t *DeliveryCreateReq) Validate() error {
	t.ProductName = strings.ToLower(t.ProductName)
	t.Category = strings.ToLower(t.Category)
	t.Union = strings.ToLower(t.Union)
	t.Time = strings.ToLower(t.Time)

	if t.Category == "drug" && t.Status == "" {
		return errors.New("status required in drug")
	}
	return validation.ValidateStruct(t,
		validation.Field(
			&t.ProductName,
			validation.Required,
		),
		validation.Field(
			&t.Category,
			validation.Required,
			validation.In("food", "drug"),
		),
		validation.Field(
			&t.Union,
			validation.Required,
		),

		validation.Field(
			&t.Time,
			validation.Required,
			validation.Date(time.DateTime),
		),
	)
}
