package repository

import (
	"bank/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AccountRepo struct {
	db *sqlx.DB
}

func NewAccountRepo(db *sqlx.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) Open(clientID string, accountType string) (*models.Account, error) {
	var id, number string
	err := r.db.QueryRow(
		`CALL open_account($1, $2, NULL, NULL)`,
		clientID, accountType,
	).Scan(&id, &number)
	if err != nil {
		return nil, err
	}
	var acc models.Account
	err = r.db.Get(&acc,
		`SELECT a.id, a.client_id, a.account_number, a.account_type, a.status,
		        a.balance, a.opened_at, a.closed_at, a.created_at, a.updated_at,
		        TRIM(c.last_name || ' ' || c.first_name || ' ' || COALESCE(c.middle_name,'')) AS client_name
		 FROM accounts a
		 JOIN clients c ON c.id = a.client_id
		 WHERE a.id = $1`, id)
	return &acc, err
}

func (r *AccountRepo) Close(id string) error {
	_, err := r.db.Exec(`CALL close_account($1)`, id)
	return err
}

func (r *AccountRepo) FindAll() ([]models.Account, error) {
	var list []models.Account
	err := r.db.Select(&list,
		`SELECT a.id, a.client_id, a.account_number, a.account_type, a.status,
		        a.balance, a.opened_at, a.closed_at, a.created_at, a.updated_at,
		        TRIM(c.last_name || ' ' || c.first_name || ' ' || COALESCE(c.middle_name,'')) AS client_name
		 FROM accounts a
		 JOIN clients c ON c.id = a.client_id
		 ORDER BY a.opened_at DESC`)
	return list, err
}

func (r *AccountRepo) FindByNumber(number string) (*models.Account, error) {
	var acc models.Account
	err := r.db.Get(&acc,
		`SELECT a.id, a.client_id, a.account_number, a.account_type, a.status,
		        a.balance, a.opened_at, a.closed_at, a.created_at, a.updated_at,
		        TRIM(c.last_name || ' ' || c.first_name || ' ' || COALESCE(c.middle_name,'')) AS client_name
		 FROM accounts a
		 JOIN clients c ON c.id = a.client_id
		 WHERE a.account_number = $1`, number)
	return &acc, err
}

func (r *AccountRepo) FindByClientLastName(lastName string) ([]models.Account, error) {
	var list []models.Account
	err := r.db.Select(&list,
		`SELECT a.id, a.client_id, a.account_number, a.account_type, a.status,
		        a.balance, a.opened_at, a.closed_at, a.created_at, a.updated_at,
		        TRIM(c.last_name || ' ' || c.first_name || ' ' || COALESCE(c.middle_name,'')) AS client_name
		 FROM accounts a
		 JOIN clients c ON c.id = a.client_id
		 WHERE c.last_name ILIKE $1
		 ORDER BY a.opened_at DESC`,
		"%"+lastName+"%")
	return list, err
}

func (r *AccountRepo) GetByID(id string) (*models.Account, error) {
	var acc models.Account
	err := r.db.Get(&acc,
		`SELECT id, client_id, account_number, account_type, status,
		        balance, opened_at, closed_at, created_at, updated_at
		 FROM accounts WHERE id = $1`, id)
	return &acc, err
}

func (r *AccountRepo) Delete(id string) error {
	if id == "00000000-0000-0000-0000-000000000002" {
		return fmt.Errorf("нельзя удалить счёт Касса")
	}
	res, err := r.db.Exec(
		`DELETE FROM accounts WHERE id = $1 AND status = 'CLOSE'`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("счёт не найден или ещё открыт")
	}
	return nil
}

func (r *AccountRepo) DeleteAllClosed() (int64, error) {
	res, err := r.db.Exec(`
		DELETE FROM accounts
		WHERE id != '00000000-0000-0000-0000-000000000002'
		  AND status = 'CLOSE'`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
