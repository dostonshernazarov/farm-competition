package postgresql

import (
	"context"
	"database/sql"
	"github.com/spf13/cast"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
)

type deliveryRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewDelivery(db *postgres.PostgresDB) repo.Delivery {
	return &deliveryRepo{
		tableName: "into_store",
		db:        db,
	}
}

func (d *deliveryRepo) Create(ctx context.Context, delivery *entity.Delivery) (*entity.Delivery, error) {
	clauses := map[string]interface{}{
		"id":         delivery.ID,
		"name":       delivery.Name,
		"category":   delivery.Category,
		"capacity":   delivery.Capacity,
		"union":      delivery.Union,
		"time":       delivery.Time,
		"created_at": delivery.CreatedAt,
		"updated_at": delivery.UpdatedAt,
	}

	queryBuilder := d.db.Sq.Builder.Insert(d.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := d.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() == 0 {
		return nil, sql.ErrNoRows
	}

	return delivery, nil
}

func (d *deliveryRepo) Update(ctx context.Context, delivery *entity.Delivery) (*entity.Delivery, error) {
	clauses := map[string]interface{}{
		"name":       delivery.Name,
		"category":   delivery.Category,
		"capacity":   delivery.Capacity,
		"union":      delivery.Union,
		"time":       delivery.Time,
		"updated_at": delivery.UpdatedAt,
	}

	queryBuilder := d.db.Sq.Builder.Update(d.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(d.db.Sq.Equal("id", delivery.ID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := d.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() == 0 {
		return nil, sql.ErrNoRows
	}

	return delivery, nil
}

func (d *deliveryRepo) Delete(ctx context.Context, deliveryID string) error {
	queryBuilder := d.db.Sq.Builder.Delete(d.tableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(d.db.Sq.Equal("id", deliveryID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	result, err := d.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *deliveryRepo) Get(ctx context.Context, deliveryID string) (*entity.Delivery, error) {
	queryBuilder := d.db.Sq.Builder.Select("id, name, category, capacity, union, time")
	queryBuilder = queryBuilder.From(d.tableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(d.db.Sq.Equal("id", deliveryID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var delivery entity.Delivery
	err = d.db.QueryRow(ctx, query, args...).Scan(
		&delivery.ID,
		&delivery.Name,
		&delivery.Category,
		&delivery.Capacity,
		&delivery.Union,
		&delivery.Time,
	)
	if err != nil {
		return nil, err
	}

	return &delivery, nil
}

func (d *deliveryRepo) List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListDelivery, error) {
	queryBuilder := d.db.Sq.Builder.Select("id, name, category, capacity, union, time")
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(d.db.Sq.ILike("name", params["name"]))
	queryBuilder = queryBuilder.Where(d.db.Sq.ILike("category", params["category"]))
	queryBuilder = queryBuilder.Where(d.db.Sq.ILike("time", cast.ToString(params["time"])+"%"))
	queryBuilder = queryBuilder.Limit(limit)
	queryBuilder = queryBuilder.Offset(limit * (page - 1))
	queryBuilder = queryBuilder.OrderBy("time DESC")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := d.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveryList := entity.ListDelivery{}
	for rows.Next() {
		delivery := entity.Delivery{}
		err := rows.Scan(
			&delivery.ID,
			&delivery.Name,
			&delivery.Category,
			&delivery.Capacity,
			&delivery.Union,
			&delivery.Time,
		)
		if err != nil {
			return nil, err
		}

		deliveryList.Deliveries = append(deliveryList.Deliveries, &delivery)
	}

	totalQueryBuilder := d.db.Sq.Builder.Select("COUNT(*)")
	totalQueryBuilder = totalQueryBuilder.Where("deleted_at IS NULL")
	totalQueryBuilder = totalQueryBuilder.Where(d.db.Sq.ILike("name", params["name"]))
	totalQueryBuilder = totalQueryBuilder.Where(d.db.Sq.ILike("category", params["category"]))
	totalQueryBuilder = totalQueryBuilder.Where(d.db.Sq.ILike("time", cast.ToString(params["time"])+"%"))
	totalQuery, totalArgs, err := totalQueryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var count = 0
	if err := d.db.QueryRow(ctx, totalQuery, totalArgs...).Scan(&count); err != nil {
		return nil, err
	}
	deliveryList.TotalCount = int64(count)

	return &deliveryList, nil
}
