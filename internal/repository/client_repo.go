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
		`INSERT INTO clients (last_name, first_name, middle_name, phone)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
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
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM accounts WHERE client_id = $1 AND status = 'OPEN'`, id,
	).Scan(&count)
	if err != nil {
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
