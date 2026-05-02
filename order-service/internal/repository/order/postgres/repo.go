package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/QosmuratSamat/order-service/internal/domain/order"
	"github.com/QosmuratSamat/order-service/internal/lib/errs"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) CreateOrder(ctx context.Context, o *order.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if o.Status == "" {
		o.Status = order.StatusPending
	}

	query := `
		INSERT INTO orders (user_id, status, total)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	if err := tx.QueryRow(ctx, query, o.UserID, o.Status, o.Total).
		Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return fmt.Errorf("create order for user %s: %w", o.UserID, err)
	}

	for i := range o.Items {
		item := &o.Items[i]
		item.OrderID = o.ID

		query := `
			INSERT INTO order_items (order_id, product_id, seller_id, quantity, price)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id`

		if err := tx.QueryRow(ctx, query,
			item.OrderID, item.ProductID, item.SellerID, item.Quantity, item.Price,
		).Scan(&item.ID); err != nil {
			return fmt.Errorf("create order item for order %s product %s: %w", o.ID, item.ProductID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetOrderByID(ctx context.Context, id string) (*order.Order, error) {
	query := `
		SELECT id, user_id, status, total, created_at, updated_at
		FROM orders
		WHERE id = $1`

	o := &order.Order{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.UserID, &o.Status, &o.Total, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.OrderNotFound
		}
		return nil, fmt.Errorf("get order %s: %w", id, err)
	}

	items, err := r.getOrderItems(ctx, o.ID)
	if err != nil {
		return nil, fmt.Errorf("get items for order %s: %w", o.ID, err)
	}
	o.Items = items

	return o, nil
}

func (r *PostgresRepo) GetOrdersByUser(ctx context.Context, userID string) ([]*order.Order, error) {
	query := `
		SELECT id, user_id, status, total, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC`

	return r.getOrders(ctx, query, userID)
}

func (r *PostgresRepo) GetOrdersBySeller(ctx context.Context, sellerID string) ([]*order.Order, error) {
	query := `
		SELECT DISTINCT o.id, o.user_id, o.status, o.total, o.created_at, o.updated_at
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id
		WHERE oi.seller_id = $1
		ORDER BY o.created_at DESC`

	return r.getOrders(ctx, query, sellerID)
}

func (r *PostgresRepo) UpdateOrderStatus(ctx context.Context, id string, status order.Status) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2`

	ct, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update order %s status: %w", id, err)
	}
	if ct.RowsAffected() == 0 {
		return errs.OrderNotFound
	}
	return nil
}

func (r *PostgresRepo) DeleteOrder(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete order %s: %w", id, err)
	}
	if ct.RowsAffected() == 0 {
		return errs.OrderNotFound
	}
	return nil
}

func (r *PostgresRepo) getOrders(ctx context.Context, query string, args ...any) ([]*order.Order, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get orders: %w", err)
	}
	defer rows.Close()

	var orders []*order.Order
	var orderIDs []string
	for rows.Next() {
		o := &order.Order{}
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.Total, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		orders = append(orders, o)
		orderIDs = append(orderIDs, o.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	itemsByOrderID, err := r.getOrderItemsByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, err
	}
	for _, o := range orders {
		o.Items = itemsByOrderID[o.ID]
	}

	return orders, nil
}

func (r *PostgresRepo) getOrderItems(ctx context.Context, orderID string) ([]order.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, seller_id, quantity, price
		FROM order_items
		WHERE order_id = $1
		ORDER BY id`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("get order items: %w", err)
	}
	defer rows.Close()

	var items []order.OrderItem
	for rows.Next() {
		var item order.OrderItem
		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID,
			&item.SellerID, &item.Quantity, &item.Price,
		); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order items: %w", err)
	}

	return items, nil
}

func (r *PostgresRepo) getOrderItemsByOrderIDs(ctx context.Context, orderIDs []string) (map[string][]order.OrderItem, error) {
	itemsByOrderID := make(map[string][]order.OrderItem, len(orderIDs))
	if len(orderIDs) == 0 {
		return itemsByOrderID, nil
	}

	placeholders := make([]string, 0, len(orderIDs))
	args := make([]any, 0, len(orderIDs))
	for i, id := range orderIDs {
		placeholders = append(placeholders, fmt.Sprintf("$%d::uuid", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT id, order_id, product_id, seller_id, quantity, price
		FROM order_items
		WHERE order_id IN (%s)
		ORDER BY order_id, id`, strings.Join(placeholders, ", "))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get items for %d orders: %w", len(orderIDs), err)
	}
	defer rows.Close()

	for rows.Next() {
		var item order.OrderItem
		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID,
			&item.SellerID, &item.Quantity, &item.Price,
		); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order items: %w", err)
	}

	return itemsByOrderID, nil
}
