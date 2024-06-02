package entity

import "time"

type Category struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListCategory struct {
	Categories []*Category
	TotalCount uint64
}
