package repository

import (
	"bank/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type JournalEntryRepo struct {
	db *sqlx.DB
}

func NewJournalEntryRepo(db *sqlx.DB) *JournalEntryRepo {
	return &JournalEntryRepo{db: db}
}

const journalSelectJoined = `
	SELECT j.id, j.entry_number, j.entry_date, j.debit_account_id, j.credit_account_id,
	       j.amount, j.payment_purpose, j.created_at,
	       da.account_number AS debit_account_number,
	       ca.account_number AS credit_account_number
	FROM journal_entries j
	JOIN accounts da ON da.id = j.debit_account_id
	JOIN accounts ca ON ca.id = j.credit_account_id`

func (r *JournalEntryRepo) FindAll() ([]models.JournalEntry, error) {
	var list []models.JournalEntry
	err := r.db.Select(&list, journalSelectJoined+` ORDER BY j.entry_date DESC, j.created_at DESC`)
	return list, err
}

func (r *JournalEntryRepo) FindByDate(date time.Time) ([]models.JournalEntry, error) {
	var list []models.JournalEntry
	err := r.db.Select(&list, journalSelectJoined+` WHERE j.entry_date = $1 ORDER BY j.created_at`, date.Format("2006-01-02"))
	return list, err
}

func (r *JournalEntryRepo) Create(e *models.JournalEntry) error {
	return r.db.QueryRow(
		`CALL add_journal_entry($1,$2,$3,$4,$5,$6,NULL)`,
		e.EntryNumber, e.EntryDate, e.DebitAccountID, e.CreditAccountID, e.Amount, e.PaymentPurpose,
	).Scan(&e.ID)
}

func (r *JournalEntryRepo) Delete(id string) error {
	_, err := r.db.Exec(`CALL delete_journal_entry($1)`, id)
	return err
}
