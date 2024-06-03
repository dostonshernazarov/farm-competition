package postgresql

import (
	"context"
	"database/sql"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
	"time"
)

type productRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewProduct(db *postgres.PostgresDB) repo.Product {
	return &productRepo{
		tableName: "products",
		db:        db,
	}
}

func (a *productRepo) Create(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	query := `
	INSERT INTO products (
	    id,
		name,
		product_union,
		description,
	    total_capacity,
		created_at,
		updated_at
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING
		id,
		name,
		product_union,
		description,
		total_capacity
	`

	var (
		createdProduct     entity.Product
		sqlNullDescription sql.NullString
	)

	err := a.db.QueryRow(ctx, query,
		product.ID,
		product.Name,
		product.Union,
		product.Description,
		product.TotalCapacity,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(
		&createdProduct.ID,
		&createdProduct.Name,
		&createdProduct.Union,
		&createdProduct.TotalCapacity,
		&sqlNullDescription,
	)

	if err != nil {
		return nil, err
	}

	if sqlNullDescription.Valid {
		createdProduct.Description = sqlNullDescription.String
	}

	return &createdProduct, nil
}

func (a *productRepo) Update(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	query := `
	UPDATE
		products
	SET
		name = $1,
		product_union = $2,
		description = $3,
	    total_capacity = $4,
		updated_at = $5
	WHERE
	    id = $6
		AND deleted_at IS NULL
	RETURNING
		id,
		name,
		product_union,
		description,
		total_capacity
	`

	var (
		updatedProduct     entity.Product
		sqlNullDescription sql.NullString
	)

	err := a.db.QueryRow(ctx, query,
		product.Name,
		product.Union,
		product.Description,
		product.TotalCapacity,
		product.UpdatedAt,
		product.ID,
	).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Union,
		&updatedProduct.TotalCapacity,
		&sqlNullDescription,
	)

	if err != nil {
		return nil, err
	}

	if sqlNullDescription.Valid {
		updatedProduct.Description = sqlNullDescription.String
	}

	return &updatedProduct, nil
}

func (a *productRepo) Delete(ctx context.Context, productID string) error {
	query := `UPDATE products SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := a.db.Exec(ctx, query, productID, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}

	if result.RowsAffected() != 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (a *productRepo) Get(ctx context.Context, productID string) (*entity.Product, error) {
	query := `
	SELECT
		id,
		name,
		product_union,
		description,
		total_capacity
	FROM
	    products
	WHERE
	    id = $1
		AND deleted_at IS NULL
	`

	var (
		product            entity.Product
		sqlNullDescription sql.NullString
	)

	err := a.db.QueryRow(ctx, query, productID).Scan(
		&product.ID,
		&product.Name,
		&product.Union,
		&product.TotalCapacity,
		&sqlNullDescription,
	)

	if err != nil {
		return nil, err
	}

	if sqlNullDescription.Valid {
		product.Description = sqlNullDescription.String
	}

	return &product, nil
}

func (a *productRepo) List(ctx context.Context, page, limit uint64) (*entity.ListProducts, error) {
	query := `
	SELECT
		id,
		name,
		product_union,
		description,
		total_capacity
	FROM
	    products
	WHERE
	    deleted_at IS NULL
	LIMIT $1
	OFFSET $2
	`

	var (
		products entity.ListProducts
		offset   = limit * (page - 1)
	)

	rows, err := a.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			product            entity.Product
			sqlNullDescription sql.NullString
		)
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Union,
			&product.TotalCapacity,
			&sqlNullDescription,
		)
		if err != nil {
			return nil, err
		}

		if sqlNullDescription.Valid {
			product.Description = sqlNullDescription.String
		}

		products.Products = append(products.Products, &product)
	}
	var (
		count      = 0
		totalQuery = `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL`
	)

	if err := a.db.QueryRow(ctx, totalQuery).Scan(&count); err != nil {
		return nil, err
	}
	products.TotalCount = uint64(count)

	return &products, nil
}
