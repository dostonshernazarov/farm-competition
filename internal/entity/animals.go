package entity

import (
	"time"
)

type Animal struct {
	ID           string
	Name         string
	CategoryName string
	Gender       string
	BirthDay     string
	Genus        string
	Weight       uint64
	IsHealth     string
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ListAnimal struct {
	Animals    []*Animal
	TotalCount uint64
}
