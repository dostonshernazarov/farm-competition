package postgresql

import (
	"context"
	"database/sql"
	"github.com/spf13/cast"
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
		&sqlNullDescription,
		&createdProduct.TotalCapacity,
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
		&sqlNullDescription,
		&updatedProduct.TotalCapacity,
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

	result, err := a.db.Exec(ctx, query, time.Now().Format(time.RFC3339), productID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (a *productRepo) Get(ctx context.Context, params map[string]string) (*entity.Product, error) {
	var (
		product            entity.Product
		sqlNullDescription sql.NullString
	)

	queryBuilder := a.db.Sq.Builder.Select("id, name, product_union, description, total_capacity")
	queryBuilder = queryBuilder.From(a.tableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(a.db.Sq.Equal(key, value))
		}
		if key == "name" {
			queryBuilder = queryBuilder.Where(a.db.Sq.Equal(key, value))
		}
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	err = a.db.QueryRow(ctx, query, args...).Scan(
		&product.ID,
		&product.Name,
		&product.Union,
		&sqlNullDescription,
		&product.TotalCapacity,
	)

	if err != nil {
		return nil, err
	}

	if sqlNullDescription.Valid {
		product.Description = sqlNullDescription.String
	}

	return &product, nil
}

func (p *productRepo) List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListProducts, error) {
	var (
		products entity.ListProducts
		offset   = limit * (page - 1)
	)

	queryBuilder := p.db.Sq.Builder.Select("id, name, product_union, total_capacity")
	queryBuilder = queryBuilder.From(p.tableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(p.db.Sq.ILike("name", "%"+cast.ToString(params["name"])+"%"))
	queryBuilder = queryBuilder.Where(p.db.Sq.ILike("product_union", "%"+cast.ToString(params["union"])+"%"))
	queryBuilder = queryBuilder.Limit(limit)
	queryBuilder = queryBuilder.Offset(offset)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(ctx, query, args...)
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

	totalQueryBuilder := p.db.Sq.Builder.Select("COUNT(*)")
	totalQueryBuilder = totalQueryBuilder.From(p.tableName)
	totalQueryBuilder = totalQueryBuilder.Where("deleted_at IS NULL")
	totalQueryBuilder = totalQueryBuilder.Where(p.db.Sq.ILike("name", "%"+cast.ToString(params["name"])+"%"))
	totalQueryBuilder = totalQueryBuilder.Where(p.db.Sq.ILike("product_union", "%"+cast.ToString(params["union"])+"%"))

	totalQuery, totalArgs, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var count = 0
	if err := p.db.QueryRow(ctx, totalQuery, totalArgs...).Scan(&count); err != nil {
		return nil, err
	}
	products.TotalCount = uint64(count)

	return &products, nil
}

func (d *productRepo) UniqueProductName(ctx context.Context, productName string) (int, error) {
	query := `SELECT COUNT(*) FROM products WHERE name = $1 AND deleted_at IS NULL`
	var count int

	if err := d.db.QueryRow(ctx, query, productName).Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}
