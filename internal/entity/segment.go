package entity

import "time"

type Segment struct {
	Id               int
	Slug             string
	AssignPercentage int
	CreatedAt        time.Time
	DeletedAt        *time.Time
}
