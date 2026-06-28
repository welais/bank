package models

import "time"

type Account struct {
	ID            string     `db:"id"             json:"id"`
	ClientID      string     `db:"client_id"      json:"client_id"`
	AccountNumber string     `db:"account_number" json:"account_number"`
	AccountType   string     `db:"account_type"   json:"account_type"`
	Status        string     `db:"status"         json:"status"`
	Balance       float64    `db:"balance"        json:"balance"`
	OpenedAt      time.Time  `db:"opened_at"      json:"opened_at"`
	ClosedAt      *time.Time `db:"closed_at"     json:"closed_at"`
	CreatedAt     time.Time  `db:"created_at"     json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"     json:"updated_at"`

	ClientName string `db:"client_name" json:"client_name,omitempty"`
}
