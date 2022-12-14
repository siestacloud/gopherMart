package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/pkg"
)

// OrderPostgres
type OrderPostgres struct {
	db *sqlx.DB
}

// NewTodoOrderPostgres
func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{
		db: db,
	}
}

// Create транзакция. Добавляю заказ в базу и связывую с клиентом
func (o *OrderPostgres) Create(userID int, order core.Order) error {
	if o.db == nil {
		return errors.New("database are not connected")
	}
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}
	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (user_order,status,sum,update_time) VALUES ($1,$2,$3,$4) RETURNING id", ordersTable)
	row := tx.QueryRow(createListQuery, order.Number, order.Status, order.Sum, NewNullString(order.CreateTime))
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, order_id) VALUES ($1, $2)", userOrderTable)
	_, err = tx.Exec(createUsersListQuery, userID, id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

// GetByNameType Получаю user из базы
func (o *OrderPostgres) GetUserByOrder(order string) (int, error) {
	if o.db == nil {
		return 0, errors.New("database are not connected")
	}
	var orderID int
	var userID int

	query := fmt.Sprintf(`SELECT id FROM %s  WHERE user_order = $1`, ordersTable)
	if err := o.db.Get(&orderID, query, order); err != nil {
		return 0, err
	}

	query = fmt.Sprintf(`SELECT user_id FROM %s  WHERE order_id = $1`,
		userOrderTable)
	if err := o.db.Get(&userID, query, orderID); err != nil {
		return 0, err
	}

	return userID, nil
}

// GetByNameType Получаю список заказов клиента
func (o *OrderPostgres) GetListOrders(userID int) ([]core.Order, error) {
	if o.db == nil {
		return nil, errors.New("database are not connected")
	}
	orderList := []core.Order{}
	query := fmt.Sprintf(`SELECT user_order, status, sum, update_time FROM %s `, ordersTable)
	if err := o.db.Select(&orderList, query); err != nil {
		pkg.ErrPrint("repository", 204, err)
		return nil, err
	}
	return orderList, nil
}
