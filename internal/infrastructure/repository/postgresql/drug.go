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

func (d *drugRepo) Create(ctx context.Context, drug *entity.Drug) (*entity.Drug, error) {
	query := `
	INSERT INTO drugs (
	    id,
	    name,
	    capacity,
	    product_union,
	    status,
	    description,
	    created_at,
	    updated_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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

	err := d.db.QueryRow(ctx, query,
		drug.ID,
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

func (d *drugRepo) Update(ctx context.Context, drug *entity.Drug) (*entity.Drug, error) {
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

	err := d.db.QueryRow(ctx, query,
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

func (d *drugRepo) Delete(ctx context.Context, drugID string) error {
	query := `UPDATE drugs SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := d.db.Exec(ctx, query, time.Now().Format(time.RFC3339), drugID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (d *drugRepo) Get(ctx context.Context, params map[string]string) (*entity.Drug, error) {
	var (
		drug               entity.Drug
		sqlNullDescription sql.NullString
	)

	queryBuilder := d.db.Sq.Builder.Select("id, name, capacity, product_union, status, description")
	queryBuilder = queryBuilder.From(d.tableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(d.db.Sq.Equal(key, value))
		}
		if key == "name" {
			queryBuilder = queryBuilder.Where(d.db.Sq.Equal(key, value))
		}
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	err = d.db.QueryRow(ctx, query, args...).Scan(
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

func (d *drugRepo) List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListDrugs, error) {
	var (
		offset = limit * (page - 1)
		drugs  = entity.ListDrugs{}
	)

	queryBuilder := d.db.Sq.Builder.Select("id, name, capacity, product_union, status, description")
	queryBuilder = queryBuilder.From(d.tableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(d.db.Sq.ILike("status", "%"+cast.ToString(params["status"])+"%"))
	queryBuilder = queryBuilder.Where(d.db.Sq.ILike("name", "%"+cast.ToString(params["name"])+"%"))
	queryBuilder = queryBuilder.Where(d.db.Sq.ILike("product_union", "%"+cast.ToString(params["union"])+"%"))
	queryBuilder = queryBuilder.Limit(limit)
	queryBuilder = queryBuilder.Offset(offset)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := d.db.Query(ctx, query, args...)
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

	totalQueryBuilder := d.db.Sq.Builder.Select("COUNT(*)")
	totalQueryBuilder = totalQueryBuilder.From(d.tableName)
	totalQueryBuilder = totalQueryBuilder.Where("deleted_at IS NULL")
	totalQueryBuilder = totalQueryBuilder.Where(d.db.Sq.ILike("status", "%"+cast.ToString(params["status"])+"%"))
	totalQueryBuilder = totalQueryBuilder.Where(d.db.Sq.ILike("name", "%"+cast.ToString(params["name"])+"%"))
	totalQueryBuilder = totalQueryBuilder.Where(d.db.Sq.ILike("product_union", "%"+cast.ToString(params["union"])+"%"))

	totalQuery, totalArgs, err := totalQueryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var count = 0
	if err := d.db.QueryRow(ctx, totalQuery, totalArgs...).Scan(&count); err != nil {
		return nil, err
	}
	drugs.TotalCount = uint64(count)

	return &drugs, nil
}

func (d *drugRepo) UniqueDrugName(ctx context.Context, drugName string) (int, error) {
	query := `SELECT COUNT(*) FROM drugs WHERE name = $1 AND deleted_at IS NULL`
	var count int

	if err := d.db.QueryRow(ctx, query, drugName).Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}
