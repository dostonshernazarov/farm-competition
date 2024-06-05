package postgresql

import (
	"context"
	"github.com/jackc/pgx/v4"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
	"time"
)

type feedingRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewFeeding(db *postgres.PostgresDB) repo.Feeding {
	return &feedingRepo{
		tableName: "animal_given_eatables",
		db:        db,
	}
}

func (f *feedingRepo) Create(ctx context.Context, feeding *entity.Feeding) (*entity.FeedingRes, error) {
	clauses := map[string]interface{}{
		"id":          feeding.ID,
		"animal_id":   feeding.AnimalID,
		"eatables_id": feeding.EatablesID,
		"category":    feeding.Category,
		"day":         feeding.Day,
		"daily":       feeding.Daily,
		"created_at":  feeding.CreatedAt,
		"updated_at":  feeding.UpdatedAt,
	}

	queryBuilder := f.db.Sq.Builder.Insert(f.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	selectQueryBuilder := f.db.Sq.Builder.Select("")
	selectQueryBuilder = selectQueryBuilder.From(f.tableName)

	return nil, nil
}

func (f *feedingRepo) Update(ctx context.Context, feeding *entity.Feeding) (*entity.FeedingRes, error) {
	clauses := map[string]interface{}{
		"animal_id":   feeding.AnimalID,
		"eatables_id": feeding.EatablesID,
		"category":    feeding.Category,
		"day":         feeding.Day,
		"daily":       feeding.Daily,
		"updated_at":  feeding.UpdatedAt,
	}

	queryBuilder := f.db.Sq.Builder.Insert(f.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	return nil, nil
}

func (f *feedingRepo) Delete(ctx context.Context, feedingID string) error {
	clauses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	queryBuilder := f.db.Sq.Builder.Update(f.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)
	queryBuilder = queryBuilder.Where("deleted_at is null")
	queryBuilder = queryBuilder.Where(f.db.Sq.Equal("id", feedingID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	result, err := f.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
