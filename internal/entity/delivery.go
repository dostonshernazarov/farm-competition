package entity

import "time"

type Delivery struct {
	ID        string
	Name      string
	Category  string
	Capacity  int64
	Union     string
	Time      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListDelivery struct {
	Deliveries []*Delivery
	TotalCount int64
}
