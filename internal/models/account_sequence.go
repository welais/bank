package models

import "time"

type AccountNumberSequence struct {
	ID         string    `db:"id"          json:"id"`
	Key        string    `db:"key"         json:"key"`
	NextNumber int64     `db:"next_number" json:"next_number"`
	Prefix     string    `db:"prefix"      json:"prefix"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
}
