package models

import "time"

type AccountStatement struct {
	ID              string    `db:"id"               json:"id"`
	AccountID       string    `db:"account_id"       json:"account_id"`
	EntryID         string    `db:"entry_id"         json:"entry_id"`
	StatementDate   time.Time `db:"statement_date"   json:"statement_date"`
	Side            string    `db:"side"             json:"side"`
	IncomingBalance float64   `db:"incoming_balance" json:"incoming_balance"`
	Amount          float64   `db:"amount"           json:"amount"`
	OutgoingBalance float64   `db:"outgoing_balance" json:"outgoing_balance"`
	CreatedAt       time.Time `db:"created_at"       json:"created_at"`
	
	EntryNumber   string `db:"entry_number"    json:"entry_number,omitempty"`
	AccountNumber string `db:"account_number"  json:"account_number,omitempty"`
}
