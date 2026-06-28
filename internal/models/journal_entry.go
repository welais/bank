package models

import "time"

type JournalEntry struct {
	ID               string    `db:"id"                json:"id"`
	EntryNumber      string    `db:"entry_number"      json:"entry_number"`
	EntryDate        time.Time `db:"entry_date"        json:"entry_date"`
	DebitAccountID   string    `db:"debit_account_id"  json:"debit_account_id"`
	CreditAccountID  string    `db:"credit_account_id" json:"credit_account_id"`
	Amount           float64   `db:"amount"            json:"amount"`
	PaymentPurpose   string    `db:"payment_purpose"   json:"payment_purpose"`
	CreatedAt        time.Time `db:"created_at"        json:"created_at"`

	// Joined
	DebitAccountNumber  string `db:"debit_account_number"  json:"debit_account_number,omitempty"`
	CreditAccountNumber string `db:"credit_account_number" json:"credit_account_number,omitempty"`
}
