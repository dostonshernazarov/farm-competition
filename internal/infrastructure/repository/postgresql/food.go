package postgresql

import (
	"context"
	"database/sql"
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
		name,
		capacity,
		product_union,
		description,
		created_at,
		updated_at
	)
	VALUES ($1, $2, $3, $4, $5, $6)
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

	result, err := a.db.Exec(ctx, query, foodID, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}

	if result.RowsAffected() != 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (a *foodRepo) Get(ctx context.Context, foodID string) (*entity.Food, error) {
	query := `
	SELECT
		id,
		name,
		capacity,
		product_union,
		description
	FROM
	    foods
	WHERE
	    id = $1
		AND deleted_at IS NULL
	`

	var (
		food               entity.Food
		sqlNullDescription sql.NullString
	)

	err := a.db.QueryRow(ctx, query, foodID).Scan(
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

func (a *foodRepo) List(ctx context.Context, page, limit uint64) (*entity.ListFoods, error) {
	query := `
	SELECT
		id,
		name,
		capacity,
		product_union,
		description
	FROM
	    foods
	WHERE
	  	deleted_at IS NULL
	LIMIT $1
	OFFSET $2
	`

	var (
		offset = limit * (page - 1)
		foods  entity.ListFoods
	)

	rows, err := a.db.Query(ctx, query, limit, offset)
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

	var (
		count      = 0
		totalQuery = `SELECT COUNT(*) FROM foods WHERE deleted_at IS NULL`
	)

	if err := a.db.QueryRow(ctx, totalQuery).Scan(&count); err != nil {
		return nil, err
	}
	foods.TotalCount = uint64(count)

	return &foods, nil
}
