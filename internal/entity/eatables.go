package entity

import "time"

type Eatables struct {
	ID        string
	AnimalID  string
	EatableID string
	Category  string
	Daily     []struct {
		Capacity int64  `json:"capacity"`
		Time     string `json:"time"`
	} `json:"daily"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EatablesRes struct {
	ID       string
	AnimalID string
	Eatable  struct {
		ID          string
		Name        string
		Status      string
		Capacity    uint64
		Union       string
		Description string
	}
	Category string
	Daily    []struct {
		Capacity int64  `json:"capacity"`
		Time     string `json:"time"`
	} `json:"daily"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EatablesFoodRes struct {
	ID       string
	AnimalID string
	Food     Food
	Daily    []struct {
		Capacity int64  `json:"capacity"`
		Time     string `json:"time"`
	} `json:"daily"`
}

type EatablesDrugRes struct {
	ID       string
	AnimalID string
	Drug     Drug
	Daily    []struct {
		Capacity int64  `json:"capacity"`
		Time     string `json:"time"`
	} `json:"daily"`
}

type ListFoodEatables struct {
	Eatables   []*EatablesFoodRes
	TotalCount uint64
}

type ListDrugEatables struct {
	Eatables   []*EatablesDrugRes
	TotalCount uint64
}
