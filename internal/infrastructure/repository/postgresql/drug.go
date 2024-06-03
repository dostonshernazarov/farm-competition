package postgresql

import (
	"context"
	"database/sql"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
	"time"
)

type drugRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewDrug(db *postgres.PostgresDB) repo.Drug {
	return &drugRepo{
		tableName: "drugs",
		db:        db,
	}
}

func (a *drugRepo) Create(ctx context.Context, drug *entity.Drug) (*entity.Drug, error) {
	query := `
	INSERT INTO drugs (
	    name,
	    capacity,
	    product_union,
	    status,
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
		status,
		description
	`

	var (
		createdDrug        entity.Drug
		sqlNullDescription sql.NullString
	)

	err := a.db.QueryRow(ctx, query,
		drug.Name,
		drug.Capacity,
		drug.Union,
		drug.Status,
		drug.Description,
		drug.CreatedAt,
		drug.UpdatedAt,
	).Scan(
		&createdDrug.ID,
		&createdDrug.Name,
		&createdDrug.Capacity,
		&createdDrug.Union,
		&createdDrug.Status,
		&sqlNullDescription,
	)

	if err != nil {
		return nil, err
	}

	if sqlNullDescription.Valid {
		createdDrug.Description = sqlNullDescription.String
	}

	return &createdDrug, nil
}

func (a *drugRepo) Update(ctx context.Context, drug *entity.Drug) (*entity.Drug, error) {
	query := `
	UPDATE
		drugs
	SET
	    name = $1,
	    capacity = $2,
	    product_union = $3,
	    status = $4,
	    description = $5,
	    updated_at = $6
	WHERE
	    id = $7
		AND deleted_at IS NULL
	RETURNING
		id,
	    name,
		capacity,
		product_union,
		status,
		description
	`

	var (
		createdDrug        entity.Drug
		sqlNullDescription sql.NullString
	)

	err := a.db.QueryRow(ctx, query,
		drug.Name,
		drug.Capacity,
		drug.Union,
		drug.Status,
		drug.Description,
		drug.UpdatedAt,
		drug.ID,
	).Scan(
		&createdDrug.ID,
		&createdDrug.Name,
		&createdDrug.Capacity,
		&createdDrug.Union,
		&createdDrug.Status,
		&sqlNullDescription,
	)

	if err != nil {
		return nil, err
	}

	if sqlNullDescription.Valid {
		createdDrug.Description = sqlNullDescription.String
	}

	return &createdDrug, nil
}

func (a *drugRepo) Delete(ctx context.Context, drugID string) error {
	query := `UPDATE drugs SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := a.db.Exec(ctx, query, drugID, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}

	if result.RowsAffected() != 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (a *drugRepo) Get(ctx context.Context, drugID string) (*entity.Drug, error) {
	query := `
	SELECT
		id,
	    name,
		capacity,
		product_union,
		status,
		description
	FROM
	    drugs
	WHERE
	    id = $1
		AND deleted_at IS NULL
	`

	var (
		drug               entity.Drug
		sqlNullDescription sql.NullString
	)

	err := a.db.QueryRow(ctx, query, drugID).Scan(
		&drug.ID,
		&drug.Name,
		&drug.Capacity,
		&drug.Union,
		&drug.Status,
		&sqlNullDescription,
	)

	if err != nil {
		return nil, err
	}

	if sqlNullDescription.Valid {
		drug.Description = sqlNullDescription.String
	}

	return &drug, nil
}

func (a *drugRepo) List(ctx context.Context, page, limit uint64) (*entity.ListDrugs, error) {
	query := `
	SELECT
		id,
	    name,
		capacity,
		product_union,
		status,
		description
	FROM
	    drugs
	WHERE
		deleted_at IS NULL
	LIMIT $1
	OFFSET $2
	`

	var (
		drugs  entity.ListDrugs
		offset = limit * (page - 1)
	)

	rows, err := a.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			drug               entity.Drug
			sqlNullDescription sql.NullString
		)
		err := rows.Scan(
			&drug.ID,
			&drug.Name,
			&drug.Capacity,
			&drug.Union,
			&drug.Status,
			&sqlNullDescription,
		)
		if err != nil {
			return nil, err
		}
		if sqlNullDescription.Valid {
			drug.Description = sqlNullDescription.String
		}

		drugs.Drugs = append(drugs.Drugs, &drug)
	}

	var (
		count      = 0
		totalQuery = `SELECT COUNT(*) FROM drugs WHERE deleted_at IS NULL`
	)

	if err := a.db.QueryRow(ctx, totalQuery).Scan(&count); err != nil {
		return nil, err
	}
	drugs.TotalCount = uint64(count)

	return &drugs, nil
}
