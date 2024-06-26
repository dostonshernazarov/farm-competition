package models

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AnimalReq struct {
	Name    string `json:"name"`
	CategoryName       string `json:"category_name"`
	Gender    string `json:"gender" example:"male"`
	DateOfBirth string `json:"date_of_birth" example:"2024-01-01"`
	Description        string `json:"description"`
	Genus      string `json:"genus"`
	Weight float32 `json:"weight"`
	IsHealth bool `json:"is_health" example:"true"`
}

type AnimalRes struct {
	Id           string `json:"id"`
	Name    string `json:"name"`
	CategoryName       string `json:"category_name"`
	Gender    string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Description        string `json:"description"`
	Genus      string `json:"genus"`
	Weight float32 `json:"weight"`
	IsHealth bool `json:"is_health"`
}

type AnimalProdactList struct {
	Id           string `json:"id"`
	Name    string `json:"name"`
	CategoryName       string `json:"category_name"`
	Gender    string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Description        string `json:"description"`
	Genus      string `json:"genus"`
	Weight float32 `json:"weight"`
	IsHealth bool `json:"is_healt"`
	Products []*AnimalProdacts `json:"products"`
}

type AnimalProdacts struct {
	Name string `json:"product_name"`
	Capacity string `json:"capacity"` // animal_products capacity + products union
	GetTime string `json:"get_time"`
}

type AnimalFoodList struct {
	Id           string `json:"id"`
	Name    string `json:"name"`
	CategoryName       string `json:"category_name"`
	Gender    string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Description        string `json:"description"`
	Genus      string `json:"genus"`
	Weight float32 `json:"weight"`
	IsHealth bool `json:"is_healt"`
	Foods []*AnimalFoods `json:"foods"`
}

type AnimalFoods struct {
	Name string `json:"food_name"`
	Capacity string `json:"capacity"` // animal_given_eatables + foods union
	GivenTime string `json:"given_time"`
}

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type ListAnimalsRes struct {
	Animals []*AnimalRes `json:"animals"`
	Count int64 `json:"count"`
}

type AnimalFieldValues struct {
	Category string `json:"category"`
	Genus string `json:"genus"`
	Gender  string `json:"gender"`
	Weight float32 `json:"weight"`
	IsHealth bool `json:"is_health"`
}

type Result struct {
	Message string `json:"message"`
}

type CategoryRes struct {
	Categories []*string `json:"categories"`
}

func (t *AnimalReq) Validate() error {
	t.CategoryName = strings.ToLower(t.CategoryName)
	t.DateOfBirth = strings.ToLower(t.DateOfBirth)
	t.Description = strings.ToLower(t.Description)
	t.Gender = strings.ToLower(t.Gender)
	t.Genus = strings.ToLower(t.Genus)
	return validation.ValidateStruct(t,
		validation.Field(
			&t.Gender,
			validation.Required,
			validation.In("male", "female"),
		),
		validation.Field(
			&t.DateOfBirth,
			validation.Required,
			validation.Date(time.DateOnly),
		),
		validation.Field(
			&t.IsHealth,
			validation.In(true, false),
		),
	)
}