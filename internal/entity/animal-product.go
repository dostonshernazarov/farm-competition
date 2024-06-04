package entity

import "time"

type AnimalProductReq struct {
	ID        string
	AnimalID  string
	ProductID string
	Capacity  int64
	GetTime   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AnimalProductRes struct {
	ID       string
	Animal   Animal
	Product  Product
	Capacity int64
	GetTime  string
}

type ListAnimalProduct struct {
	AnimalProducts []*AnimalProductRes
	TotalCount     uint64
}

type ProductsWithAnimal struct {
	Animal   Animal
	Products []*struct {
		ID            string
		Name          string
		Union         string
		Description   string
		TotalCapacity int64
	}
}

type AnimalsWithProduct struct {
	Product Product
	Animals []*struct {
		ID            string
		Name          string
		CategoryName  string
		Gender        string
		BirthDay      string
		Genus         string
		Weight        uint64
		IsHealth      string
		Description   string
		TotalCapacity int64
	}
}
