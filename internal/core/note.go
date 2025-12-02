package core

import "time"

type Note struct {
	ID        int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
