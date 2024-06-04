package models

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type DrugReq struct {
	DrugName      string `json:"drug_name"`
	Union         string `json:"union" example:"piece"`
	Description   string `json:"description"`
	TotalCapacity int64  `json:"total_capacity"`
	Status        string `json:"status"`
}

type DrugRes struct {
	Id            string `json:"id"`
	DrugName      string `json:"drug_name"`
	Union         string `json:"union"`
	Description   string `json:"description"`
	TotalCapacity int64  `json:"total_capacity"`
	Status        string `json:"status"`
}

type DrugFieldValues struct {
	Name   string `json:"name"`
	Union  string `json:"union"`
	Status string `json:"status"`
}

type ListDrugsRes struct {
	Drugs []*DrugRes `json:"drugs"`
	Count int64 `json:"count"`
}

func (t *DrugReq) Validate() error {
	t.DrugName = strings.ToLower(t.DrugName)
	t.Union = strings.ToLower(t.Union)
	t.Description = strings.ToLower(t.Description)
	t.Status = strings.ToLower(t.Status)
	return validation.ValidateStruct(t,
		validation.Field(
			&t.DrugName,
			validation.Required,
		),
		validation.Field(
			&t.Union,
			validation.Required,
		),
	)

}
