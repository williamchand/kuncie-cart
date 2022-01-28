package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/williamchand/kuncie-cart/order"

	"github.com/williamchand/kuncie-cart/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mysqlOrderRepository struct {
	Conn *sql.DB
}

// NewMysqlOrderRepository will create an object that represent the order.Repository interface
func NewMysqlOrderRepository(Conn *sql.DB) order.Repository {
	return &mysqlOrderRepository{Conn}
}

func (m *mysqlOrderRepository) GetItemsById(ctx context.Context, id []int64) (res []*models.Items, err error) {
	args := make([]interface{}, len(id))
	for i, val := range id {
		args[i] = val
	}
	query := `SELECT id,sku,name,price,inventory_quantity, updated_at, created_at
  						FROM items WHERE id IN (?` + strings.Repeat(",?", len(args)-1) + `)`
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.Items, 0)
	for rows.Next() {
		t := new(models.Items)
		err = rows.Scan(
			&t.ID,
			&t.SKU,
			&t.Name,
			&t.Price,
			&t.InventoryQuantity,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlOrderRepository) GetItems(ctx context.Context, sku []string) (res []*models.Items, err error) {
	args := make([]interface{}, len(sku))
	for i, skuid := range sku {
		args[i] = skuid
	}
	query := `SELECT id,sku,name,price,inventory_quantity, updated_at, created_at
  						FROM items WHERE sku IN (?` + strings.Repeat(",?", len(args)-1) + `)`
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.Items, 0)
	for rows.Next() {
		t := new(models.Items)
		err = rows.Scan(
			&t.ID,
			&t.SKU,
			&t.Name,
			&t.Price,
			&t.InventoryQuantity,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlOrderRepository) GetPromotions(ctx context.Context, id int64) (res *models.Promotions, err error) {
	query := `SELECT id, items_id, promo_type, promo, quantity_requirement
  						FROM promotions WHERE items_id = ?`
	rows, err := m.Conn.QueryContext(ctx, query, id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.Promotions, 0)
	for rows.Next() {
		t := new(models.Promotions)
		err = rows.Scan(
			&t.ID,
			&t.ItemsID,
			&t.PromoType,
			&t.Promo,
			&t.QuantityRequirement,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return res, nil
}

func (m *mysqlOrderRepository) GetCart(ctx context.Context) (res []*models.Cart, err error) {
	query := `SELECT id, items_id, quantity, updated_at, created_at
  						FROM cart`
	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.Cart, 0)
	for rows.Next() {
		t := new(models.Cart)
		err = rows.Scan(
			&t.ID,
			&t.ItemsID,
			&t.Quantity,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlOrderRepository) CreateCart(ctx context.Context, a *models.Cart) error {
	query := `INSERT cart SET items_id=?, quantity=?, updated_at=?, created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, a.ItemsID, a.Quantity, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = lastID
	return nil
}

func (m *mysqlOrderRepository) CreateOrder(ctx context.Context, a *models.Order) error {
	query := "INSERT `" + "order" + "` SET total_price=?, updated_at=?, created_at=?"
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx, a.TotalPrice, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = lastID
	return nil
}

func (m *mysqlOrderRepository) CreateOrderDetails(ctx context.Context, a *models.OrderDetails) error {
	query := `INSERT order_details SET order_id=? , sku=?, name=?, price=?, quantity=?, promo_type=?, updated_at=? , created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, a.OrderID, a.SKU, a.Name, a.Price, a.Quantity, a.PromoType, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = lastID
	return nil
}
func (m *mysqlOrderRepository) UpdateCart(ctx context.Context, ar *models.Cart) error {
	query := `UPDATE cart set items_id=?, quantity=?, updated_at=? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, ar.ItemsID, ar.Quantity, ar.UpdatedAt, ar.ID)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)

		return err
	}

	return nil
}

func (m *mysqlOrderRepository) UpdateItems(ctx context.Context, ar *models.Items) error {
	query := `UPDATE items set inventory_quantity= inventory_quantity - ?, updated_at=? WHERE sku = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, ar.InventoryQuantity, ar.UpdatedAt, ar.SKU)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)

		return err
	}

	return nil
}
func (m *mysqlOrderRepository) DeleteCart(ctx context.Context) error {
	query := "DELETE FROM cart"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx)
	if err != nil {

		return err
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAfected)
		return err
	}

	return nil
}

// DecodeCursor will decode cursor from user for mysql
func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(byt)
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}

// EncodeCursor will encode cursor from mysql to user
func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
