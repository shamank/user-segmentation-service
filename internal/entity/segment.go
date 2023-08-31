package entity

import "time"

type Segment struct {
	Id        int
	Slug      string
	CreatedAt time.Time
	DeletedAt *time.Time
}
