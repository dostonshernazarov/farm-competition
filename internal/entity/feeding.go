package entity

import "time"

type Feeding struct {
	ID         string
	AnimalID   string
	EatablesID string
	Category   string
	Day        string
	Daily      []struct {
		Capacity int64  `json:"capacity"`
		Time     string `json:"time"`
	} `json:"daily"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FeedingRes struct {
	ID       string
	AnimalID string
	Eatables struct {
		ID          string
		Name        string
		Status      string
		Capacity    uint64
		Union       string
		Description string
	}
	Category string
	Day      string
	Daily    []struct {
		Capacity int64  `json:"capacity"`
		Time     string `json:"time"`
	} `json:"daily"`
}
