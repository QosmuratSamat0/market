package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/QosmuratSamat0/product-service/internal/domain/product"
	"github.com/QosmuratSamat0/product-service/internal/lib/errs"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// ---------- Products ----------

func (r *PostgresRepo) CreateProduct(ctx context.Context, p *product.Product) error {
	query := `
		INSERT INTO products (name, description, price, category_id, seller_id, image_url, stock)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		p.Name, p.Description, p.Price, p.CategoryID, p.SellerID, p.ImageURL, p.Stock,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PostgresRepo) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	query := `
		SELECT id, name, description, price, category_id, seller_id,
		       COALESCE(image_url, ''), stock, created_at, updated_at
		FROM products
		WHERE id = $1`

	p := &product.Product{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price,
		&p.CategoryID, &p.SellerID, &p.ImageURL,
		&p.Stock, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ProductNotFound
		}
		return nil, fmt.Errorf("get product: %w", err)
	}
	return p, nil
}

func (r *PostgresRepo) GetAllProducts(ctx context.Context) ([]*product.Product, error) {
	query := `
		SELECT id, name, description, price, category_id, seller_id,
		       COALESCE(image_url, ''), stock, created_at, updated_at
		FROM products
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all products: %w", err)
	}
	defer rows.Close()

	var products []*product.Product
	for rows.Next() {
		p := &product.Product{}
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price,
			&p.CategoryID, &p.SellerID, &p.ImageURL,
			&p.Stock, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *PostgresRepo) GetProductsByCategory(ctx context.Context, categoryID string) ([]*product.Product, error) {
	query := `
		SELECT id, name, description, price, category_id, seller_id,
		       COALESCE(image_url, ''), stock, created_at, updated_at
		FROM products
		WHERE category_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("get products by category: %w", err)
	}
	defer rows.Close()

	var products []*product.Product
	for rows.Next() {
		p := &product.Product{}
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price,
			&p.CategoryID, &p.SellerID, &p.ImageURL,
			&p.Stock, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *PostgresRepo) GetProductsBySeller(ctx context.Context, sellerID string) ([]*product.Product, error) {
	query := `
		SELECT id, name, description, price, category_id, seller_id,
		       COALESCE(image_url, ''), stock, created_at, updated_at
		FROM products
		WHERE seller_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, sellerID)
	if err != nil {
		return nil, fmt.Errorf("get products by seller: %w", err)
	}
	defer rows.Close()

	var products []*product.Product
	for rows.Next() {
		p := &product.Product{}
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price,
			&p.CategoryID, &p.SellerID, &p.ImageURL,
			&p.Stock, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *PostgresRepo) UpdateProduct(ctx context.Context, p *product.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, category_id = $4,
		    image_url = $5, stock = $6, updated_at = NOW()
		WHERE id = $7`

	ct, err := r.db.Exec(ctx, query,
		p.Name, p.Description, p.Price, p.CategoryID,
		p.ImageURL, p.Stock, p.ID,
	)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errs.ProductNotFound
	}
	return nil
}

func (r *PostgresRepo) DeleteProduct(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errs.ProductNotFound
	}
	return nil
}

// ---------- Categories ----------

func (r *PostgresRepo) CreateCategory(ctx context.Context, c *product.Category) error {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, c.Name).Scan(&c.ID, &c.CreatedAt)
}

func (r *PostgresRepo) GetAllCategories(ctx context.Context) ([]*product.Category, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, created_at FROM categories ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("get all categories: %w", err)
	}
	defer rows.Close()

	var categories []*product.Category
	for rows.Next() {
		c := &product.Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *PostgresRepo) GetCategoryByID(ctx context.Context, id string) (*product.Category, error) {
	c := &product.Category{}
	err := r.db.QueryRow(ctx, `SELECT id, name, created_at FROM categories WHERE id = $1`, id).
		Scan(&c.ID, &c.Name, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.CategoryNotFound
		}
		return nil, fmt.Errorf("get category: %w", err)
	}
	return c, nil
}

func (r *PostgresRepo) DeleteCategory(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errs.CategoryNotFound
	}
	return nil
}
