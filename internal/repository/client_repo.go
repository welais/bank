package repository

import (
	"bank/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ClientRepo struct {
	db *sqlx.DB
}

func NewClientRepo(db *sqlx.DB) *ClientRepo {
	return &ClientRepo{db: db}
}

func (r *ClientRepo) Create(c *models.Client) error {
	return r.db.QueryRow(
		`CALL add_client($1, $2, $3, $4, NULL)`,
		c.LastName, c.FirstName, c.MiddleName, c.Phone,
	).Scan(&c.ID)
}

func (r *ClientRepo) Update(c *models.Client) error {
	_, err := r.db.Exec(
		`CALL update_client($1, $2, $3, $4, $5)`,
		c.ID, c.LastName, c.FirstName, c.MiddleName, c.Phone,
	)
	return err
}

func (r *ClientRepo) FindAll() ([]models.Client, error) {
	var list []models.Client
	err := r.db.Select(&list,
		`SELECT id, last_name, first_name, middle_name, phone, created_at, updated_at
		 FROM clients
		 ORDER BY last_name`)
	return list, err
}

func (r *ClientRepo) FindByLastName(lastName string) ([]models.Client, error) {
	var list []models.Client
	err := r.db.Select(&list,
		`SELECT id, last_name, first_name, middle_name, phone, created_at, updated_at
		 FROM clients
		 WHERE last_name ILIKE $1
		 ORDER BY last_name`,
		"%"+lastName+"%")
	return list, err
}

func (r *ClientRepo) FindByPhone(phone string) ([]models.Client, error) {
	var list []models.Client
	err := r.db.Select(&list,
		`SELECT id, last_name, first_name, middle_name, phone, created_at, updated_at
		 FROM clients
		 WHERE phone ILIKE $1`,
		"%"+phone+"%")
	return list, err
}

func (r *ClientRepo) Delete(id string) error {
	if id == "00000000-0000-0000-0000-000000000001" {
		return fmt.Errorf("нельзя удалить системного клиента Касса")
	}
	var count int
	if err := r.db.QueryRow(
		`SELECT COUNT(*) FROM accounts WHERE client_id = $1 AND status = 'OPEN'`, id,
	).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("нельзя удалить клиента: у него есть открытые счета (%d шт.)", count)
	}
	res, err := r.db.Exec(`DELETE FROM clients WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("клиент не найден")
	}
	return nil
}

func (r *ClientRepo) DeleteAllClosed() (int64, error) {
	res, err := r.db.Exec(`
		DELETE FROM clients
		WHERE id != '00000000-0000-0000-0000-000000000001'
		  AND NOT EXISTS (
		      SELECT 1 FROM accounts
		      WHERE accounts.client_id = clients.id
		        AND accounts.status = 'OPEN'
		  )`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
