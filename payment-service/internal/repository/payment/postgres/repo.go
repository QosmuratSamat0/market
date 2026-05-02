package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	domain "github.com/QosmuratSamat0/payment-service/internal/domain/payment"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pool: pool}
}

func (r *PostgresRepo) Create(ctx context.Context, p domain.Payment) error {
	metadata, err := marshalMetadata(p.Metadata)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO payments (id, order_id, user_id, amount, currency, provider, status, idempotency_key, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = r.pool.Exec(ctx, query,
		p.ID, p.OrderID, p.UserID, p.Amount, p.Currency,
		p.Provider, p.Status, p.IdempotencyKey, metadata,
	)
	return err
}

func (r *PostgresRepo) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, user_id, amount, currency, provider, provider_id, status, idempotency_key, metadata, created_at
		FROM payments 
		WHERE idempotency_key = $1
	`
	var p domain.Payment
	var providerID *string
	var metadata []byte
	err := r.pool.QueryRow(ctx, query, key).Scan(
		&p.ID, &p.OrderID, &p.UserID, &p.Amount, &p.Currency,
		&p.Provider, &providerID, &p.Status, &p.IdempotencyKey, &metadata, &p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if providerID != nil {
		p.ProviderID = *providerID
	}
	if err := unmarshalMetadata(metadata, &p.Metadata); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostgresRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus, externalID string) error {
	query := `
		UPDATE payments 
		SET status = $1, provider_id = COALESCE($2, provider_id), updated_at = NOW() 
		WHERE id = $3
	`
	cmd, err := r.pool.Exec(ctx, query, status, nullString(externalID), id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("payment %s not found", id)
	}
	return err
}

func (r *PostgresRepo) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]domain.Payment, error) {
	query := `
		SELECT id, order_id, user_id, amount, currency, provider, provider_id, status, idempotency_key, created_at
		FROM payments 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query payments: %w", err)
	}
	defer rows.Close()

	var payments []domain.Payment
	for rows.Next() {
		var p domain.Payment
		var providerID *string
		err := rows.Scan(
			&p.ID, &p.OrderID, &p.UserID, &p.Amount, &p.Currency,
			&p.Provider, &providerID, &p.Status, &p.IdempotencyKey, &p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		if providerID != nil {
			p.ProviderID = *providerID
		}
		payments = append(payments, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func marshalMetadata(metadata map[string]interface{}) ([]byte, error) {
	if metadata == nil {
		return nil, nil
	}
	body, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment metadata: %w", err)
	}
	return body, nil
}

func unmarshalMetadata(body []byte, metadata *map[string]interface{}) error {
	if len(body) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, metadata); err != nil {
		return fmt.Errorf("failed to unmarshal payment metadata: %w", err)
	}
	return nil
}

func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
