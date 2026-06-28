package repository

import (
	"bank/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type StatementRepo struct {
	db *sqlx.DB
}

func NewStatementRepo(db *sqlx.DB) *StatementRepo {
	return &StatementRepo{db: db}
}

func (r *StatementRepo) GetByAccount(accountID string, from, to time.Time) ([]models.AccountStatement, error) {
	var list []models.AccountStatement
	err := r.db.Select(&list, `
		SELECT s.id, s.account_id, s.entry_id, s.statement_date, s.side,
		       s.incoming_balance, s.amount, s.outgoing_balance, s.created_at,
		       j.entry_number, a.account_number
		FROM account_statements s
		JOIN journal_entries j ON j.id = s.entry_id
		JOIN accounts a        ON a.id = s.account_id
		WHERE s.account_id = $1
		  AND s.statement_date BETWEEN $2 AND $3
		ORDER BY s.statement_date, s.created_at`,
		accountID, from.Format("2006-01-02"), to.Format("2006-01-02"))
	return list, err
}

func (r *StatementRepo) GetAll(from, to time.Time) ([]models.AccountStatement, error) {
	var list []models.AccountStatement
	err := r.db.Select(&list, `
		SELECT s.id, s.account_id, s.entry_id, s.statement_date, s.side,
		       s.incoming_balance, s.amount, s.outgoing_balance, s.created_at,
		       j.entry_number, a.account_number
		FROM account_statements s
		JOIN journal_entries j ON j.id = s.entry_id
		JOIN accounts a        ON a.id = s.account_id
		WHERE s.statement_date BETWEEN $1 AND $2
		ORDER BY s.statement_date, s.created_at`,
		from.Format("2006-01-02"), to.Format("2006-01-02"))
	return list, err
}
