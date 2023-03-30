package model

import (
	"time"
)

type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserFilter struct {
	Name      string    `query:"name"`
	UpdatedAt time.Time `query:"updated_at"`
}
