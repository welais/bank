package models

import "time"

type Client struct {
	ID         string    `db:"id"          json:"id"`
	LastName   string    `db:"last_name"   json:"last_name"`
	FirstName  string    `db:"first_name"  json:"first_name"`
	MiddleName string    `db:"middle_name" json:"middle_name"`
	Phone      string    `db:"phone"       json:"phone"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"  json:"updated_at"`
}
