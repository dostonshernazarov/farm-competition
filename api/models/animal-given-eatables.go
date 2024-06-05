package models

import (
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AnimaGivenEatablesReq struct {
	AnimalID   string   `json:"animal_id"`
	EatablesID string   `json:"eatables_id"`
	Daily      []*Daily `json:"daily"`
	Category   string   `json:"category"`
	Day        string   `json:"day"`
}

type AnimaDrugGivenEatablesRes struct {
	ID       string   `json:"id"`
	AnimalID string   `json:"animal_id"`
	Eatables DrugRes  `json:"eatables"`
	Daily    []*Daily `json:"daily"`
	Category string   `json:"category"`
	Day      string   `json:"day"`
}

type AnimaFoodGivenEatablesRes struct {
	ID       string   `json:"id"`
	AnimalID string   `json:"animal_id"`
	Eatables FoodRes  `json:"eatables"`
	Daily    []*Daily `json:"daily"`
	Category string   `json:"category"`
	Day      string   `json:"day"`
}

type AnimaGivenEatablesRes struct {
	ID         string   `json:"id"`
	AnimalID   string   `json:"animal_id"`
	EatablesID string   `json:"eatables_id"`
	Daily      []*Daily `json:"daily"`
	Category   string   `json:"category"`
	Day      string   `json:"day"`
}

func (t *AnimaGivenEatablesReq) Validate() error {
	t.Category = strings.ToLower(t.Category)
	_, err := time.Parse(time.DateOnly, t.Day)
	if err != nil {
		return errors.New("invalide type")
	}
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
