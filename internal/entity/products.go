package entity

import "time"

type Product struct {
	ID            string
	Name          string
	Union         string
	TotalCapacity int64
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ListProducts struct {
	Products   []*Product
	TotalCount uint64
}
