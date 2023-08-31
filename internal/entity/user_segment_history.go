package entity

import "time"

type UserSegmentHistory struct {
	ID        int
	UserID    int
	SegmentID int
	Operation string
	CreatedAt time.Time
}
