package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"example.com/online-shop/product"
)

func GetProducts(db *sql.DB) ([]product.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT *
		FROM products;
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve products: %w", err)
	}
	defer rows.Close()

	var products []product.Product

	for rows.Next() {
		var p product.Product

		err := rows.Scan(
			&p.ProductID,
			&p.CategoryID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		products = append(products, p)
	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("error iterating over products: %w", err)
	}

	return products, nil
}

func GetProductsByID(db *sql.DB, id int) (product.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT *
		FROM products
		WHERE productid = $1;
	`

	var p product.Product

	err := db.QueryRowContext(ctx, query, id).Scan(
		&p.ProductID,
		&p.CategoryID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
	)

	if err != nil {
		return product.Product{}, fmt.Errorf("failed to retrieve product: %w", err)
	}

	return p, nil
}

func CreateProduct(db *sql.DB, p product.Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO products (categoryid, name, description, price, stock)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING productid;
	`

	var productID int

	err := db.QueryRowContext(
		ctx,
		query,
		p.CategoryID,
		p.Name,
		p.Description,
		p.Price,
		p.Stock,
	).Scan(&productID)

	if err != nil {
		return 0, fmt.Errorf("failed to create product: %w", err)
	}

	return productID, nil
}

func UpdateProduct(db *sql.DB, p product.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE products
		SET categoryid = $1,
			name = $2,
			description = $3,
			price = $4,
			stock = $5
		WHERE productid = $6;
	`

	result, err := db.ExecContext(
		ctx,
		query,
		p.CategoryID,
		p.Name,
		p.Description,
		p.Price,
		p.Stock,
		p.ProductID,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func PatchProduct(db *sql.DB, update product.ProductUpdate, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE products
		SET name = COALESCE($1, name),
			price = COALESCE($2, price),
			stock = COALESCE($3, stock)
		WHERE productid = $4;
	`

	result, err := db.ExecContext(
		ctx,
		query,
		update.Name,
		update.Price,
		update.Stock,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteProduct(db *sql.DB, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM products
		WHERE productid = $1;
	`

	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
