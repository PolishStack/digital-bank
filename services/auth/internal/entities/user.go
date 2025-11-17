package entities

import "time"

type User struct {
	ID        uint64
	Email     string
	Password  string
	Role      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
