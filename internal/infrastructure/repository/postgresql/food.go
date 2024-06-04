package postgresql

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/cast"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
	"time"
)

type foodRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewFood(db *postgres.PostgresDB) repo.Food {
	return &foodRepo{
		tableName: "foods",
		db:        db,
	}
}

func (a *foodRepo) Create(ctx context.Context, food *entity.Food) (*entity.Food, error) {
	query := `
	INSERT INTO foods (
	    id,
		name,
		capacity,
		product_union,
		description,
		created_at,
		updated_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING
		id,
		name,
		capacity,
		product_union,
		description
	`

	var (
		createdFood        entity.Food
		sqlNullDescription sql.NullString
	)

	err := a.db.QueryRow(ctx, query,
		food.ID,
		food.Name,
		food.Capacity,
		food.Union,
		food.Description,
		food.CreatedAt,
		food.UpdatedAt,
	).Scan(
		&createdFood.ID,
		&createdFood.Name,
		&createdFood.Capacity,
		&createdFood.Union,
		&sqlNullDescription,
	)

	if err != nil {
		return nil, err
	}

	if sqlNullDescription.Valid {
		createdFood.Description = sqlNullDescription.String
	}

	return &createdFood, nil
}

func (a *foodRepo) Update(ctx context.Context, food *entity.Food) (*entity.Food, error) {
	query := `
	UPDATE
		foods
	SET
		name = $1,
		capacity = $2,
		product_union = $3,
		description = $4,
		updated_at = $5
	WHERE
	    id = $6
		AND deleted_at IS NULL
	RETURNING
		id,
		name,
		capacity,
		product_union,
		description
	`

	var (
		updatedFood        entity.Food
		sqlNullDescription sql.NullString
	)

	err := a.db.QueryRow(ctx, query,
		food.Name,
		food.Capacity,
		food.Union,
		food.Description,
		food.UpdatedAt,
		food.ID,
	).Scan(
		&updatedFood.ID,
		&updatedFood.Name,
		&updatedFood.Capacity,
		&updatedFood.Union,
		&sqlNullDescription,
	)

	if err != nil {
		return nil, err
	}

	if sqlNullDescription.Valid {
		updatedFood.Description = sqlNullDescription.String
	}

	return &updatedFood, nil
}

func (a *foodRepo) Delete(ctx context.Context, foodID string) error {
	query := `UPDATE foods SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := a.db.Exec(ctx, query, time.Now().Format(time.RFC3339), foodID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (a *foodRepo) Get(ctx context.Context, params map[string]string) (*entity.Food, error) {
	var (
		food               entity.Food
		sqlNullDescription sql.NullString
	)

	queryBuilder := a.db.Sq.Builder.Select("id, name, capacity, product_union, description")
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
		&food.ID,
		&food.Name,
		&food.Capacity,
		&food.Union,
		&sqlNullDescription,
	)

	if err != nil {
		return nil, err
	}

	if sqlNullDescription.Valid {
		food.Description = sqlNullDescription.String
	}

	return &food, nil
}

func (a *foodRepo) List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListFoods, error) {
	var (
		offset = limit * (page - 1)
		foods  entity.ListFoods
	)

	queryBuilder := a.db.Sq.Builder.Select("id, name, capacity, product_union, description")
	queryBuilder = queryBuilder.From(a.tableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(a.db.Sq.ILike("name", "%"+cast.ToString(params["name"])+"%"))
	queryBuilder = queryBuilder.Where(a.db.Sq.ILike("product_union", "%"+cast.ToString(params["union"])+"%"))
	queryBuilder = queryBuilder.Limit(limit)
	queryBuilder = queryBuilder.Offset(offset)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			food               entity.Food
			sqlNullDescription sql.NullString
		)
		err := rows.Scan(
			&food.ID,
			&food.Name,
			&food.Capacity,
			&food.Union,
			&sqlNullDescription,
		)
		if err != nil {
			return nil, err
		}

		if sqlNullDescription.Valid {
			food.Description = sqlNullDescription.String
		}

		foods.Foods = append(foods.Foods, &food)
	}

	totalQueryBuilder := a.db.Sq.Builder.Select("COUNT(*)")
	totalQueryBuilder = totalQueryBuilder.From(a.tableName)
	totalQueryBuilder = totalQueryBuilder.Where("deleted_at IS NULL")
	totalQueryBuilder = totalQueryBuilder.Where(a.db.Sq.ILike("name", "%"+cast.ToString(params["name"])+"%"))
	totalQueryBuilder = totalQueryBuilder.Where(a.db.Sq.ILike("product_union", "%"+cast.ToString(params["union"])+"%"))

	totalQuery, totalArgs, err := totalQueryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var count = 0
	if err := a.db.QueryRow(ctx, totalQuery, totalArgs...).Scan(&count); err != nil {
		return nil, err
	}
	foods.TotalCount = uint64(count)

	return &foods, nil
}

func (d *foodRepo) UniqueFoodName(ctx context.Context, foodName string) (int, error) {
	query := `SELECT COUNT(*) FROM foods WHERE name = $1 AND deleted_at IS NULL`
	var count int

	if err := d.db.QueryRow(ctx, query, foodName).Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}
