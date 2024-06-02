package entity

import "time"

type Drug struct {
	ID          string
	Name        string
	Status      string
	Capacity    uint64
	Union       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ListDrugs struct {
	Drugs      []*Drug
	TotalCount uint64
}
