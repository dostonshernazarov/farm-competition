package entity

import "time"

type Food struct {
	ID          string
	Name        string
	Capacity    uint64
	Union       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ListFoods struct {
	Foods      []*Food
	TotalCount uint64
}
